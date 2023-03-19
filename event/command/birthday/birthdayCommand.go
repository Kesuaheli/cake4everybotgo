// Copyright 2022-2023 Kesuaheli
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
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"

	"cake4everybot/database"
	"cake4everybot/event/command/util"
)

type Birthday struct{}

func (c Birthday) AppCmd() *discordgo.ApplicationCommand {

	names := map[discordgo.Locale]string{
		discordgo.German: "geburtstag",
	}

	descriptions := map[discordgo.Locale]string{
		discordgo.German: "Verschiedene Einstellungen für den Geburtstagsbot",
	}
	var minValue = 1.0
	cmd := &discordgo.ApplicationCommand{
		Name:                     "birthday",
		NameLocalizations:        &names,
		Description:              "Various settings for the birthday bot",
		DescriptionLocalizations: &descriptions,
		// Type:        discordgo.ChatApplicationCommand,
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionSubCommand,
				Name:        "set",
				Description: "Enter or change your birthday",
				Options: []*discordgo.ApplicationCommandOption{
					{
						Type:         discordgo.ApplicationCommandOptionInteger,
						Name:         "day",
						Description:  "On wich day of month is your birthday?",
						Required:     true,
						Autocomplete: true,
						MinValue:     &minValue,
						MaxValue:     31,
					},
					{
						Type:         discordgo.ApplicationCommandOptionInteger,
						Name:         "month",
						Description:  "On wich month of the year is your birthday?",
						Required:     true,
						Autocomplete: true,
						MinValue:     &minValue,
						MaxValue:     12,
					},
					{
						Type:        discordgo.ApplicationCommandOptionInteger,
						Name:        "year",
						Description: "In wich year were you born?",
					},
					{
						Type:        discordgo.ApplicationCommandOptionBoolean,
						Name:        "visible",
						Description: "Should your name and birthday be discoverable by others? (defaults to \"Yes\")",
						Choices: []*discordgo.ApplicationCommandOptionChoice{
							{Name: "Yes", Value: true},
							{Name: "No", Value: false},
						},
					},
				},
			},
		},
	}

	return cmd

}

func (c Birthday) CmdHandler() func(s *discordgo.Session, i *discordgo.InteractionCreate) {

	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		data := i.ApplicationCommandData()
		options := data.Options[0].Options
		switch data.Options[0].Name {
		case "set":
			birthdaySetHandler(s, i, options)
		}
	}
}

func birthdaySetHandler(s *discordgo.Session, i *discordgo.InteractionCreate, options []*discordgo.ApplicationCommandInteractionDataOption) {
	var (
		iu      util.InteractionUtil = util.InteractionUtil{Session: s, Interaction: i}
		day     *discordgo.ApplicationCommandInteractionDataOption
		month   *discordgo.ApplicationCommandInteractionDataOption
		year    *discordgo.ApplicationCommandInteractionDataOption
		visible *discordgo.ApplicationCommandInteractionDataOption
	)
	for _, opt := range options {
		switch opt.Name {
		case "day":
			day = opt
		case "month":
			month = opt
		case "year":
			year = opt
		case "visible":
			visible = opt
		}
	}

	switch i.Type {
	case discordgo.InteractionApplicationCommandAutocomplete:
		var choices []*discordgo.ApplicationCommandOptionChoice

		if day.Focused {
			i, _ := strconv.Atoi(day.Value.(string))
			choices = append(choices, dayChoice(i))
		} else if month.Focused {
			// jan, mar, may, jul, aug, oct, dec
			choices = append(choices,
				monthChoice(1),
				monthChoice(3),
				monthChoice(5),
				monthChoice(7),
				monthChoice(8),
				monthChoice(10),
				monthChoice(12),
			)
			if day.IntValue() < 31 {
				// apr, jun, sep, nov
				choices = append(choices,
					monthChoice(4),
					monthChoice(6),
					monthChoice(9),
					monthChoice(11),
				)
			}
			if day.IntValue() == 29 && (year.IntValue()%4) == 0 {
				// feb (leap year)
				choices = append(choices,
					monthChoice(2),
				)
			}
			if day.IntValue() < 29 {
				// feb
				choices = append(choices,
					monthChoice(2),
				)
			}
		}
		iu.ReplyAutocomplete(choices)
	case discordgo.InteractionApplicationCommand:
		authorID, err := strconv.ParseUint(i.Member.User.ID, 10, 64)
		if err != nil {
			log.Printf("Error on parse author id of birthday command: %v\n", err)
			iu.ReplyError()
			return
		}

		var (
			newDay     int = int(day.IntValue())
			newMonth   int = int(month.IntValue())
			newYear    int
			newVisible bool = true
		)
		if year != nil {
			newYear = int(year.IntValue())
		}
		if visible != nil {
			newVisible = visible.BoolValue()
		}

		hasBDay, err := hasBirthday(authorID)
		if err != nil {
			log.Printf("Error on getting birthday data: %v\n", err)
			iu.ReplyError()
			return
		}

		if hasBDay {
			updateBirthday(&iu, authorID, newDay, newMonth, newYear, newVisible)
		} else {
			setBirthday(&iu, authorID, newDay, newMonth, newYear, newVisible)
		}
	}
}

func monthChoice(month int) (choice *discordgo.ApplicationCommandOptionChoice) {
	choice = &discordgo.ApplicationCommandOptionChoice{
		Name:  "InternalError",
		Value: month,
	}
	switch month {
	case 1:
		choice.Name = "January"
	case 2:
		choice.Name = "February"
	case 3:
		choice.Name = "March"
	case 4:
		choice.Name = "April"
	case 5:
		choice.Name = "May"
	case 6:
		choice.Name = "June"
	case 7:
		choice.Name = "July"
	case 8:
		choice.Name = "August"
	case 9:
		choice.Name = "September"
	case 10:
		choice.Name = "October"
	case 11:
		choice.Name = "November"
	case 12:
		choice.Name = "December"
	}
	return
}

func dayChoice(day int) (choice *discordgo.ApplicationCommandOptionChoice) {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  fmt.Sprint(day),
		Value: day,
	}
}

// getBirthday returns all birthday fields of
// the given user id.
//
// If that user is not found it returns
// sql.ErrNoRows.
func getBirthday(id uint64) (day int, month int, year int, visible bool, err error) {
	row := database.QueryRow("SELECT day,month,year,visible FROM birthdays WHERE id=?", id)
	err = row.Scan(&day, &month, &year, &visible)
	return
}

// hasBirthday returns true whether the given
// user id has entered their birthday.
func hasBirthday(id uint64) (hasBirthday bool, err error) {
	err = database.QueryRow("SELECT id FROM birthdays WHERE id=?", id).Err()

	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil, err
}

// setBirthday inserts a new birthday entry with
// the given values into the database.
func setBirthday(iu *util.InteractionUtil, id uint64, day int, month int, year int, visible bool) {
	_, err := database.Exec("INSERT INTO birthdays(id,day,month,year,visible) VALUES(?,?,?,?,?);", id, day, month, year, visible)
	if err != nil {
		log.Printf("Error on set birthday: %v", err)
		iu.ReplyError()
		return
	}

	// notify the user
	if visible {
		iu.Replyf("Added your Birthday on %d.%d.%d!", day, month, year)
	} else {
		iu.ReplyHiddenf("Added your Birthday on %d.%d.%d!\nYour can close this now", day, month, year)
	}
}

// updateBirthday updates an existing birthday
// entry with the given values to database.
func updateBirthday(iu *util.InteractionUtil, id uint64, day int, month int, year int, visible bool) {
	_, err := database.Exec("UPDATE birthdays SET day=?,month=?,year=?,visible=? WHERE id=?;", day, month, year, visible, id)
	if err != nil {
		log.Printf("Error on update birthday: %v\n", err)
		iu.ReplyError()
		return
	}

	// notify the user
	if visible {
		iu.Replyf("Updated your Birthday to '%d.%d.%d'!", day, month, year)
	} else {
		iu.ReplyHiddenf("Updated your Birthday to '%d.%d.%d'!\nYour can close this now.", day, month, year)
	}
}
