package twitch

import (
	"cake4everybot/webserver/twitch"

	"github.com/bwmarrin/discordgo"
)

func Announce(s *discordgo.Session, e *twitch.RawEvent) {
	log.Printf("New twitch event:\n%+v\n", e)
}
