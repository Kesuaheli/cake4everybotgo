package twitch

import (
	"cake4everybot/database"
	webTwitch "cake4everybot/webserver/twitch"
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func HandleChannelUpdate(s *discordgo.Session, e *webTwitch.ChannelUpdateEvent) {
	announcements, err := database.GetAnnouncement("twitch", e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	log.Printf("Channel were updated ('%d' server): %v", len(announcements), e)
}

func HandleStreamOnline(s *discordgo.Session, e *webTwitch.StreamOnlineEvent) {
	announcements, err := database.GetAnnouncement("twitch", e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	log.Printf("Stream went online ('%d' server): %v", len(announcements), e)
}

func HandleStreamOfflineEvent(s *discordgo.Session, e *webTwitch.StreamOfflineEvent) {
	announcements, err := database.GetAnnouncement("twitch", e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	log.Printf("Stream went offline ('%d' server): %v", len(announcements), e)
}
