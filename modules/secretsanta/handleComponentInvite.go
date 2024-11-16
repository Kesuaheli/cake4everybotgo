package secretsanta

import (
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

func (c Component) handleInvite(s *discordgo.Session, ids []string) {
	switch util.ShiftL(ids) {
	case "show_match":
		c.handleInviteShowMatch(s, ids)
		return
	case "set_address":
		c.handleInviteSetAddress(s, ids)
		return
	case "show_address":
		c.handleInviteShowAddress(s, ids)
		return
	default:
		log.Printf("Unknown component interaction ID: %s", c.data.CustomID)
	}
}

func (c Component) handleInviteShowMatch(s *discordgo.Session, ids []string) {
	guildID := util.ShiftL(ids)
	_, err := c.getPlayers()
	if err != nil {
		log.Printf("ERROR: could not get players: %+v", err)
		c.ReplyError()
		return
	}
	players := allPlayers[guildID]
	if len(players) == 0 {
		log.Printf("ERROR: no players in guild %s", guildID)
		c.ReplyError()
		return
	}
	var player *player
	for _, p := range players {
		if p.User.ID == c.Interaction.User.ID {
			player = p
		}
	}

	if player == nil {
		log.Printf("ERROR: could not find player %s in guild %s: %+v", c.Interaction.User.ID, guildID, c.Interaction.User.ID)
		c.ReplyError()
		return
	}

	e := util.AuthoredEmbed(s, player.Match.Member, tp+"display")

	util.SetEmbedFooter(s, tp+"display", e)
	c.ReplyHiddenEmbed(e)
}

func (c Component) handleInviteSetAddress(s *discordgo.Session, ids []string) {

}

func (c Component) handleInviteShowAddress(s *discordgo.Session, ids []string) {

}
