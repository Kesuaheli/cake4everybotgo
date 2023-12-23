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

package twitch

import (
	"log"

	"github.com/kesuaheli/twitchgo"
)

func MessageHandler(t *twitchgo.Twitch, message *twitchgo.Message) {
	log.Printf("Twitch: [%s] <%s> %s", message.Command.Arguments[0], message.Source, message.Command.Data)

}
