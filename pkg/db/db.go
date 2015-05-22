package db

import (
	"flag"
	"log"

	"github.com/boltdb/bolt"
)

// DB is used to access the DB
var DB Database

var flgDB string

// Database wraps a bolt.DB instance to add any future overrides we might need
type Database struct {
	*bolt.DB
}

func init() {
	flag.StringVar(&flgDB, "db", "meep.db", "path to meep bot db file")
}

func Open() {
	db, err := bolt.Open(flgDB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	DB = Database{db}
}
