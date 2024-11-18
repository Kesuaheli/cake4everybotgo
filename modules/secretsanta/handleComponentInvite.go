package secretsanta

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	"fmt"

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
	e.Title = fmt.Sprintf(lang.GetDefault(tp+"msg.invite.show_match.title"), player.Match.Member.DisplayName())
	e.Description = lang.GetDefault(tp + "msg.invite.show_match.description")
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:  lang.GetDefault(tp + "msg.invite.show_match.address"),
		Value: fmt.Sprintf("```\n%s\n```", player.Match.Address),
	})
	if player.Match.Address == "" {
		log.Printf("%s has no address set: %+v", player.Match.Member.DisplayName(), player.Match)
		e.Fields[0].Value = lang.GetDefault(tp + "msg.invite.show_match.address_not_set")
	}

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
