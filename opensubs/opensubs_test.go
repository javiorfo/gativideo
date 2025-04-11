package opensubs

import "testing"

func TestOpenSubs(t *testing.T) {
	subtitles, err := GetSubs("2008", "Sleepwalking")

	if err != nil {
		t.Fatal("subtitles must not be empty")
	}

	code := subtitles[0].GetDownloadSubtitleCode()
	if code == "" {
		t.Fatal("code must not be empty")
	}
}

func TestCleanString(t *testing.T) {
	s := getOpenSubsMovieName("Movie: Name")

	if s != "movie-name" {
		t.Fatalf("Result %s", s)
	}

	s = getOpenSubsMovieName("Movie, Name. 1990")
	if s != "movie-name-1990" {
		t.Fatalf("Result %s", s)
	}
}
