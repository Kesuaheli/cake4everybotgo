package twitch

import (
	webTwitch "cake4everybot/webserver/twitch"

	"github.com/bwmarrin/discordgo"
)

func HandleChannelUpdate(s *discordgo.Session, e *webTwitch.ChannelUpdateEvent) {
	log.Printf("Channel were updated: %v", e)
}

func HandleStreamOnline(s *discordgo.Session, e *webTwitch.StreamOnlineEvent) {
	log.Printf("Stream went online: %v", e)
}

func HandleStreamOfflineEvent(s *discordgo.Session, e *webTwitch.StreamOfflineEvent) {
	log.Printf("Stream went offline: %v", e)
}
