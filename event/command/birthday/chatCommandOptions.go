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

import "github.com/bwmarrin/discordgo"

var minValue = 1.0

func subCommandSet() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionSetDay(),
		commandOptionSetMonth(),
		commandOptionSetYear(),
		commandOptionSetVisible(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "set",
		Description: "Enter or change your birthday",
		Options:     options,
	}
}

func commandOptionSetDay() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionInteger,
		Name:         "day",
		Description:  "On wich day of month is your birthday?",
		Required:     true,
		Autocomplete: true,
		MinValue:     &minValue,
		MaxValue:     31,
	}
}

func commandOptionSetMonth() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:         discordgo.ApplicationCommandOptionInteger,
		Name:         "month",
		Description:  "On wich month of the year is your birthday?",
		Required:     true,
		Autocomplete: true,
		MinValue:     &minValue,
		MaxValue:     12,
	}
}

func commandOptionSetYear() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionInteger,
		Name:        "year",
		Description: "In wich year were you born?",
	}
}

func commandOptionSetVisible() *discordgo.ApplicationCommandOption {
	choices := []*discordgo.ApplicationCommandOptionChoice{
		{Name: "Yes", Value: true},
		{Name: "No", Value: false},
	}

	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionBoolean,
		Name:        "visible",
		Description: "Should your name and birthday be discoverable by others? (defaults to \"Yes\")",
		Choices:     choices,
	}
}

func subCommandRemove() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:        discordgo.ApplicationCommandOptionSubCommand,
		Name:        "remove",
		Description: "Remove your entered Birthday from the bot",
	}
}
