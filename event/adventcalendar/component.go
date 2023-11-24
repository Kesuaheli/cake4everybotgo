package adventcalendar

import (
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

type Component struct {
	AdventCalendar
}

func (c Component) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {

	c.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	c.member = i.Member
	c.user = i.User
	if i.Member != nil {
		c.user = i.Member.User
	} else if i.User != nil {
		c.member = &discordgo.Member{User: i.User}
	}
}

func (Component) ID() string {
	return "adventcalendar"
}
