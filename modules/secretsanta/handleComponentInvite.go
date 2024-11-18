package secretsanta

import (
	"cake4everybot/data/lang"
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

	var player *player
	for _, p := range players {
		if p.User.ID == c.Interaction.User.ID {
			player = p
		}
	}
	if player == nil {
		log.Printf("ERROR: could not find player %s in guild %s: %+v", c.Interaction.User.ID, c.Interaction.GuildID, c.Interaction.User.ID)
		c.ReplyError()
		return
	}

	c.ReplyModal("secretsanta.set_address."+c.Interaction.GuildID, lang.GetDefault(tp+"msg.invite.modal.set_address.title"), discordgo.ActionsRow{Components: []discordgo.MessageComponent{
		discordgo.TextInput{
			CustomID:    "address",
			Label:       lang.GetDefault(tp + "msg.invite.modal.set_address.label"),
			Style:       discordgo.TextInputParagraph,
			Placeholder: lang.GetDefault(tp + "msg.invite.modal.set_address.placeholder"),
			Value:       player.Address,
			Required:    true,
		},
	}})
}

func (c Component) handleInviteShowAddress(ids []string) {

}
