package chessgame

import (
	"cake4everybot/data/lang"
	"cake4everybot/util"
	logger "log"

	"github.com/bwmarrin/discordgo"
)

var (
	log = logger.New(logger.Writer(), "[Chess] ", logger.LstdFlags|logger.Lmsgprefix)
)

// The Chat (slash) command of the chess package.
type Chat struct {
	chessBase

	ID string
}

// AppCmd (ApplicationCommand) returns the definition of the chat command
func (cmd Chat) AppCmd() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:                     lang.GetDefault(tp + "base"),
		NameLocalizations:        util.TranslateLocalization(tp + "base"),
		Description:              lang.GetDefault(tp + "base.description"),
		DescriptionLocalizations: util.TranslateLocalization(tp + "base.description"),
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:                     discordgo.ApplicationCommandOptionUser,
				Name:                     lang.GetDefault(tp + "option.opponent"),
				NameLocalizations:        *util.TranslateLocalization(tp + "option.opponent"),
				Description:              lang.GetDefault(tp + "option.opponent.description"),
				DescriptionLocalizations: *util.TranslateLocalization(tp + "option.opponent.description"),
				Required:                 true,
			},
		},
	}
}

// Handle handles the functionality of a command
func (cmd Chat) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	cmd.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
	cmd.member = i.Member
	cmd.user = i.User
	if i.Member != nil {
		cmd.user = i.Member.User
	} else if i.User != nil {
		cmd.member = &discordgo.Member{User: i.User}
	}

	var opponent *discordgo.Member
	for _, option := range i.ApplicationCommandData().Options {
		switch option.Name {
		case lang.GetDefault(tp + "option.opponent"):
			opponentUser := option.UserValue(s)
			if cmd.user.ID == opponentUser.ID {
				cmd.ReplyHidden(lang.GetDefault(tp + "msg.opponent.no_self"))
				return
			}
			opponent = &discordgo.Member{User: opponentUser}
		}
	}

	game := NewGame(opponent, cmd.member)
	log.Printf("new game:\n%s", game.Position().Board().Draw())
	board := game.Display()
	log.Printf("board length: %d", len(board))
	cmd.Reply(board)
}

// SetID sets the registered command ID for internal uses after uploading to discord
func (cmd *Chat) SetID(id string) {
	cmd.ID = id
}

// GetID gets the registered command ID
func (cmd Chat) GetID() string {
	return cmd.ID
}
