package component

import "github.com/bwmarrin/discordgo"

// Component is an interface wrapper for all message components.
type Component interface {
	// Function of a component.
	// All things that should happen after submitting or pressing a button.
	Handle(*discordgo.Session, *discordgo.InteractionCreate)

	// Custom ID of the modal to identify the module
	ID() string
}

// ComponentMap holds all active components. It maps them from a unique string identifier to the
// corresponding Component.
var ComponentMap = make(map[string]Component)
