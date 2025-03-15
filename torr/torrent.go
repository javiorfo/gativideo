package torr

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	torrLog "github.com/anacrolix/log"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"github.com/javiorfo/bitsmuggler/config"
)

var configuration = config.GetConfiguration()

func MovieTorrentName(torrentPath string) (string, error) {
	resp, err := http.Get(torrentPath)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read torrent data: %w", err)
	}

	mi, err := metainfo.Load(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to parse torrent metadata: %w", err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return "", fmt.Errorf("failed to unmarshal info: %w", err)
	}

	return info.Name, nil
}

func DownloadTorrentFile(tFile string) {
	resp, err := http.Get(tFile)
	if err != nil {
		log.Fatalf("failed to get torrent file from %s: %v", tFile, err)
	}
	defer resp.Body.Close()

	destFile := filepath.Join(configuration.DownloadFolder, filepath.Base(tFile))
	out, err := os.Create(destFile)
	if err != nil {
		log.Fatalf("failed to create file: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatalf("failed to copy content: %v", err)
	}
}

type OnDownloadMsg string

var Update OnDownloadMsg

func DownloadMovies() {
	clientConfig := torrent.NewDefaultClientConfig()
	clientConfig.DataDir = configuration.DownloadFolder
	clientConfig.Logger.SetHandlers(torrLog.DiscardHandler)

	c, _ := torrent.NewClient(clientConfig)
	defer c.Close()

	// TODO search .torrent
	t, _ := c.AddTorrentFromFile(filepath.Join(configuration.DownloadFolder, "eric.torrent"))
	<-t.GotInfo()
	t.DownloadAll()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for !t.Seeding() {
		select {
		case <-ticker.C:
			// 			stats := t.Stats()
			totalSize := t.Info().TotalLength()
			s := t.BytesCompleted()
			// 			fmt.Printf("Mp4: %d \n", s)
			Update = OnDownloadMsg(fmt.Sprintf("Mp4: %d \n", s))
			/* 			fmt.Printf("Peers: %d/%d\n", stats.ActivePeers, stats.TotalPeers)
			   			fmt.Printf("Total: %d \n", totalSize)
			   			fmt.Printf("Downloaded: %.2f%%\n", (float64(s)/float64(totalSize))*100) */
			if t.BytesCompleted() == totalSize {
				// 				log.Print("Torrent downloaded")
				return
			}
		}
	}
}
