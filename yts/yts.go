package yts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/javiorfo/bitsmuggler/config"
	"github.com/javiorfo/steams"
)

var configuration = config.GetConfiguration()

type Movie struct {
	Name      string
	Year      string
	Genre     string
	Rate      string
	TechSpecs []TechSpec
}

func (m Movie) GetTechSpec() TechSpec {
	quality := strconv.Itoa(int(configuration.YTSQuality))
	return steams.OfSlice(m.TechSpecs).FindOne(func(t TechSpec) bool {
		switch quality {
		case "2160":
			if strings.Contains(t.TorrentFileLink, "2160") {
				return true
			}
			fallthrough
		case "1080":
			if strings.Contains(t.TorrentFileLink, "1080") {
				return true
			}
			fallthrough
		case "720":
			if strings.Contains(t.TorrentFileLink, "720") {
				return true
			}
		}
		return false
	}).OrElseGet(techSpecNotFound)
}

type TechSpec struct {
	TorrentFileLink string
	Size            string
	Resolution      string
	Duration        string
	Language        string
}

func techSpecNotFound() TechSpec {
	return TechSpec{
		Size:       "Unknown",
		Resolution: "Unknown",
		Duration:   "Unknown",
		Language:   "Unknown",
	}
}

func GetMovies(host, keyword, genre, rating, year, order string, page int) (int, []Movie) {
	c := colly.NewCollector()

	var movies []Movie
	c.OnHTML("div#movie-content", func(e *colly.HTMLElement) {
		data := strings.Split(e.ChildText("div div#movie-info div"), "\n")
		tFiles := e.ChildAttrs("div div#movie-info div p a", "href")

		techSpecInfo := strings.Split(e.ChildText("div#movie-tech-specs div.tech-spec-info div div"), "\n")

		techSpecs := make([]TechSpec, len(tFiles))
		size := 0
		resolution := 1
		language := 2
		duration := 7
		for i, t := range tFiles {
			if !strings.Contains(t, ".torrent") {
				continue
			}
			techSpecs[i].TorrentFileLink = configuration.YTSHost + t
			sizeStr := strings.TrimSpace(techSpecInfo[size])
			if strings.Contains(sizeStr, "P/S") {
				i := strings.LastIndex(sizeStr[:len(sizeStr)-3], " ")
				sizeStr = sizeStr[i+1:]
			}

			techSpecs[i].Size = sizeStr
			techSpecs[i].Resolution = strings.TrimSpace(techSpecInfo[resolution])
			techSpecs[i].Language = strings.TrimSpace(techSpecInfo[language])
			techSpecs[i].Duration = strings.TrimSpace(techSpecInfo[duration])
			size += 8
			resolution += 8
			language += 8
			duration += 8
		}

		movie := Movie{
			Name:      strings.TrimSpace(data[0]),
			Year:      strings.TrimSpace(data[1]),
			Genre:     strings.TrimSpace(data[2]),
			Rate:      strings.TrimSpace(data[len(data)-1]),
			TechSpecs: techSpecs,
		}
		movies = append(movies, movie)
	})

	var total int
	c.OnHTML("div.browse-content", func(e *colly.HTMLElement) {
		total, _ = strconv.Atoi(e.ChildText("div h2 b"))
		movies := e.ChildAttrs("div section div div a.browse-movie-title", "href")
		for _, v := range movies {
			e.Request.Visit(v)
		}
	})

	c.Visit(fmt.Sprintf("%s/browse-movies?keyword=%s&quality=all&genre=%s&rating=%s&year=%s&order_by=%s&page=%d", host, keyword, genre, rating, year, order, page))
	return total, movies
}
