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
	players = derangementMatch(players)

	inviteMessage := &discordgo.MessageSend{
		Embeds: make([]*discordgo.MessageEmbed, 1),
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
			}},
		},
	}

	var failedToSend string
	for _, player := range players {
		var DMChannel *discordgo.Channel
		DMChannel, err = c.Session.UserChannelCreate(player.User.ID)
		if err != nil {
			log.Printf("ERROR: could not create DM channel for user %s: %+v", player.User.ID, err)
			failedToSend += "\n- " + player.Mention()
			continue
		}

		inviteMessage.Embeds[0] = player.InviteEmbed(c.Session)
		var msg *discordgo.Message
		msg, err = c.Session.ChannelMessageSendComplex(DMChannel.ID, inviteMessage)
		if err != nil {
			log.Printf("ERROR: could not send invite: %+v", err)
			failedToSend += "\n- " + player.Mention()
			continue
		}
		player.MessageID = msg.ID
		log.Printf("Sent invite to user %s in channel %s", player.User.ID, DMChannel.ID)
	}

	if failedToSend != "" {
		c.ReplyHiddenf("Failed to send invites to:%s", failedToSend)
		return
	}

	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not save players to file: %+v", err)
		c.ReplyError()
		return
	}

	c.ReplyHidden(lang.GetDefault(tp + "msg.setup.success"))
}
