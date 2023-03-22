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

	"github.com/bwmarrin/discordgo"

	"cake4everybot/database"
	"cake4everybot/event/command/util"
)

type birthdayBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

// getBirthday returns all birthday fields of
// the given user id.
//
// If that user is not found it returns
// sql.ErrNoRows.
func (cmd birthdayBase) getBirthday(id uint64) (day int, month int, year int, visible bool, err error) {
	row := database.QueryRow("SELECT day,month,year,visible FROM birthdays WHERE id=?", id)
	err = row.Scan(&day, &month, &year, &visible)
	return
}

// hasBirthday returns true whether the given
// user id has entered their birthday.
func (cmd birthdayBase) hasBirthday(id uint64) (hasBirthday bool, err error) {
	err = database.QueryRow("SELECT EXISTS(SELECT id FROM birthdays WHERE id=?)", id).Scan(&hasBirthday)
	return hasBirthday, err
}

// setBirthday inserts a new birthday entry with
// the given values into the database.
func (cmd birthdayBase) setBirthday(id uint64, day int, month int, year int, visible bool) {
	_, err := database.Exec("INSERT INTO birthdays(id,day,month,year,visible) VALUES(?,?,?,?,?);", id, day, month, year, visible)
	if err != nil {
		log.Printf("Error on set birthday: %v", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if visible {
		cmd.Replyf("Added your Birthday on %d.%d.%d!", day, month, year)
	} else {
		cmd.ReplyHiddenf("Added your Birthday on %d.%d.%d!\nYou can close this now", day, month, year)
	}
}

// updateBirthday updates an existing birthday
// entry with the given values to database.
func (cmd birthdayBase) updateBirthday(id uint64, day int, month int, year int, visible bool) {
	_, err := database.Exec("UPDATE birthdays SET day=?,month=?,year=?,visible=? WHERE id=?;", day, month, year, visible, id)
	if err != nil {
		log.Printf("Error on update birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if visible {
		cmd.Replyf("Updated your Birthday to '%d.%d.%d'!", day, month, year)
	} else {
		cmd.ReplyHiddenf("Updated your Birthday to '%d.%d.%d'!\nYou can close this now.", day, month, year)
	}
}

// removeBirthday deletes the existing birthday
// entry for the given id.
func (cmd birthdayBase) removeBirthday(id uint64) {
	day, month, year, visible, err := cmd.getBirthday(id)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	_, err = database.Exec("DELETE FROM birthdays WHERE id=?;", id)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if visible {
		cmd.Replyf("Removed your Birthday from the bot!\nWas on '%d.%d.%d'.", day, month, year)
	} else {
		cmd.ReplyHiddenf("Removed your Birthday from the bot!\nWas on '%d.%d.%d'.\nYou can close this now.", day, month, year)
	}
}
