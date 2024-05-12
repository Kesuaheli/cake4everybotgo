package chessgame

import (
	"cake4everybot/util"

	"github.com/bwmarrin/discordgo"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => chess
	tp = "discord.command.chess."
)

type chessBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}
