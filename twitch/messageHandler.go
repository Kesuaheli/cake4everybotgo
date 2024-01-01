// Copyright 2023 Kesuaheli
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package twitch

import (
	"cake4everybot/data/lang"
	"cake4everybot/database"
	"cake4everybot/tools/streamelements"
	logger "log"
	"strings"

	"github.com/kesuaheli/twitchgo"
	"github.com/spf13/viper"
)

const tp string = "twitch.command."

var log logger.Logger = *logger.New(logger.Writer(), "[Twitch] ", logger.LstdFlags|logger.Lmsgprefix)
var se *streamelements.Streamelements

// MessageHandler handles new messages from the twitch chat(s). It will be called on every new
// message.
func MessageHandler(t *twitchgo.Twitch, channel string, user *twitchgo.User, message string) {
	log.Printf("<%s@%s> %s", user.Nickname, channel, message)
}

// HandleCmdJoin is the handler for a command in a twitch chat. This handler buys a giveaway ticket
// and removes the configured cost amount for a ticket.
func HandleCmdJoin(t *twitchgo.Twitch, channel string, user *twitchgo.User, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	const tp = tp + "join."

	log.Printf("[%s@%s] executed join command with %d args: %v", user.Nickname, channel, len(args), args)
	seChannel, err := se.GetChannel(channel)
	if err != nil {
		log.Printf("Error getting streamelements channel '%s': %v", channel, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	sePoints, err := se.GetPoints(seChannel.ID, user.Nickname)
	if err != nil {
		log.Printf("Error getting streamelements points '%s(%s)/%s' : %v", seChannel.ID, channel, user.Nickname, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	joinCost := viper.GetInt("event.twitch_giveaway.ticket_cost")
	if sePoints.Points < joinCost {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.too_few_points"), user.Nickname, sePoints.Points, joinCost-sePoints.Points, joinCost)
		return
	}
	entry := database.AddGiveawayWeight("tw11", user.Nickname, 1)
	if entry.UserID == "" {
		log.Println("Error getting database giveaway entry", seChannel.ID, channel, user.Nickname, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	err = se.AddPoints(seChannel.ID, user.Nickname, -joinCost)
	if err != nil {
		log.Printf("Error adding points for '%s(%s)/%s/-%d': %v", seChannel.ID, channel, user.Nickname, joinCost, err)
		t.SendMessage(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	t.SendMessagef(channel, lang.GetDefault(tp+"msg.success"), user.Nickname, joinCost, entry.Weight, sePoints.Points-joinCost)
}

// HandleCmdDraw is the handler for the draw command in a twitch chat. This handler selects a random
// winner and removes their tickets.
func HandleCmdDraw(t *twitchgo.Twitch, channel string, user *twitchgo.User, args []string) {
	channel, _ = strings.CutPrefix(channel, "#")
	const tp = tp + "draw."

	//only accept broadcaster
	if channel != user.Nickname {
		return
	}

	p, err := database.NewGiveawayPrize(viper.GetString("event.twitch_giveaway.prizes"))
	if err != nil {
		log.Printf("Error reading prizes file: %v", err)
		t.SendMessagef(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
	prize, ok := p.GetNextPrize()
	if !ok {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.no_prizes"), user.Nickname)
		return
	}

	winner, totalTickets := database.DrawGiveawayWinner(database.GetAllGiveawayEntries("tw11"))
	if totalTickets == 0 {
		t.SendMessagef(channel, lang.GetDefault(tp+"msg.no_entries"), user.Nickname)
		return
	}

	t.SendMessagef(channel, lang.GetDefault(tp+"msg.winner"), winner.UserID, prize.Name, winner.Weight, float64(winner.Weight*100)/float64(totalTickets))

	err = database.DeleteGiveawayEntry(winner.UserID)
	if err != nil {
		log.Printf("Error deleting database giveaway entry: %v", err)
		t.SendMessagef(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}

	prize.Winner = winner.UserID
	err = p.SaveFile()
	if err != nil {
		log.Printf("Error saving prizes file: %v", err)
		t.SendMessagef(channel, lang.GetDefault("twitch.command.generic.error"))
		return
	}
}
