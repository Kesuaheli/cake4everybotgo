package secretsanta

import (
	"cake4everybot/util"
	"encoding/json"
	"fmt"
	logger "log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => adventcalendar
	tp = "discord.command.secretsanta."
)

var log = logger.New(logger.Writer(), "[SecretSanta] ", logger.LstdFlags|logger.Lmsgprefix)

type secretSantaBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

// getPlayers returns the list of players for the current guild. If it is the first time, it loads
// the players from the file or creates an empty file.
func (ssb secretSantaBase) getPlayers() ([]*player, error) {
	if allPlayers != nil {
		return allPlayers[ssb.Interaction.GuildID], nil
	}

	log.Println("First time getting players. Loading from file...")
	playersPath := viper.GetString("event.secretsanta.players")
	playersData, err := os.ReadFile(playersPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("read players file: %v", err)
		}
		allPlayers = make(AllPlayers)
		playersData, err = json.Marshal(allPlayers)
		if err != nil {
			return nil, fmt.Errorf("marshal players file: %v", err)
		}
		err = os.WriteFile(playersPath, playersData, 0644)
		if err != nil {
			return nil, fmt.Errorf("write players file: %v", err)
		}
		log.Printf("Created players file: %s\n", playersPath)
		return []*player{}, nil
	}
	allPlayersUnresolved := AllPlayersUnresolved{}
	err = json.Unmarshal(playersData, &allPlayersUnresolved)
	if err != nil {
		return nil, fmt.Errorf("unmarshal players file: %v", err)
	}
	err = allPlayersUnresolved.Resolve(ssb.Session)
	if err != nil {
		return nil, fmt.Errorf("resolve players file: %v", err)
	}
	log.Printf("Got %d guilds from file", len(allPlayers))

	return allPlayers[ssb.Interaction.GuildID], nil
}

// setPlayers sets the players for the current guild.
func (ssb secretSantaBase) setPlayers(players []*player) (err error) {
	if _, err = ssb.getPlayers(); err != nil {
		return err
	}

	allPlayers[ssb.Interaction.GuildID] = players
	playersData, err := json.Marshal(allPlayers)
	if err != nil {
		return fmt.Errorf("marshal players file: %v", err)
	}
	err = os.WriteFile(viper.GetString("event.secretsanta.players"), playersData, 0644)
	if err != nil {
		return fmt.Errorf("write players file: %v", err)
	}
	return nil
}

// player is a player in the secret santa game
type player struct {
	*discordgo.Member

	// Match is the matched player
	Match *player
	// Address is the address of the player
	Address string
}

type playerUnresolved struct {
	ID      string `json:"id"`
	MatchID string `json:"match"`
	Address string `json:"address"`
}

// AllPlayers is a map from guild ID to a list of players
type AllPlayers map[string][]*player

// allPlayers is the current state of all players.
// See [AllPlayers]
var allPlayers AllPlayers

// MarshalJSON implements json.Marshaler
func (allPlayers AllPlayers) MarshalJSON() ([]byte, error) {
	m := make(AllPlayersUnresolved)
	for guildID, players := range allPlayers {
		for _, player := range players {
			m[guildID] = append(m[guildID], &playerUnresolved{
				ID:      player.User.ID,
				MatchID: player.Match.User.ID,
				Address: player.Address,
			})
		}

	}
	return json.Marshal(m)
}

// AllPlayersUnresolved is a map from guild ID to a list of unresolved players.
// Unresolved players have no member but only an ID
type AllPlayersUnresolved map[string][]*playerUnresolved

// Resolve resolves allPlayersUnresolved into allPlayers
func (allPlayersUnresolved AllPlayersUnresolved) Resolve(s *discordgo.Session) (err error) {
	allPlayers = make(AllPlayers)
	for guildID, unresolvedPlayers := range allPlayersUnresolved {
		resolvedPlayers := map[string]*player{}
		for _, up := range unresolvedPlayers {
			member, err := s.GuildMember(guildID, up.ID)
			if err != nil {
				return fmt.Errorf("failed to get guild member %s/%s: %v", guildID, up.ID, err)
			}
			resolvedPlayers[up.ID] = &player{
				Member:  member,
				Match:   resolvedPlayers[up.MatchID],
				Address: up.Address,
			}
		}
		for _, rp := range resolvedPlayers {
			if rp.Match != nil {
				continue
			}
			rp.Match = resolvedPlayers[rp.User.ID]

			allPlayers[guildID] = append(allPlayers[guildID], rp)
		}
	}
	return nil
}

// derangementMatch matches the players in a way that no one gets matched to themselves.
func derangementMatch(players []*player) []*player {
	n := len(players)
	players2 := make([]*player, n)
	copy(players2, players)

	for i := 0; i < n-1; i++ {
		j := i + rand.Intn(n-i-1) + 1
		players2[i], players2[j] = players2[j], players2[i]
	}

	for i, item := range players {
		item.Match = players2[i]
	}

	return players
}
