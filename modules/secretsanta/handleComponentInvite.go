package secretsanta

import (
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleInvite(ids []string) {
	switch util.ShiftL(ids) {
	case "show_match":
		c.handleInviteShowMatch(ids)
		return
	case "set_address":
		c.handleInviteSetAddress(ids)
		return
	case "show_address":
		c.handleInviteShowAddress(ids)
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}
}

func (c Component) handleInviteShowMatch(ids []string) {
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

	e := util.AuthoredEmbed(c.Session, player.Match.Member, tp+"display")

	util.SetEmbedFooter(c.Session, tp+"display", e)
	c.ReplyHiddenEmbed(e)
}

func (c Component) handleInviteSetAddress(ids []string) {

}

func (c Component) handleInviteShowAddress(ids []string) {

}
