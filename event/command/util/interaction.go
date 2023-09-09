// Copyright 2022-2023 Kesuaheli
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
	"cake4everybot/data/lang"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// InteractionUtil is a helper for discords application interactions.
// It add useful methods for simpler and faster coding.
type InteractionUtil struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	response    *discordgo.InteractionResponse
}

// Replyf formats according to a format specifier
// and prints the result as reply to the user who
// executes the command.
func (i *InteractionUtil) Replyf(format string, a ...any) {
	i.Reply(fmt.Sprintf(format, a...))
}

// Reply prints the given message as reply to the
// user who executes the command.
func (i *InteractionUtil) Reply(message string) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
		},
	}
	i.respond()
}

// ReplyEmbed prints the given embeds as reply to the
// user who executes the command.
func (i *InteractionUtil) ReplyEmbed(embeds ...*discordgo.MessageEmbed) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
		},
	}
	i.respond()
}

// ReplyComponents sends a message or embed along with the
// provided message components.
func (i *InteractionUtil) ReplyComponents(message string, components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Embeds:     embeds,
			Components: components,
		},
	}
	i.respond()
}

// ReplyHiddenf formats according to a format specifier
// and prints the result as ephemral reply to
// the user who executes the command.
func (i *InteractionUtil) ReplyHiddenf(format string, a ...any) {
	i.ReplyHidden(fmt.Sprintf(format, a...))
}

// ReplyHidden prints the given message as ephemral reply
// to the user who executes the command.
func (i *InteractionUtil) ReplyHidden(message string) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	}
	i.respond()
}

// ReplyHiddenEmbed prints the given embeds as ephemral reply to the user who
// executes the command. Automatically append "hidden reply note" to
// last embed if hiddenSelf is set tot true. See AddReplyHiddenField() for more.
func (i *InteractionUtil) ReplyHiddenEmbed(hiddenSelf bool, embeds ...*discordgo.MessageEmbed) {
	l := len(embeds)
	if l == 0 {
		return
	}
	if hiddenSelf {
		AddReplyHiddenField(embeds[l-1])
	}

	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: embeds,
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	}
	i.respond()
}

// ReplyHiddenComponents sends an ephemeral message or
// embed along with the provided message components.
func (i *InteractionUtil) ReplyHiddenComponents(message string, components []discordgo.MessageComponent, embeds ...*discordgo.MessageEmbed) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content:    message,
			Embeds:     embeds,
			Flags:      discordgo.MessageFlagsEphemeral,
			Components: components,
		},
	}
	i.respond()
}

// ReplyAutocomplete returns the given choices to
// the user. When this is called on an interaction
// type outside form an applicationCommandAutocomplete
// nothing will happen.
func (i *InteractionUtil) ReplyAutocomplete(choices []*discordgo.ApplicationCommandOptionChoice) {
	if i.Interaction.Type != discordgo.InteractionApplicationCommandAutocomplete {
		return
	}

	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	}
	i.respond()
}

// ReplyError sends a simple message to the user to indicate, that
// something failed or unexpected happened during the execution of
// the interaction.
func (i *InteractionUtil) ReplyError() {
	i.ReplyHidden("Somthing went wrong :(")
}

// ReplyModal displays a modal (popup) with the specified components to the user.
//
// Params:
//
//	tp // The translation prefix of the command
//	id // To identify the modal when parsing the interaction event
//	components // One or more message components to display in this modal
func (i *InteractionUtil) ReplyModal(tp, id string, components ...discordgo.MessageComponent) {
	i.response = &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID:   lang.GetDefault(tp+"base") + "." + id,
			Title:      lang.Get(tp+"modal."+id+".title", lang.FallbackLang()),
			Components: components,
		},
	}
	i.respond()
}

func (i *InteractionUtil) respond() {
	err := i.Session.InteractionRespond(i.Interaction.Interaction, i.response)
	if err != nil {
		log.Printf("Error while sending command response: %v\n", err)
	}
}
