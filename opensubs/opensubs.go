package opensubs

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/javiorfo/bitsmuggler/config"
	"github.com/javiorfo/steams"
)

var configuration = config.GetConfiguration()

type response struct {
	Data [][]string `json:"data"`
}

type Subtitle struct {
	Name      string
	Downloads string
	Date      string
	Link      string
}

func (s Subtitle) GetDownloadSubtitleCode() string {
	c := colly.NewCollector()

	var subCode string
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "opensubtitles.org/"+configuration.OpenSubsLanguage+"/subtitles") {
			subCode = href[strings.LastIndex(href, "/")+1:]
		}
	})

	c.Visit(s.Link)

	return subCode
}

func GetSubs(movieYear, movieName string) ([]Subtitle, error) {
	movieName = getOpenSubsMovieName(movieName)

	url := fmt.Sprintf("https://www.opensubtitles.com/%s/%s/features/%s-%s/subtitles.json", configuration.OpenSubsLanguage, configuration.OpenSubsLanguage, movieYear, movieName)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	var subtitles []Subtitle
	for index := range response.Data {
		var sub Subtitle
		for i, v := range response.Data[index] {
			if i == 3 {
				sub.Date = v
			}
			if i == 2 {
				name := strings.TrimSuffix(v, "</a>")
				link := strings.TrimPrefix(v, "<a href=\"")
				sub.Name = name[strings.Index(name, "\">")+2 : strings.Index(name, "</a>")]
				sub.Link = "https://www.opensubtitles.com" + link[:strings.Index(link, "\">")]
			}
			if i == 8 {
				download := v[strings.Index(v, "\">")+2:]
				sub.Downloads = strings.TrimSuffix(download, "</a>")
			}
		}
		subtitles = append(subtitles, sub)
	}
	return steams.OfSlice(subtitles).Sorted(sorted).Collect(), nil
}

func sorted(s1 Subtitle, s2 Subtitle) bool {
	lower := strings.ToLower(s1.Name)
	if strings.Contains(lower, "yify") || strings.Contains(lower, "yifi") || strings.Contains(lower, "yts") {
		return true
	}
	a, _ := strconv.Atoi(s1.Downloads)
	b, _ := strconv.Atoi(s2.Downloads)
	return a > b
}

func getOpenSubsMovieName(input string) string {
	re := regexp.MustCompile(`[.,:]`).ReplaceAllString(input, "")
	return strings.ReplaceAll(strings.ToLower(re), " ", "-")
}

func DownloadSubtitle(code, movieName string) error {
	resp, err := http.Get("https://dl.opensubtitles.org/en/download/sub/" + code)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, resp.Body)
	if err != nil {
		return err
	}

	zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configuration.DownloadFolder, 0755); err != nil {
		return err
	}

	for _, file := range zipReader.File {
		if file.FileInfo().IsDir() || !strings.Contains(file.Name, ".srt") {
			continue
		}

		filePath := filepath.Join(configuration.DownloadFolder, movieName+".srt")
		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			fileReader.Close()
			return err
		}

		_, err = io.Copy(outFile, fileReader)
		if err != nil {
			fileReader.Close()
			outFile.Close()
			return err
		}

		fileReader.Close()
		outFile.Close()
	}
	return nil
}
