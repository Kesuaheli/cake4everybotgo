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
	"github.com/bwmarrin/discordgo"

	"cake4everybot/data/lang"
	"cake4everybot/event/command/util"
)

var minValue = 1.0

func subCommandSet() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionSetDay(),
		commandOptionSetMonth(),
		commandOptionSetYear(),
		commandOptionSetVisible(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.set"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set"),
		Description:              lang.GetDefault(tp + "option.set.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.description"),
		Options:                  options,
	}
}

func commandOptionSetDay() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.day"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.day"),
		Description:              lang.GetDefault(tp + "option.set.option.day.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.day.description"),
		Required:                 true,
		Autocomplete:             true,
		MinValue:                 &minValue,
		MaxValue:                 31,
	}
}

func commandOptionSetMonth() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.month"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.month"),
		Description:              lang.GetDefault(tp + "option.set.option.month.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.month.description"),
		Required:                 true,
		Autocomplete:             true,
		MinValue:                 &minValue,
		MaxValue:                 12,
	}
}

func commandOptionSetYear() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.year"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.year"),
		Description:              lang.GetDefault(tp + "option.set.option.year.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.year.description"),
		Autocomplete:             true,
	}
}

func commandOptionSetVisible() *discordgo.ApplicationCommandOption {
	choices := []*discordgo.ApplicationCommandOptionChoice{
		{
			Name:              lang.GetDefault("discord.command.generic.yes"),
			NameLocalizations: *util.TranslateLocalization("discord.command.generic.yes"),
			Value:             true,
		},
		{
			Name:              lang.GetDefault("discord.command.generic.no"),
			NameLocalizations: *util.TranslateLocalization("discord.command.generic.no"),
			Value:             false,
		},
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionBoolean,
		Name:                     lang.GetDefault(tp + "option.set.option.visible"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.visible"),
		Description:              lang.GetDefault(tp + "option.set.option.visible.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.visible.description"),
		Choices:                  choices,
	}
}

func subCommandRemove() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.remove"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.remove"),
		Description:              lang.GetDefault(tp + "option.remove.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.remove.description"),
	}
}

func subCommandList() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionListMonth(),
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "option.list"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.list"),
		Description:              lang.GetDefault(tp + "option.list.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.list.description"),
		Options:                  options,
	}
}

func commandOptionListMonth() *discordgo.ApplicationCommandOption {
	var choices []*discordgo.ApplicationCommandOptionChoice
	for m := 1; m <= 12; m++ {
		choices = append(choices, monthChoice(m))
	}

	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.list.option.month"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.list.option.month"),
		Description:              lang.GetDefault(tp + "option.list.option.month.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.list.option.month.description"),
		Required:                 true,
		Choices:                  choices,
		MinValue:                 &minValue,
		MaxValue:                 12,
	}
}
