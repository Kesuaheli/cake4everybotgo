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

package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GiveawayEntry represents a giveaway entry from the database
type GiveawayEntry struct {
	// Identification of the entry
	UserID string
	// The current weight or number of tickets in this entry
	Weight int
	// The day of last entry. Useful to check when only one ticket per day is allowed.
	LastEntry time.Time
}

// ToEmbedField formats the giveaway entry to an discord message embed field.
func (e GiveawayEntry) ToEmbedField(s *discordgo.Session, totalTickets int) (f *discordgo.MessageEmbedField) {
	var name string
	if u, err := s.User(e.UserID); err != nil {
		log.Printf("Error on getting user '%s': %v", e.UserID, err)
		name = "???"
	} else {
		name = u.Username
	}

	return &discordgo.MessageEmbedField{
		Name:   name,
		Value:  fmt.Sprintf("<@%s>\n%d tickets\nChance: %.2f%%\nlast entry: <t:%d:R>", e.UserID, e.Weight, float64(e.Weight*100)/float64(totalTickets), e.LastEntry.Unix()),
		Inline: true,
	}
}

// GetGiveawayEntry gets the giveaway entry for the given user identifier, if their last entry was
// prefixed with prefix.
//
// If an error occours or it doesn't match prefix, an emtpy GiveawayEntry is returned instead.
func GetGiveawayEntry(prefix, userID string) GiveawayEntry {
	var (
		weight      int
		lastEntryID string
	)
	err := QueryRow("SELECT weight,last_entry_id FROM giveaway WHERE id=?", userID).Scan(&weight, &lastEntryID)
	if err == sql.ErrNoRows {
		return GiveawayEntry{UserID: userID, Weight: 0}
	}
	if err != nil {
		log.Printf("Database failed to get giveaway entries for '%s': %v", userID, err)
		return GiveawayEntry{}
	}

	if lastEntryID == "" {
		return GiveawayEntry{UserID: userID, Weight: weight}
	}

	dateValue, ok := strings.CutPrefix(lastEntryID, prefix+"-")
	if !ok {
		return GiveawayEntry{}
	}

	lastEntry, err := time.Parse(time.DateOnly, dateValue)
	if err != nil {
		log.Printf("could not convert last_entry_id '%s' to time: %v", lastEntryID, err)
		return GiveawayEntry{}
	}
	return GiveawayEntry{userID, weight, lastEntry}
}

// AddGiveawayWeight adds amount to the given user identifier.
//
// However if their last entry wasn't prefixed with prefix, their weight will be resetted and starts
// at amount. If you dont want it to be resetted check with GetGiveawayEntry first.
//
// If there was no error the modified entry is returned. If there was an error, an emtpy
// GiveawayEntry is returned instead.
func AddGiveawayWeight(prefix, userID string, amount int) GiveawayEntry {
	var (
		weight      int
		lastEntryID string
		new         bool
	)
	err := QueryRow("SELECT weight,last_entry_id FROM giveaway WHERE id=?", userID).Scan(&weight, &lastEntryID)
	if err == sql.ErrNoRows {
		new = true
	} else if err != nil {
		log.Printf("Database failed to get giveaway weight for '%s': %v", userID, err)
		return GiveawayEntry{}
	}

	// validate prefix
	if _, ok := strings.CutPrefix(lastEntryID, prefix+"-"); !ok {
		weight = 0
	}

	weight += amount
	dateValue := time.Now().Format(time.DateOnly)
	lastEntryID = fmt.Sprintf("%s-%s", prefix, dateValue)
	lastEntry, _ := time.Parse(time.DateOnly, dateValue)

	if new {
		_, err = Exec("INSERT INTO giveaway (id,weight,last_entry_id) VALUES (?,?,?)", userID, weight, lastEntryID)
		if err != nil {
			log.Printf("Database failed to insert giveaway for '%s': %v", userID, err)
			return GiveawayEntry{}
		}
		return GiveawayEntry{userID, weight, lastEntry}
	}
	_, err = Exec("UPDATE giveaway SET weight=?,last_entry_id=? WHERE id=?", weight, lastEntryID, userID)
	if err != nil {
		log.Printf("Database failed to update weight (new: %d) for '%s': %v", weight, userID, err)
		return GiveawayEntry{}
	}
	return GiveawayEntry{userID, weight, lastEntry}
}

// GetAllGiveawayEntries gets all giveaway entries that matches prefix.
func GetAllGiveawayEntries(prefix string) []GiveawayEntry {
	rows, err := Query("SELECT id,weight,last_entry_id FROM giveaway")
	if err != nil {
		log.Printf("ERROR: could not get entries from database: %v", err)
		return []GiveawayEntry{}
	}
	defer rows.Close()

	var entries []GiveawayEntry
	for rows.Next() {
		var (
			userID      string
			weight      int
			lastEntryID string
		)
		err = rows.Scan(&userID, &weight, &lastEntryID)
		if err != nil {
			log.Printf("Warning: could not scan variables from row")
			continue
		}

		if lastEntryID == "" {
			entries = append(entries, GiveawayEntry{UserID: userID, Weight: weight})
			continue
		}

		dateValue, ok := strings.CutPrefix(lastEntryID, prefix+"-")
		if !ok {
			continue
		}

		lastEntry, err := time.Parse(time.DateOnly, dateValue)
		if err != nil {
			log.Printf("ERROR: could not convert last_entry_id '%s' to time: %v", lastEntryID, err)
			continue
		}
		entries = append(entries, GiveawayEntry{userID, weight, lastEntry})
	}
	return entries
}

func DrawGiveawayWinner(a []GiveawayEntry) (winner GiveawayEntry, totalTickets int) {
	var entries []GiveawayEntry
	for _, e := range a {
		for i := 0; i < e.Weight; i++ {
			entries = append(entries, e)
		}
	}
	totalTickets = len(entries)
	if totalTickets == 0 {
		return GiveawayEntry{}, 0
	}

	rand.Shuffle(len(entries), func(i, j int) {
		entries[i], entries[j] = entries[j], entries[i]
	})
	return entries[rand.Intn(totalTickets-1)], totalTickets
}
