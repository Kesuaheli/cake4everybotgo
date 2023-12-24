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
	"cake4everybot/event/command"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

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
		username = member.Nick
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

// SetEmbedFooter takes a pointer to an embeds and sets the standard
// footer with the given name.
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

// AddReplyHiddenField appends the standard field for ephemral
// embeds to the existing fields of the given embed.
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
	if command.CommandMap[cBase] == nil {
		return ""
	}
	cID := command.CommandMap[cBase].GetID()

	var cSub string
	for _, sub := range subcommand {
		cSub = cSub + " " + lang.GetDefault(sub)
	}

	return fmt.Sprintf("</%s%s:%s>", cBase, cSub, cID)
}

// CreateButtonComponent returns a simple button component
// with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to identify the button when pressed (automatically prefixed)
//	style              // Style of the button (see https://discord.com/developers/docs/interactions/message-components#button-object-button-styles)
//	Optional: emoji    // An emoji to put in the label, can be empty
func CreateButtonComponent(tp, component_group_id, id string, style discordgo.ButtonStyle, emoji discordgo.ComponentEmoji) discordgo.Button {
	return discordgo.Button{
		CustomID: lang.GetDefault(tp+"base") + "." + id,
		Label:    lang.Get(tp+"component."+component_group_id+".button."+id+".label", lang.FallbackLang()),
		Style:    style,
		Emoji:    emoji,
	}
}

// CreateURLButtonComponent returns a URL button component
// with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to generate labels
//	url                // The link to open when clicked
//	Optional: emoji    // An emoji to put in the label, can be empty
func CreateURLButtonComponent(tp, component_group_id, id, url string, emoji discordgo.ComponentEmoji) discordgo.Button {
	return discordgo.Button{
		Label: lang.Get(tp+"component."+component_group_id+".button."+id+".label", lang.FallbackLang()),
		Style: discordgo.LinkButton,
		Emoji: emoji,
		URL:   url,
	}
}

// CreateTextInputComponent returns a text input form for
// modals with the specified configurations.
// Params:
//
//	tp                 // The translation prefix used to generate ids and labels
//	component_group_id // A group id to generate labels
//	id                 // Custom id to identify the input field after submitting
//	style              // Single or multi line
//	required           // If this has to be not empty
//	minLength          // Minimum number of characters that has to be entered
//	maxLength          // Maximum number of characters that are able to be entered
func CreateTextInputComponent(tp, component_group_id, id string, style discordgo.TextInputStyle, requred bool, minLength, maxLength int) discordgo.TextInput {
	return discordgo.TextInput{
		CustomID:    id,
		Label:       lang.Get(tp+"component."+component_group_id+".text_input."+id+".label", lang.FallbackLang()),
		Style:       style,
		Placeholder: lang.Get(tp+"component."+component_group_id+".text_input."+id+".placeholder", lang.FallbackLang()),
		Value:       lang.Get(tp+"component."+component_group_id+".text_input."+id+".value", lang.FallbackLang()),
		Required:    requred,
		MinLength:   minLength,
		MaxLength:   maxLength,
	}
}

func CreateMultiselectComponent() discordgo.SelectMenu {
	return discordgo.SelectMenu{}
}
