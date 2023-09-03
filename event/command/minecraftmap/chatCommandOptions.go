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

package minecraftmap

import (
	"cake4everybot/data/lang"
	"cake4everybot/event/command/util"

	"github.com/bwmarrin/discordgo"
)

func subCommandSet() *discordgo.ApplicationCommandOption {
	options := []*discordgo.ApplicationCommandOption{
		commandOptionSetSet(),
		commandOptionSetID(),
		commandOptionSetLabel(),
		commandOptionSetWorld(),
		commandOptionSetPosX(),
		commandOptionSetPosY(),
		commandOptionSetPosZ(),
		commandOptionSetIcon(),
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

func commandOptionSetSet() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     lang.GetDefault(tp + "option.set.option.set"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.set"),
		Description:              lang.GetDefault(tp + "option.set.option.set.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.set.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetID() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     lang.GetDefault(tp + "option.set.option.ID"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.ID"),
		Description:              lang.GetDefault(tp + "option.set.option.ID.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.ID.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetLabel() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     lang.GetDefault(tp + "option.set.option.label"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.label"),
		Description:              lang.GetDefault(tp + "option.set.option.label.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.label.description"),
		Required:                 true,
		Autocomplete:             false,
	}
}

func commandOptionSetWorld() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     lang.GetDefault(tp + "option.set.option.world"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.world"),
		Description:              lang.GetDefault(tp + "option.set.option.world.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.world.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetPosX() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.PosX"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.PosX"),
		Description:              lang.GetDefault(tp + "option.set.option.PosX.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.PosX.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetPosY() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.PosY"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.PosY"),
		Description:              lang.GetDefault(tp + "option.set.option.PosY.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.PosY.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetPosZ() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionInteger,
		Name:                     lang.GetDefault(tp + "option.set.option.PosZ"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.PosZ"),
		Description:              lang.GetDefault(tp + "option.set.option.PosZ.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.PosZ.description"),
		Autocomplete:             false,
	}
}

func commandOptionSetIcon() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionString,
		Name:                     lang.GetDefault(tp + "option.set.option.icon"),
		NameLocalizations:        *util.TranslateLocalization(tp + "option.set.option.icon"),
		Description:              lang.GetDefault(tp + "option.set.option.icon.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "option.set.option.icon.description"),
		Autocomplete:             false,
	}
}
