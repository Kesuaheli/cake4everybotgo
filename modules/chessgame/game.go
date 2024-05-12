package chessgame

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/notnil/chess"
)

const (
	numOfSquaresInRow = 8

	iconCorner = ":blue_square:"

	iconBBBS = "<:bbbs:1239220881255960648>"
	iconBBWS = "<:bbws:1239220882778488882>"
	iconBKBS = "<:bkbs:1239220884334706780>"
	iconBKWS = "<:bkws:1239220886083731467>"
	iconBNBS = "<:bnbs:1239220887862116492>"
	iconBNWS = "<:bnws:1239220889271144519>"
	iconBPBS = "<:bpbs:1239220890676494426>"
	iconBPWS = "<:bpws:1239220892073197648>"
	iconBQBS = "<:bqbs:1239220893524430939>"
	iconBQWS = "<:bqws:1239221244285550642>"
	iconBRBS = "<:brbs:1239220896816697396>"
	iconBRWS = "<:brws:1239220898871902328>"
	iconBS   = "<:bs:1239221245661151263>"

	iconWBBS = "<:wbbs:1239220903532040313>"
	iconWBWS = "<:wbws:1239220906744745984>"
	iconWKBS = "<:wkbs:1239221247078961263>"
	iconWKWS = "<:wkws:1239220910226145351>"
	iconWNBS = "<:wnbs:1239221248572264549>"
	iconWNWS = "<:wnws:1239220914738954292>"
	iconWPBS = "<:wpbs:1239220917771440129>"
	iconWPWS = "<:wpws:1239221250132283423>"
	iconWQBS = "<:wqbs:1239220921328205924>"
	iconWQWS = "<:wqws:1239220923811496006>"
	iconWRBS = "<:wrbs:1239221251541700699>"
	iconWRWS = "<:wrws:1239220927502352495>"
	iconWS   = "<:ws:1239221242821742652>"
)

var (
	iconFiles = [numOfSquaresInRow]string{
		":regional_indicator_a:",
		":regional_indicator_b:",
		":regional_indicator_c:",
		":regional_indicator_d:",
		":regional_indicator_e:",
		":regional_indicator_f:",
		":regional_indicator_g:",
		":regional_indicator_h:",
	}
	iconRanks = [numOfSquaresInRow]string{
		":one:",
		":two:",
		":three:",
		":four:",
		":five:",
		":six:",
		":seven:",
		":eight:",
	}
)

type Game struct {
	*chess.Game
	PlayerWhite *discordgo.Member
	PlayerBlack *discordgo.Member
}

// NewGame creates a new game of chess for the two given player
func NewGame(white *discordgo.Member, black *discordgo.Member) *Game {
	return &Game{
		Game:        chess.NewGame(),
		PlayerWhite: white,
		PlayerBlack: black,
	}
}

// Turn returns the player whoose turn it is
func (g *Game) Turn() *discordgo.Member {
	switch g.Position().Turn() {
	case chess.White:
		return g.PlayerWhite
	case chess.Black:
		return g.PlayerBlack
	default:
		panic(fmt.Sprintf("The game's current turn is not valid: %s!", g.Position().Turn()))
	}
}

// Display returns a graphical representation of the current game using emojis
func (g *Game) Display() string {
	board := strings.Builder{}

	//board.WriteString(iconCorner + strings.Join(iconFiles[:], "") + iconCorner)
	//board.WriteByte('\n')

	for r := numOfSquaresInRow - 1; r >= 0; r-- {
		board.WriteString(iconRanks[r])
		for f := 0; f < numOfSquaresInRow; f++ {
			piece := g.Position().Board().Piece(chess.NewSquare(chess.File(f), chess.Rank(r)))
			bgColor := chess.Color((r+f)%2 + 1)

			icon := PieceToIcon(piece, bgColor)
			board.WriteString(icon)
		}
		//board.WriteString(iconRanks[r])
		board.WriteByte('\n')
	}

	board.WriteString(iconCorner + strings.Join(iconFiles[:], ""))

	return board.String()
}

// PieceToIcon returns the emoij string for a given piece
func PieceToIcon(piece chess.Piece, bg chess.Color) string {
	if bg != chess.Black && bg != chess.White {
		panic(fmt.Sprintf("pice has no background color: %s (%d)", bg, bg))
	}

	switch piece {
	case chess.BlackBishop:
		if bg == chess.Black {
			return iconBBBS
		}
		return iconBBWS
	case chess.BlackKing:
		if bg == chess.Black {
			return iconBKBS
		}
		return iconBKWS
	case chess.BlackKnight:
		if bg == chess.Black {
			return iconBNBS
		}
		return iconBNWS
	case chess.BlackPawn:
		if bg == chess.Black {
			return iconBPBS
		}
		return iconBPWS
	case chess.BlackQueen:
		if bg == chess.Black {
			return iconBQBS
		}
		return iconBQWS
	case chess.BlackRook:
		if bg == chess.Black {
			return iconBRBS
		}
		return iconBRWS
	case chess.WhiteBishop:
		if bg == chess.Black {
			return iconWBBS
		}
		return iconWBWS
	case chess.WhiteKing:
		if bg == chess.Black {
			return iconWKBS
		}
		return iconWKWS
	case chess.WhiteKnight:
		if bg == chess.Black {
			return iconWNBS
		}
		return iconWNWS
	case chess.WhitePawn:
		if bg == chess.Black {
			return iconWPBS
		}
		return iconWPWS
	case chess.WhiteQueen:
		if bg == chess.Black {
			return iconWQBS
		}
		return iconWQWS
	case chess.WhiteRook:
		if bg == chess.Black {
			return iconWRBS
		}
		return iconWRWS
	case chess.NoPiece:
		if bg == chess.Black {
			return iconBS
		}
		return iconWS
	default:
		panic(fmt.Sprintf("piece has no icon: %s (%d)", piece, piece))
	}
}
