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

package youtube

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setServer(f http.HandlerFunc) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(f))
}

func Test_handleYTGet_as_post(t *testing.T) {
	server := setServer(HandleGet)

	resp, err := http.Post(server.URL, "*/*", nil)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 but got %d", resp.StatusCode)
	}
}

func Test_handleYTGet_without_query(t *testing.T) {
	server := setServer(HandleGet)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

func Test_handleYTGet_with_unknown_mode(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "GEhceZmXvC1thlYAC19J"
	challenge := "aOCutOf4xe"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=ebircsbus"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_with_wrong_topic_host(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "7yNIWPVJ9YqAIYm53qV7"
	challenge := "rwu71lWVJ"
	server.URL += "?hub.topic=https://www.example.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=subscribe"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_with_wrong_topic_path(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "3liKhDvwVjyv7W54TjB7"
	challenge := "YcAHL4cyp"
	server.URL += "?hub.topic=https://www.youtube.com/foo/bar.xml?channelI_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.challenge=subscribe"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_without_channel(t *testing.T) {
	server := setServer(HandleGet)

	challenge := "62xrFENuj"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml"
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.challenge=subscribe"

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_without_channel_subscription(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "HE513kwtQPFgR4KUNdng"
	challenge := "lTmnRsn0o"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=subscribe"

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_unsubscribe_used_channel(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "gRdEy5KIYyEatbYH6oQm"
	challenge := "FHbOv8Tie"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=unsubscribe"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTGet_return_challenge(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "MH2Fx521amDXyaBIJhKw"
	challenge := "MNKXlkUTc"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=subscribe"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Errorf("Expected 2xx but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != challenge {
		t.Errorf("Expected challenge '%s' but got '%s'", challenge, string(b))
	}
}

func Test_handleYTGet_unsubscribe_again_after_subscribe(t *testing.T) {
	server := setServer(HandleGet)

	channelID := "XeFfTreP5kIokcOgdsbF"
	challenge := "tvI3iU10O"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channelID
	server.URL += "&hub.challenge=" + challenge
	server.URL += "&hub.mode=subscribe"

	SubscribeChannel(channelID)
	defer UnsubscribeChannel(channelID)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		t.Errorf("Expected 2xx but got %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) != challenge {
		t.Errorf("Expected challenge '%s' but got '%s'", challenge, string(b))
	}

	UnsubscribeChannel(channelID)

	resp, err = http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}

	b, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}
	if string(b) == challenge {
		t.Error("Expected to not return challenge")
	}
}

func Test_handleYTPost_as_get(t *testing.T) {
	server := setServer(HandlePost)

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405 but got %d", resp.StatusCode)
	}
}

func Test_handleYTPost_with_wrong_content_type(t *testing.T) {
	server := setServer(HandlePost)

	resp, err := http.Post(server.URL, "foo/bar", nil)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusUnsupportedMediaType {
		t.Errorf("Expected 415 but got %d", resp.StatusCode)
	}
}

func Test_handleYTPost_with_invalid_content(t *testing.T) {
	server := setServer(HandlePost)

	resp, err := http.Post(server.URL, "application/atom+xml", nil)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	body := strings.NewReader("foo bar")
	resp, err = http.Post(server.URL, "application/atom+xml", body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

func Test_handleYTPost_with_incomplete_content(t *testing.T) {
	server := setServer(HandlePost)

	bodyEmptyFeed := strings.NewReader(`
	<feed>
	</feed>
	`)
	resp, err := http.Post(server.URL, "application/atom+xml", bodyEmptyFeed)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	bodyEmptyEntry := strings.NewReader(`
	<feed>
		<entry>
		</entry>
	</feed>
	`)
	resp, err = http.Post(server.URL, "application/atom+xml", bodyEmptyEntry)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

func Test_handleYTPost_with_invalid_entry_feed(t *testing.T) {
	server := setServer(HandlePost)

	// no yt namespace set
	body := strings.NewReader(`
	<feed>
		<entry>
			<title>This is an automated test</title>
			<yt:videoId>go-test</yt:videoId>
			<yt:channelId>go-test</yt:channelId>
		</entry>
	</feed>
	`)
	resp, err := http.Post(server.URL, "application/atom+xml", body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}

	// wrong yt namespace set
	body = strings.NewReader(`
	<feed xmlns:yt="http://www.cake4everyone.de/">
		<entry>
			<title>This is an automated test</title>
			<yt:videoId>go-test</yt:videoId>
			<yt:channelId>go-test</yt:channelId>
		</entry>
	</feed>
	`)
	resp, err = http.Post(server.URL, "application/atom+xml", body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

func Test_handleYTPost_with_valid_content(t *testing.T) {
	server := setServer(HandlePost)

	body := strings.NewReader(`
	<feed xmlns:yt="http://www.youtube.com/">
		<entry>
			<title>This is an automated test</title>
			<yt:videoId>go-test</yt:videoId>
			<yt:channelId>go-test</yt:channelId>
		</entry>
	</feed>
	`)
	resp, err := http.Post(server.URL, "application/atom+xml", body)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202 but got %d", resp.StatusCode)
	}
}
