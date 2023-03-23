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

package birthday

import (
	"log"
	"sort"

	"github.com/bwmarrin/discordgo"

	"cake4everybot/database"
	"cake4everybot/event/command/util"
)

type birthdayBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

type birthdayEntry struct {
	id               uint64
	day, month, year int
	visible          bool
}

// getBirthday copies all birthday fields into
// the struct pointed at by b.
//
// If the user from b.id is not found it returns
// sql.ErrNoRows.
func (cmd birthdayBase) getBirthday(b *birthdayEntry) (err error) {
	row := database.QueryRow("SELECT day,month,year,visible FROM birthdays WHERE id=?", b.id)
	return row.Scan(&b.day, &b.month, &b.year, &b.visible)
}

// hasBirthday returns true whether the given
// user id has entered their birthday.
func (cmd birthdayBase) hasBirthday(id uint64) (hasBirthday bool, err error) {
	err = database.QueryRow("SELECT EXISTS(SELECT id FROM birthdays WHERE id=?)", id).Scan(&hasBirthday)
	return hasBirthday, err
}

// setBirthday inserts a new database entry with
// the values from b.
func (cmd birthdayBase) setBirthday(b birthdayEntry) {
	_, err := database.Exec("INSERT INTO birthdays(id,day,month,year,visible) VALUES(?,?,?,?,?);", b.id, b.day, b.month, b.year, b.visible)
	if err != nil {
		log.Printf("Error on set birthday: %v", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if b.visible {
		cmd.Replyf("Added your Birthday on %d.%d.%d!", b.day, b.month, b.year)
	} else {
		cmd.ReplyHiddenf("Added your Birthday on %d.%d.%d!\nYou can close this now", b.day, b.month, b.year)
	}
}

// updateBirthday updates an existing database
// entry with the values from b.
func (cmd birthdayBase) updateBirthday(b birthdayEntry) {
	_, err := database.Exec("UPDATE birthdays SET day=?,month=?,year=?,visible=? WHERE id=?;", b.day, b.month, b.year, b.visible, b.id)
	if err != nil {
		log.Printf("Error on update birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if b.visible {
		cmd.Replyf("Updated your Birthday to '%d.%d.%d'!", b.day, b.month, b.year)
	} else {
		cmd.ReplyHiddenf("Updated your Birthday to '%d.%d.%d'!\nYou can close this now.", b.day, b.month, b.year)
	}
}

// removeBirthday deletes the existing birthday
// entry for the given id.
func (cmd birthdayBase) removeBirthday(id uint64) {
	b := birthdayEntry{id: id}
	err := cmd.getBirthday(&b)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	_, err = database.Exec("DELETE FROM birthdays WHERE id=?;", b.id)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if b.visible {
		cmd.Replyf("Removed your Birthday from the bot!\nWas on '%d.%d.%d'.", b.day, b.month, b.year)
	} else {
		cmd.ReplyHiddenf("Removed your Birthday from the bot!\nWas on '%d.%d.%d'.\nYou can close this now.", b.day, b.month, b.year)
	}
}

// getBirthdaysMonth return a sorted slice of
// birthday entries that matches the given month.
func (cmd birthdayBase) getBirthdaysMonth(month int) (birthdays []birthdayEntry, err error) {
	var numOfEntries int64
	err = database.QueryRow("SELECT COUNT(*) FROM birthdays WHERE month=?", month).Scan(&numOfEntries)
	if err != nil {
		return nil, err
	}

	birthdays = make([]birthdayEntry, numOfEntries)
	if len(birthdays) == 0 {
		return birthdays, nil
	}

	rows, err := database.Query("SELECT id,day,year,visible FROM birthdays WHERE month=?", month)
	if err != nil {
		return birthdays, err
	}
	defer rows.Close()

	for rows.Next() {
		b := birthdayEntry{month: month}
		err = rows.Scan(&b.id, &b.day, &b.year, &b.visible)
		if err != nil {
			return birthdays, err
		}
		birthdays = append(birthdays, b)
	}

	sort.Slice(birthdays, func(i, j int) bool {
		return birthdays[i].day < birthdays[j].day
	})

	return birthdays, nil
}
