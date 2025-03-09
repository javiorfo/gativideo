package yts

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type Movie struct {
	Name     string
	Year     string
	Genre    string
	Rate     string
	Torrents []Torrent
}

type Torrent struct {
	File       string
	Size       string
	Resolution string
	Duration   string
	Language   string
}

func GetMovies(host, keyword, quality, genre, rating, year, order string, page int) (int, []Movie) {
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
			Genre:    strings.ReplaceAll(strings.TrimSpace(data[2]), "\\u00a0", ""),
			Rate:     strings.TrimSpace(data[len(data)-1]),
			Torrents: torrents,
		}
		movies = append(movies, movie)
	})

	// 	c.OnHTML("div.browse-movie-bottom", func(e *colly.HTMLElement) {
    var total int
	c.OnHTML("div.browse-content", func(e *colly.HTMLElement) {
        total, _ = strconv.Atoi(e.ChildText("div h2 b"))
		movies := e.ChildAttrs("div section div div a.browse-movie-title", "href")
		for _, v := range movies {
			e.Request.Visit(v)
		}
		// 		e.Request.Visit(e.ChildAttr("div section div div a.browse-movie-title", "href"))
	})

	c.Visit(fmt.Sprintf("%s/browse-movies?keyword=%s&quality=%s&genre=%s&rating=%s&year=%s&order_by=%s&page=%d", host, keyword, quality, genre, rating, year, order, page))
	return total, movies
}
