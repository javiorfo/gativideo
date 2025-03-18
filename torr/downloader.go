package torr

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	torrLog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/javiorfo/bitsmuggler/config"
	"github.com/javiorfo/nilo"
)

var configuration = config.GetConfiguration()

type OnDownload string

var Status OnDownload

type Downloader struct {
	torrentPath string
	torrentName string
	Downloading bool
}

func (d *Downloader) Scan(cancelDownload <-chan struct{}) {
	searchTorrent().IfPresent(func(torrentPath string) {
		d.Downloading = true
		d.torrentPath = torrentPath
		go d.downloadMovie(cancelDownload)
	})
}

func (d *Downloader) DownloadTorrentFile(tFile string) error {
	resp, err := http.Get(tFile)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	d.torrentPath = filepath.Join(configuration.DownloadFolder, filepath.Base(tFile))
	out, err := os.Create(d.torrentPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func (d *Downloader) downloadMovie(cancelDownload <-chan struct{}) {
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = configuration.DownloadFolder
	clientConfig.Logger.SetHandlers(torrLog.DiscardHandler)

	c, _ := torrent.NewClient(clientConfig)
	defer c.Close()

	t, _ := c.AddTorrentFromFile(d.torrentPath)

	<-t.GotInfo()
	t.DownloadAll()
	d.torrentName = t.Name()
	for !t.Seeding() {
		select {
		case <-cancelDownload:
			Status = OnDownload(fmt.Sprintf("  %s Canceled!", t.Name()))
			_ = os.Remove(d.torrentPath)
			_ = os.RemoveAll(filepath.Join(configuration.DownloadFolder, d.torrentName))
			return
		default:
			stats := t.Stats()
			totalSize := t.Info().TotalLength()
			completed := t.BytesCompleted()
			Status = OnDownload(fmt.Sprintf("  %s | Progress %.2f%% | Peers %d/%d", d.torrentName, (float64(completed)/float64(totalSize))*100.0, stats.ActivePeers, stats.TotalPeers))
			if completed == totalSize {
				Status = OnDownload(fmt.Sprintf("  %s Completed!", d.torrentName))
				_ = d.purge()
				return
			}
		}
		time.Sleep(1 * time.Second)
	}
	c.WaitAll()
}

func (d *Downloader) purge() error {
	_ = os.Remove(d.torrentPath)
	movieTorrentDir := filepath.Join(configuration.DownloadFolder, d.torrentName)
	subtitlePath := movieTorrentDir + ".srt"
	subtitleDestPath := filepath.Join(movieTorrentDir, d.torrentName+".srt")

	if _, err := os.Stat(subtitlePath); err != nil {
		return err
	}

	err := os.Rename(subtitlePath, subtitleDestPath)
	if err != nil {
		return err
	}

	files, err := os.ReadDir(movieTorrentDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			filePath := filepath.Join(movieTorrentDir, fileName)
			if strings.HasSuffix(fileName, ".mp4") {
				_ = os.Rename(filePath, filepath.Join(movieTorrentDir, d.torrentName+".mp4"))
			} else {
				filePath := filepath.Join(movieTorrentDir, fileName)
				_ = os.Remove(filePath)
			}
		}
	}
	return nil
}

func MovieTorrentName(torrentPath string) (string, error) {
	resp, err := http.Get(torrentPath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	mi, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", err
	}

	return info.Name, nil
}

func searchTorrent() nilo.Optional[string] {
	var optional nilo.Optional[string]
	err := filepath.Walk(configuration.DownloadFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(info.Name()) == ".torrent" {
			optional = nilo.Of(path)
		}
		return nil
	})

	if err != nil {
		return nilo.Empty[string]()
	}
	return optional
}
