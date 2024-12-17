// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"

	"cake4everybot/data/lang"
	"cake4everybot/database"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

var commandIDMap map[string]string

// SetCommandMap sets the map from command names to ther registered ID.
// TODO: move the original command.CommandMap in a seperate Package to avoid this.
func SetCommandMap(m map[string]string) {
	commandIDMap = m
}

// AuthoredEmbed returns a new Embed with an author and footer set.
//
//	author:
//		The name and icon in the author field
//		of the embed.
//	sectionName:
//		The translation key used in the standard footer.
func AuthoredEmbed[T *discordgo.User | *discordgo.Member](s *discordgo.Session, author T, sectionName string) *discordgo.MessageEmbed {
	var username string
	user, ok := any(author).(*discordgo.User)
	if !ok {
		member, ok := any(author).(*discordgo.Member)
		if !ok {
			panic("Given generic type is not an discord user or member")
		}
		user = member.User
		username = member.DisplayName()
	}

	if username == "" {
		username = user.Username
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    username,
			IconURL: user.AvatarURL(""),
		},
	}
	SetEmbedFooter(s, sectionName, embed)
	return embed
}

// SetEmbedFooter takes a pointer to an embeds and sets the standard footer with the given name.
//
//	sectionName:
//		translation key for the name
func SetEmbedFooter(s *discordgo.Session, sectionName string, e *discordgo.MessageEmbed) {
	botName := viper.GetString("discord.name")
	name := lang.Get(sectionName, lang.FallbackLang())

	if e == nil {
		e = &discordgo.MessageEmbed{}
	}
	e.Footer = &discordgo.MessageEmbedFooter{
		Text:    fmt.Sprintf("%s > %s", botName, name),
		IconURL: s.State.User.AvatarURL(""),
	}
}

// AddEmbedField is a short hand for appending one field to the embed
func AddEmbedField(e *discordgo.MessageEmbed, name, value string, inline bool) {
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{Name: name, Value: value, Inline: inline})
}

// AddReplyHiddenField appends the standard field for ephemral embeds to the existing fields of the
// given embed.
func AddReplyHiddenField(e *discordgo.MessageEmbed) {
	AddEmbedField(e,
		lang.Get("discord.command.generic.msg.self_hidden", lang.FallbackLang()),
		lang.Get("discord.command.generic.msg.self_hidden.desc", lang.FallbackLang()),
		false,
	)
}

// MentionCommand returns the mention string for a slashcommand
func MentionCommand(base string, subcommand ...string) string {
	cBase := lang.GetDefault(base)

	cID := commandIDMap[cBase]
	if cID == "" {
		return ""
	}

	var cSub string
	for _, sub := range subcommand {
		cSub = cSub + " " + lang.GetDefault(sub)
	}

	return fmt.Sprintf("</%s%s:%s>", cBase, cSub, cID)
}

// GetChannelsFromDatabase returns a map from guild IDs to channel IDs
func GetChannelsFromDatabase(s *discordgo.Session, channelName string) (map[string]string, error) {
	rows, err := database.Query("SELECT id," + channelName + " FROM guilds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	IDMap := map[string]string{}
	for rows.Next() {
		var guildInt, channelInt uint64
		err := rows.Scan(&guildInt, &channelInt)
		if err != nil {
			return nil, err
		}
		if channelInt == 0 {
			continue
		}
		guildID := fmt.Sprint(guildInt)
		channelID := fmt.Sprint(channelInt)

		// validate channel
		channel, err := s.Channel(channelID)
		if err != nil {
			log.Printf("Warning: could not get %s channel for id '%s: %+v\n", channelName, channelID, err)
			continue
		}
		if channel.GuildID != guildID {
			log.Printf("Warning: tried to get %s channel (from channel/%s/%s), but this channel is from guild: '%s'\n", channelName, guildID, channelID, channel.GuildID)
			continue
		}

		IDMap[guildID] = channelID
	}

	return IDMap, nil
}

// GetConfigComponentEmoji returns a configured [discordgo.ComponentEmoji] for the given name.
func GetConfigComponentEmoji(name string) *discordgo.ComponentEmoji {
	e := GetConfigEmoji(name)
	return &discordgo.ComponentEmoji{
		Name:     e.Name,
		ID:       e.ID,
		Animated: e.Animated,
	}
}

// GetConfigEmoji returns a configured [discordgo.Emoji] for the given name.
func GetConfigEmoji(name string) (e *discordgo.Emoji) {
	override := viper.GetString("event.emoji." + name)
	if override != "" && override != name {
		return GetConfigEmoji(override)
	}
	e = &discordgo.Emoji{
		Name:     viper.GetString("event.emoji." + name + ".name"),
		ID:       viper.GetString("event.emoji." + name + ".id"),
		Animated: viper.GetBool("event.emoji." + name + ".animated"),
	}
	if e.Name == "" && e.ID == "" {
		log.Printf("Warning: tried to get emoji '%s', but its not configured or empty\n", name)
	}
	return e
}

// CompareEmoji returns true if the two emoji are the same
func CompareEmoji[E1, E2 *discordgo.Emoji | *discordgo.ComponentEmoji](e1 E1, e2 E2) bool {
	return *componentEmoji(e1) == *componentEmoji(e2)
}

// componentEmoji returns a [discordgo.ComponentEmoji] for the given [discordgo.Emoji] or [discordgo.ComponentEmoji].
func componentEmoji[E *discordgo.Emoji | *discordgo.ComponentEmoji](e E) *discordgo.ComponentEmoji {
	if ee, ok := any(e).(*discordgo.Emoji); ok {
		return &discordgo.ComponentEmoji{
			Name:     ee.Name,
			ID:       ee.ID,
			Animated: ee.Animated,
		}
	}
	if ce, ok := any(e).(*discordgo.ComponentEmoji); ok {
		return ce
	}
	panic("Given generic type is not an emoji or component emoji")
}

// MessageComplexEdit converts a [discordgo.MessageSend] to a [discordgo.MessageEdit]
func MessageComplexEdit(src *discordgo.MessageSend, channel, id string) *discordgo.MessageEdit {
	return &discordgo.MessageEdit{
		Content:         &src.Content,
		Components:      &src.Components,
		Embeds:          &src.Embeds,
		AllowedMentions: src.AllowedMentions,
		Flags:           src.Flags,
		Files:           src.Files,

		Channel: channel,
		ID:      id,
	}
}

// MessageComplexSend converts a [discordgo.MessageEdit] to a [discordgo.MessageSend]
func MessageComplexSend(src *discordgo.MessageEdit) *discordgo.MessageSend {
	return &discordgo.MessageSend{
		Content:         *src.Content,
		Components:      *src.Components,
		Embeds:          *src.Embeds,
		AllowedMentions: src.AllowedMentions,
		Flags:           src.Flags,
		Files:           src.Files,
	}

}
