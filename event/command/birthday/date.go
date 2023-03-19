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

	"github.com/bwmarrin/discordgo"
)

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
