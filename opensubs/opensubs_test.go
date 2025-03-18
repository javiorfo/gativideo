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
