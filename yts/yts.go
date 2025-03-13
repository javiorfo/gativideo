package yts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/javiorfo/bitsmuggler/config"
	"github.com/javiorfo/steams"
)

var configuration = config.GetConfig()

type Movie struct {
	Name     string
	Year     string
	Genre    string
	Rate     string
	Torrents []Torrent
}

func (m Movie) GetTorrent() Torrent {
	quality := strconv.Itoa(int(configuration.YTSQuality))
	return steams.OfSlice(m.Torrents).FindOne(func(t Torrent) bool {
		switch quality {
		case "2160":
			if strings.Contains(t.File, "2160") {
				return true
			}
			fallthrough
		case "1080":
			if strings.Contains(t.File, "1080") {
				return true
			}
			fallthrough
		case "720":
			if strings.Contains(t.File, "720") {
				return true
			}
		}
		return false
	}).OrElseGet(torrentNotFound)
}

type Torrent struct {
	File       string
	Size       string
	Resolution string
	Duration   string
	Language   string
}

func torrentNotFound() Torrent {
	return Torrent{
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
		files := e.ChildAttrs("div div#movie-info div p a", "href")

		torrentInfo := strings.Split(e.ChildText("div#movie-tech-specs div.tech-spec-info div div"), "\n")

		torrents := make([]Torrent, len(files))
		size := 0
		resolution := 1
		language := 2
		duration := 7
		for i, v := range files {
			if !strings.Contains(v, ".torrent") {
				continue
			}
			torrents[i].File = v
			sizeStr := strings.TrimSpace(torrentInfo[size])
			if strings.Contains(sizeStr, "P/S") {
				i := strings.LastIndex(sizeStr[:len(sizeStr)-3], " ")
				sizeStr = sizeStr[i+1:]
			}

			torrents[i].Size = sizeStr
			torrents[i].Resolution = strings.TrimSpace(torrentInfo[resolution])
			torrents[i].Language = strings.TrimSpace(torrentInfo[language])
			torrents[i].Duration = strings.TrimSpace(torrentInfo[duration])
			size += 8
			resolution += 8
			language += 8
			duration += 8
		}

		movie := Movie{
			Name:     strings.TrimSpace(data[0]),
			Year:     strings.TrimSpace(data[1]),
			Genre:    strings.TrimSpace(data[2]),
			Rate:     strings.TrimSpace(data[len(data)-1]),
			Torrents: torrents,
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
