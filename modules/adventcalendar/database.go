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

package adventcalendar

import (
	"cake4everybot/database"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func getEntry(userID string) giveawayEntry {
	var (
		weight      int
		lastEntryID string
	)
	err := database.QueryRow("SELECT weight,last_entry_id FROM giveaway WHERE id=?", userID).Scan(&weight, &lastEntryID)
	if err == sql.ErrNoRows {
		return giveawayEntry{userID: userID, weight: 0}
	}
	if err != nil {
		log.Printf("Database failed to get giveaway entries for '%s': %+v", userID, err)
		return giveawayEntry{}
	}

	if lastEntryID == "" {
		return giveawayEntry{userID: userID, weight: weight}
	}

	dateValue, ok := strings.CutPrefix(lastEntryID, "xmas-")
	if !ok {
		return giveawayEntry{}
	}

	lastEntry, err := time.Parse(time.DateOnly, dateValue)
	if err != nil {
		log.Printf("could not convert last_entry_id '%s' to time: %+v", lastEntryID, err)
		return giveawayEntry{}
	}
	return giveawayEntry{userID, weight, lastEntry}
}

func addGiveawayWeight(userID string, amount int) giveawayEntry {
	var weight int
	var new bool
	err := database.QueryRow("SELECT weight FROM giveaway WHERE id=?", userID).Scan(&weight)
	if err == sql.ErrNoRows {
		new = true
	} else if err != nil {
		log.Printf("Database failed to get giveaway weight for '%s': %+v", userID, err)
		return giveawayEntry{}
	}

	weight += amount
	dateValue := time.Now().Format(time.DateOnly)
	lastEntryID := fmt.Sprintf("xmas-%s", dateValue)
	lastEntry, _ := time.Parse(time.DateOnly, dateValue)

	if new {
		_, err = database.Exec("INSERT INTO giveaway (id,weight,last_entry_id) VALUES (?,?,?)", userID, weight, lastEntryID)
		if err != nil {
			log.Printf("Database failed to insert giveaway for '%s': %+v", userID, err)
			return giveawayEntry{}
		}
		return giveawayEntry{userID, weight, lastEntry}
	}
	_, err = database.Exec("UPDATE giveaway SET weight=?,last_entry_id=? WHERE id=?", weight, lastEntryID, userID)
	if err != nil {
		log.Printf("Database failed to update weight (new: %d) for '%s': %+v", weight, userID, err)
		return giveawayEntry{}
	}
	return giveawayEntry{userID, weight, lastEntry}
}

func getGetAllEntries() []giveawayEntry {
	rows, err := database.Query("SELECT id,weight,last_entry_id FROM giveaway")
	if err != nil {
		log.Printf("ERROR: could not get entries from database: %+v", err)
		return []giveawayEntry{}
	}
	defer rows.Close()

	var entries []giveawayEntry
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
			entries = append(entries, giveawayEntry{userID: userID, weight: weight})
			continue
		}

		dateValue, ok := strings.CutPrefix(lastEntryID, "xmas-")
		if !ok {
			continue
		}

		lastEntry, err := time.Parse(time.DateOnly, dateValue)
		if err != nil {
			log.Printf("ERROR: could not convert last_entry_id '%s' to time: %+v", lastEntryID, err)
			continue
		}
		entries = append(entries, giveawayEntry{userID, weight, lastEntry})
	}
	return entries
}
