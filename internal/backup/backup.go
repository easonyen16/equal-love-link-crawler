package backup

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/easonyen16/equal-love-link-crawler/api/message"
)

const pageSize = 100

var jstLoc *time.Location

func init() {
	var err error
	jstLoc, err = time.LoadLocation("Asia/Tokyo")
	if err != nil {
		jstLoc = time.FixedZone("JST", 9*60*60)
	}
}

func toJSTFilename(unixSec int64) string {
	return time.Unix(unixSec, 0).In(jstLoc).Format("20060102150405")
}

func setFileTimes(path string, unixSec int64) error {
	t := time.Unix(unixSec, 0)
	return os.Chtimes(path, t, t)
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func downloadFile(url, destPath string, unixSec int64) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, resp.Body); err != nil {
		return err
	}

	return setFileTimes(destPath, unixSec)
}

func saveMessage(dir string, msg message.ChatMessage) (alreadyExists bool, err error) {
	baseName := toJSTFilename(msg.PostedDate)
	txtPath := filepath.Join(dir, baseName+".txt")

	if fileExists(txtPath) {
		return true, nil
	}

	if err := os.WriteFile(txtPath, []byte(msg.TextContent), 0644); err != nil {
		return false, err
	}
	if err := setFileTimes(txtPath, msg.PostedDate); err != nil {
		return false, err
	}

	if len(msg.ChatMedia) == 1 {
		m := msg.ChatMedia[0]
		mediaPath := filepath.Join(dir, baseName+"."+m.FileExtension)
		if err := downloadFile(m.URL, mediaPath, msg.PostedDate); err != nil {
			return false, fmt.Errorf("download media: %w", err)
		}
	} else if len(msg.ChatMedia) > 1 {
		for i, m := range msg.ChatMedia {
			mediaPath := filepath.Join(dir, fmt.Sprintf("%s-%d.%s", baseName, i+1, m.FileExtension))
			if err := downloadFile(m.URL, mediaPath, msg.PostedDate); err != nil {
				return false, fmt.Errorf("download media %d: %w", i+1, err)
			}
		}
	}

	return false, nil
}

func Room(accessToken, downloadDir string, room message.TalkRoom) error {
	dir := filepath.Join(downloadDir, room.Name)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	page := 1
	pageStartID := 0

	for {
		chatPage, err := message.GetChat(accessToken, room.ID, page, pageStartID)
		if err != nil {
			return err
		}

		for _, msg := range chatPage.Messages {
			if msg.PostedArtistID == 0 {
				continue
			}

			alreadyExists, err := saveMessage(dir, msg)
			if err != nil {
				return fmt.Errorf("save message %d: %w", msg.ID, err)
			}
			if alreadyExists {
				fmt.Printf("  [完了] 既存ファイルを検出、%s のバックアップを終了\n", room.Name)
				return nil
			}
		}

		fmt.Printf("  %s: %d 件保存 (page %d)\n", room.Name, len(chatPage.Messages), page)

		if len(chatPage.Messages) < pageSize || chatPage.NextPageID == 0 {
			break
		}

		page++
		pageStartID = chatPage.NextPageID
	}

	return nil
}

func All(accessToken, downloadDir string, rooms []message.TalkRoom) {
	for _, room := range rooms {
		if !room.IsAccessible {
			continue
		}
		fmt.Printf("\n%s のバックアップ開始...\n", room.Name)
		if err := Room(accessToken, downloadDir, room); err != nil {
			fmt.Printf("  エラー: %v\n", err)
		}
	}
}
