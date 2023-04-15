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
	"fmt"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"cake4everybot/database"
	"cake4everybot/event/command/util"
)

const (
	// Prefix for translation key, i.e.:
	//   key := tp+"base" // => birthday
	tp = "discord.command.birthday."
)

type birthdayBase struct {
	util.InteractionUtil
	member *discordgo.Member
	user   *discordgo.User
}

type birthdayEntry struct {
	ID      uint64 `database:"id"`
	Day     int    `database:"day"`
	Month   int    `database:"month"`
	Year    int    `database:"year"`
	Visible bool   `database:"visible"`
}

// Returns a readable Form of the date
func (b birthdayEntry) String() string {
	bTime, err := time.Parse(time.DateOnly, fmt.Sprintf("%d-%02d-%02d", b.Year, b.Month, b.Day))
	if err != nil {
		log.Printf("couldn't parse date: %s", err)
		if b.Year == 0 {
			return fmt.Sprintf("%d.%d.", b.Day, b.Month)
		}
		return fmt.Sprintf("%d.%d.%d", b.Day, b.Month, b.Year)
	}

	layout := "Mon, _2 Jan 2006"
	if b.Year == 0 {
		layout = "_2 Jan"
	}
	return fmt.Sprintf(bTime.Format(layout))

}

// getBirthday copies all birthday fields into
// the struct pointed at by b.
//
// If the user from b.ID is not found it returns
// sql.ErrNoRows.
func (cmd birthdayBase) getBirthday(b *birthdayEntry) (err error) {
	row := database.QueryRow("SELECT day,month,year,visible FROM birthdays WHERE id=?", b.ID)
	return row.Scan(&b.Day, &b.Month, &b.Year, &b.Visible)
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
	_, err := database.Exec("INSERT INTO birthdays(id,day,month,year,visible) VALUES(?,?,?,?,?);", b.ID, b.Day, b.Month, b.Year, b.Visible)
	if err != nil {
		log.Printf("Error on set birthday: %v", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if b.Visible {
		cmd.Replyf("Added your Birthday on %d.%d.%d!", b.Day, b.Month, b.Year)
	} else {
		cmd.ReplyHiddenf("Added your Birthday on %d.%d.%d!\nYou can close this now", b.Day, b.Month, b.Year)
	}
}

// updateBirthday updates an existing database
// entry with the values from b.
func (cmd birthdayBase) updateBirthday(b birthdayEntry) {
	bdayOld := b
	if err := cmd.getBirthday(&bdayOld); err != nil {
		log.Printf("Error on update birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	var (
		updateNames []string
		updateVars  []any
		oldV        reflect.Value = reflect.ValueOf(bdayOld)
		v           reflect.Value = reflect.ValueOf(b)
	)
	for i := 0; i < v.NumField(); i++ {
		var (
			oldF = oldV.Field(i).Interface()
			f    = v.Field(i).Interface()
		)
		if f != oldF {
			tag := v.Type().Field(i).Tag.Get("database")
			if tag == "" {
				continue
			}
			updateNames = append(updateNames, tag)
			updateVars = append(updateVars, f)
		}

	}

	if len(updateNames) == 0 {
		cmd.ReplyHidden("Nothing changed! You set your birthday already to this date.")
		return
	}

	updateString := strings.Join(updateNames, "=?,") + "=?"
	_, err := database.Exec("UPDATE birthdays SET "+updateString+";", updateVars...)
	if err != nil {
		log.Printf("Error on update birthday: %v\n", err)
		cmd.ReplyError()
		return
	}

	// notify the user
	if b.Visible {
		cmd.Replyf("Updated your Birthday to '%d.%d.%d'!", b.Day, b.Month, b.Year)
	} else {
		cmd.ReplyHiddenf("Updated your Birthday to '%d.%d.%d'!\nYou can close this now.", b.Day, b.Month, b.Year)
	}
}

// removeBirthday deletes the existing birthday entry for the given
// id and returns the previously entered birthday.
func (cmd birthdayBase) removeBirthday(id uint64) birthdayEntry {
	b := birthdayEntry{ID: id}
	err := cmd.getBirthday(&b)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
		return b
	}

	_, err = database.Exec("DELETE FROM birthdays WHERE id=?;", b.ID)
	if err != nil {
		log.Printf("Error on remove birthday: %v\n", err)
		cmd.ReplyError()
	}
	return b
}

// getBirthdaysMonth return a sorted slice of
// birthday entries that matches the given month.
func (cmd birthdayBase) getBirthdaysMonth(month int) (birthdays []birthdayEntry, err error) {
	var numOfEntries int64
	err = database.QueryRow("SELECT COUNT(*) FROM birthdays WHERE month=?", month).Scan(&numOfEntries)
	if err != nil {
		return nil, err
	}

	birthdays = make([]birthdayEntry, 0, numOfEntries)
	if numOfEntries == 0 {
		return birthdays, nil
	}

	rows, err := database.Query("SELECT id,day,year,visible FROM birthdays WHERE month=?", month)
	if err != nil {
		return birthdays, err
	}
	defer rows.Close()

	for rows.Next() {
		b := birthdayEntry{Month: month}
		err = rows.Scan(&b.ID, &b.Day, &b.Year, &b.Visible)
		if err != nil {
			return birthdays, err
		}
		birthdays = append(birthdays, b)
	}

	sort.Slice(birthdays, func(i, j int) bool {
		return birthdays[i].Day < birthdays[j].Day
	})

	return birthdays, nil
}

// getBirthdaysDate return a slice of birthday
// entries that matches the given date.
func getBirthdaysDate(day int, month int) (birthdays []birthdayEntry, err error) {
	var numOfEntries int64
	err = database.QueryRow("SELECT COUNT(*) FROM birthdays WHERE day=? AND month=?", day, month).Scan(&numOfEntries)
	if err != nil {
		return nil, err
	}

	birthdays = make([]birthdayEntry, 0, numOfEntries)
	if numOfEntries == 0 {
		return birthdays, nil
	}

	rows, err := database.Query("SELECT id,year,visible FROM birthdays WHERE day=? AND month=?", day, month)
	if err != nil {
		return birthdays, err
	}
	defer rows.Close()

	for rows.Next() {
		b := birthdayEntry{Day: day, Month: month}
		err = rows.Scan(&b.ID, &b.Year, &b.Visible)
		if err != nil {
			return birthdays, err
		}
		birthdays = append(birthdays, b)
	}

	return birthdays, nil
}
