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

func Test_handleYTGet_with_wrong_topic_host(t *testing.T) {
	server := setServer(HandleGet)
	server.URL += "?hub.topic=https://www.example.com/"
	server.URL += "&hub.challenge=123aBc456"
	server.URL += "&hub.mode=subscribe"

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusForbidden {
		t.Errorf("Expected 403 but got %d", resp.StatusCode)
	}
}

func Test_handleYTGet_without_channel(t *testing.T) {
	server := setServer(HandleGet)
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml"
	server.URL += "&hub.challenge=123aBc456"
	server.URL += "&hub.challenge=subscribe"

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected 400 but got %d", resp.StatusCode)
	}
}

func Test_handleYTGet_return_challenge(t *testing.T) {
	server := setServer(HandleGet)
	channel := "UC6sb0bkXREewXp2AkSOsOqg"
	server.URL += "?hub.topic=https://www.youtube.com/xml/feeds/videos.xml?channel_id=" + channel
	server.URL += "&hub.challenge=123aBc456"
	server.URL += "&hub.mode=subscribe"

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

	if string(b) != "123aBc456" {
		t.Errorf("Expected challenge '123aBc456' but got '%s'", string(b))
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

func Test_handleYTPost_with_valid_content(t *testing.T) {
	server := setServer(HandlePost)

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

	if resp.StatusCode != http.StatusAccepted {
		t.Errorf("Expected 202 but got %d", resp.StatusCode)
	}
}
