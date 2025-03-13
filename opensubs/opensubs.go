package opensubs

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/javiorfo/bitsmuggler/config"
	"github.com/javiorfo/nilo"
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

func (s Subtitle) GetDownloadSubCode() nilo.Optional[string] {
	c := colly.NewCollector()

	var subCode string
	c.OnHTML("a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		if strings.Contains(href, "opensubtitles.org/"+configuration.OpenSubsLanguage+"/subtitles") {
			subCode = href[strings.LastIndex(href, "/")+1:]
		}
	})

	c.Visit(s.Link)

	if subCode == "" {
		return nilo.Empty[string]()
	}
	return nilo.Of(subCode)
}

func GetSubs(movieYear, movieName string) []Subtitle {
	movieName = strings.ReplaceAll(strings.ToLower(movieName), " ", "-")

	url := fmt.Sprintf("https://www.opensubtitles.com/%s/%s/features/%s-%s/subtitles.json", configuration.OpenSubsLanguage, configuration.OpenSubsLanguage, movieYear, movieName)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	var response response
	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return nil
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
	return steams.OfSlice(subtitles).Sorted(sorted).Collect()
}

func sorted(s1 Subtitle, s2 Subtitle) bool {
	a, _ := strconv.Atoi(s1.Downloads)
	b, _ := strconv.Atoi(s2.Downloads)
	return a > b
}

func DownloadZip(code, zipName string) {
	resp, err := http.Get("https://dl.opensubtitles.org/en/download/sub/" + code)
	if err != nil {
		log.Fatalf("failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(zipName)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("failed to copy content: %v", err)
	}
}
