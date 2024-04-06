package twitch

import (
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/util"
	webTwitch "cake4everybot/webserver/twitch"
	"database/sql"
	"fmt"

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
	if announcement.MessageID != "" {
		return s.ChannelMessage(announcement.ChannelID, announcement.MessageID)
	}

	msgs, err := s.ChannelMessages(announcement.ChannelID, 1, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("get last message: %v", err)
	}

	if len(msgs) == 0 {
		return newAnnouncementMessage(s, announcement)
	}

	msg = msgs[0]
	if msg.Author.ID != s.State.User.ID {
		msg, err = newAnnouncementMessage(s, announcement)
	}

	return msg, err
}

func newAnnouncementMessage(s *discordgo.Session, announcement *database.Announcement) (msg *discordgo.Message, err error) {
	embed := &discordgo.MessageEmbed{
		Description: "-",
	}
	util.SetEmbedFooter(s, "module.twitch.embed_footer", embed)

	msg, err = s.ChannelMessageSendEmbed(announcement.ChannelID, embed)
	if err != nil {
		return
	}
	return msg, announcement.UpdateAnnouncementMessage(msg.ID)
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
