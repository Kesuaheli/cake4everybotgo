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
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"

	"cake4everybot/data/lang"
	"cake4everybot/event/command/util"
)

var (
	day_choices_prefix = [][]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 20, 30},         // 0
		{1, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},     // 1
		{2, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 12}, // 2
		{3, 30, 31, 13, 23},                             // 3
		{4, 14, 24},                                     // 4
		{5, 15, 25},                                     // 5
		{6, 16, 26},                                     // 6
		{7, 17, 27},                                     // 7
		{8, 18, 28},                                     // 8
		{9, 19, 29},                                     // 9
	}
)

func dayChoices(start string, month int, leapYear bool) (choices []*discordgo.ApplicationCommandOptionChoice) {
	i, _ := strconv.Atoi(start)
	if i < 0 || i > getDays(month, leapYear) {
		return choices
	}

	if i >= len(day_choices_prefix) || len(day_choices_prefix[i]) == 0 {
		choices = append(choices, intChoice(i))
		return choices
	}

	for _, c := range day_choices_prefix[i] {
		if c > getDays(month, leapYear) {
			continue
		}
		choices = append(choices, intChoice(c))
	}
	return choices
}

func monthChoices(start string, day int, leapYear bool) (choices []*discordgo.ApplicationCommandOptionChoice) {
	i, err := strconv.Atoi(start)
	if err != nil {
		for month := 1; month <= 12; month++ {
			if day > getDays(month, leapYear) {
				continue
			}
			key := fmt.Sprintf("%smonth.%d", tp, month-1)

			hasPrefix := strings.Contains(start, fmt.Sprint(month))
			hasPrefix = hasPrefix || strings.Contains(fmt.Sprint(month), start)

			for _, l := range lang.GetLangs() {
				name := lang.Get(key, l)
				hasPrefix = hasPrefix || strings.Contains(name, start)
				hasPrefix = hasPrefix || strings.Contains(start, name)
			}

			if hasPrefix {
				choices = append(choices, monthChoice(month))
			}
		}
		return choices
	}

	if i < 0 || i > 12 {
		return choices
	}

	if i == 1 {
		choices = append(choices, monthChoice(1))
		choices = append(choices, monthChoice(10))
		choices = append(choices, monthChoice(11))
		choices = append(choices, monthChoice(12))
		return choices
	}
	if i == 2 {
		if day <= getDays(2, leapYear) {
			choices = append(choices, monthChoice(2))
		}
		choices = append(choices, monthChoice(12))
		return choices
	}
	if util.ContainsInt([]int{3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, i) {
		choices = append(choices, monthChoice(i))
		return choices
	}

	for month := 1; month <= 12; month++ {
		if day > getDays(month, leapYear) {
			continue
		}
		choices = append(choices, monthChoice(month))
	}
	return choices
}

func yearChoices(start string, day, month int) (choices []*discordgo.ApplicationCommandOptionChoice) {
	maxDate := time.Now().AddDate(-16, 0, 0)

	var decades []int
	cur_decade := maxDate.Year() / 10 * 10
	for y := cur_decade; y >= cur_decade-100; y = y - 10 {
		decades = append(decades, y)
	}

	y, err := strconv.Atoi(start)
	if err != nil || y == 0 {
		for _, dec := range decades {
			choices = append(choices, intChoice(dec))
		}
		return choices
	}
	y = int(math.Abs(float64(y)))

	digits := len(fmt.Sprint(y))

	rm := func(s []int, i int) []int {
		if i < 0 || i >= len(s) {
			return s
		}
		if i == len(s)-1 {
			return s[:i]
		}
		return append(s[:i], s[i+1:]...)
	}
	decades_cp := make([]int, len(decades))
	copy(decades_cp, decades)
	for i := len(decades_cp) - 1; i >= 0; i-- {
		dec := fmt.Sprint(decades_cp[i])
		if len(dec) < digits {
			decades = rm(decades, i)
			continue
		}
		if !strings.HasPrefix(dec, fmt.Sprint(y)) {
			decades = rm(decades, i)
			continue
		}
	}

	mustLeapYear := day == 29 && month == 2
	if !mustLeapYear && len(decades) > 2 ||
		mustLeapYear && len(decades) > 10 ||
		len(decades) == 0 {
		for _, dec := range decades {
			choices = append(choices, intChoice(dec))
		}
		return choices
	}

	for _, dec := range decades {
		for y := 0; y < 10; y++ {
			if time.Date(dec+y, time.Month(month), day, 0, 0, 0, 0, time.Local).After(maxDate) {
				continue
			}
			if mustLeapYear && (dec+y)%4 != 0 {
				continue
			}
			choices = append(choices, intChoice(dec+y))
		}
	}
	return choices
}

func intChoice(i int) (choice *discordgo.ApplicationCommandOptionChoice) {
	return &discordgo.ApplicationCommandOptionChoice{
		Name:  fmt.Sprint(i),
		Value: i,
	}
}

func monthChoice(month int) (choice *discordgo.ApplicationCommandOptionChoice) {
	key := fmt.Sprintf("%smonth.%d", tp, month-1)

	return &discordgo.ApplicationCommandOptionChoice{
		Name:              lang.GetDefault(key),
		NameLocalizations: *util.TranslateLocalization(key),
		Value:             month,
	}
}

// getDays returns the maximum number of days in the given month.
// When the given month is february (month: 2), getDays returns 29,
// as it is the max. number of day the february can have.
func getDays(month int, leapYear bool) int {
	if util.ContainsInt([]int{2}, month) {
		if leapYear {
			return 29
		} else {
			return 28
		}
	} else if util.ContainsInt([]int{4, 6, 9, 11}, month) {
		return 30
	} else if util.ContainsInt([]int{1, 3, 5, 7, 8, 10, 12}, month) {
		return 31
	}
	return 0
}
