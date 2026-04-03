package main

import (
	"fmt"
	"strings"

	"github.com/easonyen16/equal-love-link-crawler/api/message"
)

func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r > 0x7F {
			w += 2
		} else {
			w++
		}
	}
	return w
}

func padRight(s string, width int) string {
	diff := width - displayWidth(s)
	if diff <= 0 {
		return s
	}
	return s + strings.Repeat(" ", diff)
}

func printRooms(rooms []message.TalkRoom, maxWidth int) {
	for _, room := range rooms {
		fmt.Printf("  %s  %s\n", padRight(room.Name, maxWidth), room.RomajiName)
	}
}

func printTalkRoomSummary(talkRooms []message.TalkRoom) {
	var subscribed, unsubscribed []message.TalkRoom
	maxWidth := 0
	maxRomajiWidth := 0
	for _, room := range talkRooms {
		if w := displayWidth(room.Name); w > maxWidth {
			maxWidth = w
		}
		if w := len(room.RomajiName); w > maxRomajiWidth {
			maxRomajiWidth = w
		}
		if room.IsAccessible {
			subscribed = append(subscribed, room)
		} else {
			unsubscribed = append(unsubscribed, room)
		}
	}

	lineWidth := maxWidth + 2 + maxRomajiWidth
	printTitle := func(title string) {
		pad := (lineWidth - displayWidth(title)) / 2
		fmt.Printf("\n%s%s\n", strings.Repeat(" ", pad), title)
	}

	total := len(talkRooms)
	printTitle(fmt.Sprintf("=== 購読済み (%d/%d) ===", len(subscribed), total))
	printRooms(subscribed, maxWidth)

	printTitle(fmt.Sprintf("=== 未購読 (%d/%d) ===", len(unsubscribed), total))
	printRooms(unsubscribed, maxWidth)
}
