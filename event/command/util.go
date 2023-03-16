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

package command

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type InteractionUtil struct {
	session     *discordgo.Session
	interaction *discordgo.InteractionCreate
	response    *discordgo.InteractionResponse
}

// Replyf formats according to a format specifier
// and prints the result as reply to the user who
// executes the command.
func (i *InteractionUtil) Replyf(format string, a ...any) {
	i.Reply(fmt.Sprintf(format, a...))
}

// Prints the given message as reply to the
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

// Replyf formats according to a format specifier
// and prints the result as emphemral reply to
// the user who executes the command.
func (i *InteractionUtil) ReplyHiddenf(format string, a ...any) {
	i.ReplyHidden(fmt.Sprintf(format, a...))
}

// Prints the given message as emphemral reply
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

func (i *InteractionUtil) ReplyError() {
	i.ReplyHidden("Somthing went wrong :(")
}

func (i *InteractionUtil) respond() {
	err := i.session.InteractionRespond(i.interaction.Interaction, i.response)
	if err != nil {
		fmt.Printf("Error while sending command response: %v", err)
	}
}
