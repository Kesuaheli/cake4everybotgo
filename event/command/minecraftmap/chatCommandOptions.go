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

func subCommandMarker() *discordgo.ApplicationCommandOption {
	return &discordgo.ApplicationCommandOption{
		Type:                     discordgo.ApplicationCommandOptionSubCommand,
		Name:                     lang.GetDefault(tp + "subcommand.marker"),
		NameLocalizations:        *util.TranslateLocalization(tp + "subcommand.marker"),
		Description:              lang.GetDefault(tp + "subcommand.marker.description"),
		DescriptionLocalizations: *util.TranslateLocalization(tp + "subcommand.marker.description"),
	}
}
