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

package webserver

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

// Feed is the object that holds a incomming notification feed from
// youtube. This could be a new video (upload/publish) or an update
// of an existing one.
type Feed struct {
	XMLName xml.Name `xml:"feed" json:"-"`
	ID      string   `xml:"entry>videoId"`
	Channel string   `xml:"entry>channelId"`
	Title   string   `xml:"entry>title"`
}

func handleYTGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Println(r.RequestURI)
	r.Header.Write(log.Writer())

	topic := r.FormValue("hub.topic")
	challenge := r.FormValue("hub.challenge")
	mode := r.FormValue("hub.mode")

	if topic == "" || challenge == "" || mode == "" {
		log.Println("Missing at least one of topic, challenge, mode in query parameters")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if ok, _ := regexp.MatchString("(?:un)?subscribe", mode); !ok {
		log.Printf("Unsupported mode '%s'", mode)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("Topic:", topic)
	log.Println("Challenge: ", challenge)
	log.Println("Mode: ", mode)

	topicURL, err := url.Parse(topic)
	if err != nil {
		log.Printf("Error on parse topic url '%s': %v\n", topic, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// only accept youtube video feed
	if topicURL.Host != "www.youtube.com" {
		log.Printf("Topic host is not youtube: %s\n", topicURL.Host)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	//Todo: check for path too: "https://www.youtube.com/xml/feeds/videos.xml?channel_id={channel_id}"

	channelID := topicURL.Query().Get("channel_id")
	log.Println("ChannelID: ", channelID)

	if channelID != "UC6sb0bkXREewXp2AkSOsOqg" {
		log.Printf("Requested unknown channel: %s\n", channelID)
		w.WriteHeader(http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(challenge))
	log.Printf("Accepted '%s' from %s for channel %s\n", mode, topicURL.Host, channelID)
}

func handleYTPost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.Header.Write(log.Writer())

	// only accept atom feed
	content := r.Header.Get("Content-Type")
	if content != "application/atom+xml" {
		log.Printf("Content-Type '%s' not supported\n", content)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	buf, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	feed := Feed{}
	err = xml.Unmarshal(buf, &feed)
	if err != nil {
		log.Printf("Error on parse XML body: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Println("New Notification :)")
	log.Println("Parsed Feed:", feed)
}
