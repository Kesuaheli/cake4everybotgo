package twitch

import (
	webTwitch "cake4everybot/webserver/twitch"
	"encoding/json"

	"github.com/bwmarrin/discordgo"
)

func Announce(s *discordgo.Session, e *webTwitch.RawEvent) {
	log.Printf("New twitch event:\n%+v\n", e)
	data, _ := json.Marshal(e.Event)
	switch e.Subscription.Type {
	case "channel.update":
		var e2 webTwitch.ChannelUpdateEvent
		err := json.Unmarshal(data, &e2)
		if err != nil {
			log.Printf("failed to parse channel update event: %v", err)
			return
		}
		log.Printf("Channel were updated: %v", e2)
	case "stream.online":
		var e2 webTwitch.StreamOnlineEvent
		err := json.Unmarshal(data, &e2)
		if err != nil {
			log.Printf("failed to parse online event: %v", err)
			return
		}
		log.Printf("Stream went online: %v", e2)
	case "stream.offline":
		var e2 webTwitch.StreamOfflineEvent
		err := json.Unmarshal(data, &e2)
		if err != nil {
			log.Printf("failed to parse offline event: %v", err)
			return
		}
		log.Printf("Stream went offline: %v", e2)
	}
}
