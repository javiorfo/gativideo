package yts

import "testing"

func TestGetMovies(t *testing.T) {
	total, movies := GetMovies(configuration.YTSHost, "", "all", "0", "2000", "rating", 1)

	if len(movies) != 20 {
		t.Fatal("movies must be 20")
	}

	if total == 0 {
		t.Fatal("total must be 20")
	}
}
