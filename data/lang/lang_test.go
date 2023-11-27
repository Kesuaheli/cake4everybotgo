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

package lang

import (
	"testing"

	"github.com/spf13/viper"
)

func TestIsLoaded(t *testing.T) {
	type args struct {
		lang string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{name: "de", args: args{lang: "de"}, want: true},
		{name: "de", args: args{lang: "DE"}, want: true},
		{name: "de", args: args{lang: "en"}, want: true},
		{name: "en", args: args{lang: "en_us"}, want: true},
		{name: "en", args: args{lang: "en-US"}, want: true},
		{name: "en", args: args{lang: "en US"}, want: true},
		{name: "en", args: args{lang: "en.us"}, want: false},
	}

	langsMap["de"] = viper.New()
	langsMap["en"] = viper.New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsLoaded(tt.args.lang); got != tt.want {
				t.Errorf("IsLoaded() = %v, want %v", got, tt.want)
			}
		})
	}

	delete(langsMap, "de")
	delete(langsMap, "en")
}
