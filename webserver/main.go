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
	"cake4everybot/webserver/youtube"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func initHTTP() http.Handler {
	r := mux.NewRouter()
	r.Use(Logger)

	r.NotFoundHandler = http.HandlerFunc(handle404)

	r.HandleFunc("/favicon.ico", favicon)
	r.HandleFunc("/api/yt_pubsubhubbub/", youtube.HandleGet).Methods("GET")
	r.HandleFunc("/api/yt_pubsubhubbub/", youtube.HandlePost).Methods("POST")

	return r
}

// Run starts the webserver at the given address
func Run(addr string) {
	handler := initHTTP()

	go func() {
		err := http.ListenAndServe(addr, handler)
		if err != nil {
			log.Printf("Webserver ended with error: %v\n", err)
		} else {
			log.Println("Webserver ended!")
		}
	}()
}

func favicon(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, viper.GetString("webserver.favicon"))
}
