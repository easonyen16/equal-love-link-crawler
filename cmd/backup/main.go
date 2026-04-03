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

	loginResp, err := message.Login(email, password)
	if err != nil {
		log.Fatal(err)
	}

	talkRooms, err := message.GetTalkRooms(loginResp.AccessToken)
	if err != nil {
		log.Fatal(err)
	}

	printTalkRoomSummary(talkRooms)
	backup.All(loginResp.AccessToken, "download", talkRooms)
}
