package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleSetup(ids []string) {
	switch util.ShiftL(ids) {
	case "invite":
		c.handleSetupInvite()
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}
}

func (c Component) handleSetupInvite() {
	players, err := c.getPlayers()
	if err != nil {
		log.Printf("ERROR: could not get players: %+v", err)
		c.ReplyError()
		return
	}
	c.ReplyDeferedHidden()

	inviteMessage := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Title:  lang.GetDefault(tp + "msg.invite.title"),
			Fields: []*discordgo.MessageEmbedField{},
		}},
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{Components: []discordgo.MessageComponent{
				util.CreateButtonComponent(
					fmt.Sprintf("secretsanta.invite.show_match.%s", c.Interaction.GuildID),
					lang.GetDefault(tp+"msg.invite.button.show_match"),
					discordgo.PrimaryButton,
					util.GetConfigComponentEmoji("secretsanta.invite.show_match"),
				),
				util.CreateButtonComponent(
					fmt.Sprintf("secretsanta.invite.set_address.%s", c.Interaction.GuildID),
					lang.GetDefault(tp+"msg.invite.button.set_address"),
					discordgo.SecondaryButton,
					util.GetConfigComponentEmoji("secretsanta.invite.set_address"),
				),
				util.CreateButtonComponent(
					fmt.Sprintf("secretsanta.invite.show_address.%s", c.Interaction.GuildID),
					lang.GetDefault(tp+"msg.invite.button.show_address"),
					discordgo.SecondaryButton,
					util.GetConfigComponentEmoji("secretsanta.invite.show_address"),
				),
			}},
		},
	}

	var errCount int
	for _, player := range players {
		var DMChannel *discordgo.Channel
		DMChannel, err = c.Session.UserChannelCreate(player.User.ID)
		if err != nil {
			log.Printf("ERROR: could not create DM channel for user %s: %+v", player.User.ID, err)
			errCount++
			continue
		}

		_, err = c.Session.ChannelMessageSendComplex(DMChannel.ID, inviteMessage)
		if err != nil {
			log.Printf("ERROR: could not send invite: %+v", err)
			errCount++
			continue
		}
		log.Printf("Sent invite to user %s in channel %s", player.User.ID, DMChannel.ID)
	}

	if errCount > 0 {
		c.ReplyHiddenf("Failed to send %d invites!", errCount)
		return
	}

	c.ReplyHidden(lang.GetDefault(tp + "msg.setup.success"))
}
