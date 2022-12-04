// Copyright 2022 Kesuaheli
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

package main

import (
	"fmt"
	"os"
	"os/signal"

	"cake4everybot/event"

	"github.com/bwmarrin/discordgo"
)

const banner string = "\n" +
	"   ______      __        __ __  ______                                     \n" +
	"  / ____/___ _/ /_____  / // / / ____/   _____  _______  ______  ____  ___ \n" +
	" / /   / __ `/ //_/ _ \\/ // /_/ __/ | | / / _ \\/ ___/ / / / __ \\/ __ \\/ _ \\\n" +
	"/ /___/ /_/ / ,< /  __/__  __/ /___ | |/ /  __/ /  / /_/ / /_/ / / / /  __/\n" +
	"\\____/\\__,_/_/|_|\\___/  /_/ /_____/ |___/\\___/_/   \\__, /\\____/_/ /_/\\___/ \n" +
	"                                                  /____/                   \n" +
	"      ____  _                          __            ____        __        \n" +
	"     / __ \\(_)_____________  _________/ /           / __ )____  / /_       \n" +
	"    / / / / / ___/ ___/ __ \\/ ___/ __  /  ______   / __  / __ \\/ __/       \n" +
	"   / /_/ / (__  ) /__/ /_/ / /  / /_/ /  /_____/  / /_/ / /_/ / /_         \n" +
	"  /_____/_/____/\\___/\\____/_/   \\__,_/           /_____/\\____/\\__/         \n" +
	"\n" +
	"Cake4Everybot, developed by @Kesuaheli#5868 and the ideas of the community ♥\n" +
	"Copyright 2022 Kesuaheli\n\n"

func main() {

	fmt.Print(banner)

	dcToken, err := os.ReadFile("lib/dcToken.0")
	if err != nil {
		panic("could not read token file")
	}
	s, err := discordgo.New("Bot " + string(dcToken))
	if err != nil {
		panic(fmt.Sprintf("invalid bot parameters: %v", err))
	}

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		fmt.Printf("Logged in as %s#%s\n", s.State.User.Username, s.State.User.Discriminator)
	})

	// Add event listeners
	event.AddListeners(s)

	// open connection to Discord and login
	err = s.Open()
	if err != nil {
		panic(fmt.Sprintf("could not open the discord session: %v", err))
	}
	defer s.Close()

	// register all command and co.

	dcGuildID, err := os.ReadFile("lib/dcGuildID.0")
	if err != nil {
		panic("could not read guildID file")
	}
	err = event.Register(s, string(dcGuildID))
	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	// Wait to end the bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop

	fmt.Println("\nGracefully shutting down. Byee")
}
