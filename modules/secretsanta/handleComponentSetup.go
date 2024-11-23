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
	players, err = derangementMatch(players)
	if err != nil {
		log.Printf("ERROR: could not match players: %+v", err)
		c.ReplySimpleEmbed(0xFF0000, lang.GetDefault(tp+"msg.setup.match_error"))
		return
	}

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
	}

	if failedToSend != "" {
		c.ReplyHiddenSimpleEmbedf(0xFF0000, lang.GetDefault(tp+"msg.setup.invite.error"), failedToSend)
		return
	}

	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not save players to file: %+v", err)
		c.ReplyError()
		return
	}

	c.ReplyHiddenSimpleEmbed(0x690042, lang.GetDefault(tp+"msg.setup.success"))
}
