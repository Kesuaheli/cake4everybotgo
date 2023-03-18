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

package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

var db *sql.DB

type connection_config struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
}

func Connect() {
	log.Println("Connecting to Database...")

	// setting default values
	config := connection_config{
		Host: "localhost",
		Port: 3306,
	}
	err := viper.UnmarshalKey("mysql", &config)
	if err != nil {
		log.Fatalf("Could not read msql connection data from config: %v", err)
	}

	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.User, config.Password, config.Host, config.Port, config.Database)

	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatalf("Could not open database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Could not ping database: %v", err)
	}

	log.Printf("Connected to database %s@%s:%d/%s\n", config.User, config.Host, config.Port, config.Database)
}

func Close() {
	db.Close()
	log.Println("Closed connection to database")
}

func (c connection_config) String() string {
	return fmt.Sprintf("user='%s', host:port='%s:%d', database='%s'", c.User, c.Host, c.Port, c.Database)
}
