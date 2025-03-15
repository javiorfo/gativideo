package opensubs

import "testing"

func TestOpenSubs(t *testing.T) {
	subtitles := GetSubs("2008", "Sleepwalking")

	if len(subtitles) == 0 {
		t.Fatal("subtitles must not be empty")
	}

	code := subtitles[0].GetDownloadSubtitleCode()
	if code.IsEmpty() {
		t.Fatal("code must not be empty")
	}
}
