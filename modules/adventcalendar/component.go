package adventcalendar

import (
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

// The Component of the advent calendar package.
type Component struct {
	adventcalendarBase
}

// Handle handles the functionality of a component.
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

// ID returns the custom ID of the modal to identify the module
func (Component) ID() string {
	return "adventcalendar"
}
