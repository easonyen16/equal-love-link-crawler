package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/easonyen16/equal-love-link-crawler/api/message"
	"github.com/easonyen16/equal-love-link-crawler/internal/backup"
	"golang.org/x/term"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("グループ選択:")
	fmt.Println("  1) =LOVE")
	fmt.Println("  2) ≠ME")
	fmt.Println("  3) ≒JOY")
	fmt.Print("グループを選んでください (1/2/3): ")
	scanner.Scan()
	platformInput := strings.TrimSpace(scanner.Text())

	var platform message.Platform
	switch platformInput {
	case "2":
		platform = message.PlatformNotEqualMe
	case "3":
		platform = message.PlatformNearlyEqualJoy
	default:
		platform = message.PlatformEqualLove
	}

	client := message.NewClient(platform)

	fmt.Print("Email: ")
	scanner.Scan()
	email := strings.TrimSpace(scanner.Text())

	fmt.Print("Password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		log.Fatal(err)
	}
	password := string(passwordBytes)

	loginResp, err := client.Login(email, password)
	if err != nil {
		log.Fatal(err)
	}

	talkRooms, err := client.GetTalkRooms(loginResp.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	printTalkRoomSummary(talkRooms)
	backup.All(client, loginResp.AccessToken, "download", talkRooms)
}
