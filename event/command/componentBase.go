package command

import "github.com/bwmarrin/discordgo"

type Component interface {
	// Function of a component.
	// All things that should happen after submitting or
	// pressing a button.
	Handle(*discordgo.Session, *discordgo.InteractionCreate)

	// Custom ID of the modal to identify the module
	ID() string
}

var ComponentMap = make(map[string]Component)
