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
