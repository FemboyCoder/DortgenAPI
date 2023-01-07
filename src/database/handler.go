package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var (
	Connection       *DatabaseConnection
	GenerateCooldown int64
)

/*
Startup ~ Used to start up the database and create the apikey table if it doesn't exist
*/
func Startup(dataFolder string, generateCooldown int) error {
	GenerateCooldown = int64(generateCooldown)

	// open connection to the database
	conn, err := sql.Open("sqlite3", dataFolder+"/database.sqlite")
	if err != nil {
		return err
	}

	// set the database connection
	Connection = &DatabaseConnection{
		Database: conn,
	}

	// create apikey table
	err = Connection.CreateApiKeyTable()
	if err != nil {
		return err
	}
	err = Connection.CreateAltListTable()
	if err != nil {
		return err
	}
	key, err := Connection.CreateAdminUser()
	if err != nil {
		return err
	}
	log.Println("Created admin user with key:", key)

	return nil
}
