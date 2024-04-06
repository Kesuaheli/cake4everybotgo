package twitch

import (
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/util"
	webTwitch "cake4everybot/webserver/twitch"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

// HandleChannelUpdate is the event handler for the "channel.update" event from twitch.
func HandleChannelUpdate(s *discordgo.Session, e *webTwitch.ChannelUpdateEvent) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformTwitch, e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	updateEmbed := func(embed *discordgo.MessageEmbed) {
		embed.Description = e.Title
		if len(embed.Fields) == 0 {
			embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{})
		}
		embed.Fields[0].Name = lang.GetDefault("module.twitch.embed_category")
		embed.Fields[0].Value = e.CategoryName
	}

	for _, announcement := range announcements {
		err = updateAnnouncementMessage(s, announcement, updateEmbed)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

// HandleStreamOnline is the event handler for the "stream.online" event from twitch.
func HandleStreamOnline(s *discordgo.Session, e *webTwitch.StreamOnlineEvent) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformTwitch, e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	updateEmbed := func(embed *discordgo.MessageEmbed) {
		embed.Title = fmt.Sprintf("🔴 %s", e.BroadcasterUserName)
		embed.Color = 9520895
	}

	for _, announcement := range announcements {
		err = updateAnnouncementMessage(s, announcement, updateEmbed)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

// HandleStreamOffline is the event handler for the "stream.offline" event from twitch.
func HandleStreamOffline(s *discordgo.Session, e *webTwitch.StreamOfflineEvent) {
	announcements, err := database.GetAnnouncement(database.AnnouncementPlatformTwitch, e.BroadcasterUserID)
	if err == sql.ErrNoRows {
		return
	} else if err != nil {
		log.Printf("Error on get announcement: %v", err)
		return
	}

	updateEmbed := func(embed *discordgo.MessageEmbed) {
		embed.Title = fmt.Sprintf("⚫ %s", e.BroadcasterUserName)
		embed.Color = 2829358
	}

	for _, announcement := range announcements {
		err = updateAnnouncementMessage(s, announcement, updateEmbed)
		if err != nil {
			log.Printf("Error: %v", err)
		}
	}
}

func getAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement) (msg *discordgo.Message, err error) {
	channel, err := s.Channel(announcement.ChannelID)
	if err != nil {
		return nil, fmt.Errorf("get channel '%s': %v", announcement, err)
	}

	if channel.LastMessageID == "" {
		return newAnnouncementMessage(s, channel)
	}

	msg, err = s.ChannelMessage(channel.ID, channel.LastMessageID)
	if err != nil {
		if restErr, ok := err.(*discordgo.RESTError); ok {
			// if the lastMessageID returns a 404, i.e. it was deleted, create a new one
			if restErr.Response.StatusCode == http.StatusNotFound {
				return newAnnouncementMessage(s, channel)
			}
		}
		return nil, err
	}

	if msg.Author.ID != s.State.User.ID {
		msg, err = newAnnouncementMessage(s, channel)
	}

	return msg, err
}

func newAnnouncementMessage(s *discordgo.Session, channel *discordgo.Channel) (*discordgo.Message, error) {
	embed := &discordgo.MessageEmbed{
		Description: "-",
	}
	util.SetEmbedFooter(s, "module.twitch.embed_footer", embed)
	return s.ChannelMessageSendEmbed(channel.ID, embed)
}

func updateAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement, updateEmbed func(embed *discordgo.MessageEmbed)) error {
	msg, err := getAnnouncementMessage(s, announcement)
	if err != nil {
		return fmt.Errorf("get announcement in channel '%s': %v", announcement, err)
	}

	var embed *discordgo.MessageEmbed
	if len(msg.Embeds) == 0 {
		embed = &discordgo.MessageEmbed{}
	} else {
		embed = msg.Embeds[0]
	}

	updateEmbed(embed)

	_, err = s.ChannelMessageEditEmbed(announcement.ChannelID, msg.ID, embed)
	if err != nil {
		return fmt.Errorf("update announcement in channel '%s': %v", announcement, err)
	}
	return nil
}
