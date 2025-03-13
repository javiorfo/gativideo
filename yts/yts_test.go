package yts

import "testing"

func TestGetMovies(t *testing.T) {
	total, movies := GetMovies(configuration.YTSHost, "", "all", "0", "2000", "rating", 1)
	t.Log(total, movies)
}
