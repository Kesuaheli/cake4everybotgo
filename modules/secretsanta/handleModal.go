package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleModalSetAddress(ids []string) {
	c.Interaction.GuildID = util.ShiftL(ids)
	players, err := c.getPlayers()
	if err != nil {
		log.Printf("ERROR: could not get players: %+v", err)
		c.ReplyError()
		return
	}
	if len(players) == 0 {
		log.Printf("ERROR: no players in guild %s", c.Interaction.GuildID)
		c.ReplyError()
		return
	}
	player, ok := players[c.Interaction.User.ID]
	if !ok {
		log.Printf("ERROR: could not find player %s in guild %s: %+v", c.Interaction.User.ID, c.Interaction.GuildID, c.Interaction.User.ID)
		c.ReplyError()
		return
	}

	addressFiled := c.modal.Components[0].(*discordgo.ActionsRow).Components[0].(*discordgo.TextInput)
	if addressFiled.Value == player.Address {
		c.ReplyHidden(lang.GetDefault(tp + "msg.invite.set_address.not_changed"))
		return
	}

	player.Address = addressFiled.Value
	err = c.setPlayers(players)
	if err != nil {
		log.Printf("ERROR: could not set players: %+v", err)
		c.ReplyError()
		return
	}

	e := &discordgo.MessageEmbed{
		Color: 0x00FF00,
		Fields: []*discordgo.MessageEmbedField{{
			Name:  lang.GetDefault(tp + "msg.invite.set_address.changed"),
			Value: fmt.Sprintf("```\n%s\n```", player.Address),
		}},
	}

	util.SetEmbedFooter(c.Session, tp+"display", e)
	c.ReplyHiddenEmbed(e)
}
