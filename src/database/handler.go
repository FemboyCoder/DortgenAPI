package database

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"math/rand"
	"time"
)

var (
	database         *sql.DB
	GenerateCooldown int64
)

type ApiKey struct {
	ApiKey        string
	LastGenerated int64
	Created       int64
	Uses          int `json:"uses,omitempty"`
	Disabled      bool
	Owner         string `json:"owner,omitempty"`
	Notes         string
}

/*
Startup ~ Used to start up the database and create the apikey table if it doesn't exist
*/
func Startup(dataFolder string, generateCooldown int) error {

	GenerateCooldown = int64(generateCooldown)

	// open connection to the database
	databaseConnection, err := sql.Open("sqlite3", dataFolder+"/database.sqlite")
	if err != nil {
		return err
	}

	// create apikey table
	_, err = databaseConnection.Exec(`CREATE TABLE IF NOT EXISTS apikeys(
    								apikey TEXT NOT NULL PRIMARY KEY UNIQUE, -- api key for access
    								lastgenerated INTEGER NOT NULL DEFAULT 0, -- last time the key was used to generate a combo in unix millis
    								created INTEGER NOT NULL DEFAULT (strftime('%s', 'now')), -- when the key was created in unix millis
    								uses INTEGER NOT NULL DEFAULT 0, -- how many times the key has been used to generate a combo
    								disabled INTEGER NOT NULL DEFAULT 0, -- if the key is allowed to generate combos
    								owner TEXT NOT NULL DEFAULT 'unknown', -- who owns the key
    								notes TEXT NOT NULL DEFAULT '' -- notes about the key
    								);`)
	if err != nil {
		return err
	}

	database = databaseConnection

	return nil
}

/*
CreateApiKey ~ Returns the api key or an error if it fails to create the api key for some reason
*/
func CreateApiKey(owner string) (*ApiKey, error) {

	// check to see if user already has a key
	apikey, err := DoesUserExist(owner)
	if err != nil {
		return nil, err
	}
	if apikey != nil {
		return apikey, nil
	}

	// generate a random api key until it is unique
	key := generateRandomKey()
	for {
		apiKey, err := DoesKeyExist(key)
		if err != nil {
			return nil, err
		}
		if apiKey != nil {
			key = generateRandomKey()
		} else {
			break
		}
	}

	// add the new apikey to the database
	err = addUser(key, owner)
	if err != nil {
		return nil, err
	}
	return apikey, nil
}

/*
addUser ~ Used internally to add a user to
*/
func addUser(apikey string, owner string) error {
	_, err := database.Exec("INSERT INTO apikeys (apikey, owner) VALUES (?, ?)", apikey, owner)
	return err
}

/*
generateRandomKey ~ Used to generate a new random api key to put into the database
*/
func generateRandomKey() string {
	prefix := "dortgen-"
	keylength := 12
	randomString := ""
	for i := 0; i < keylength; i++ {
		randomString += string(rune(65 + rand.Intn(25)))
	}
	return prefix + randomString
}

/*
DoesKeyExist ~ Used to check if a key exists in the database
*/
func DoesKeyExist(key string) (*ApiKey, error) {
	// checks database to see if api key exists
	result, err := database.Query("SELECT * FROM apikeys WHERE apikey = ?", key)
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("result is nil")
	}

	var apiKey ApiKey
	if result.Next() {
		err := result.Scan(&apiKey.ApiKey, &apiKey.LastGenerated, &apiKey.Created, &apiKey.Uses, &apiKey.Disabled, &apiKey.Owner, &apiKey.Notes)
		if err != nil {
			return nil, err
		}
	}

	return &apiKey, nil
}

/*
DoesUserExist ~ Used to check if someone has a key or not already before creating a new one and returns their key if they have one
*/
func DoesUserExist(owner string) (*ApiKey, error) {
	// checks database to see if api key exists
	result, err := database.Query("SELECT * FROM apikeys WHERE owner = ?", owner)
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("result is nil")
	}

	if result.Next() {
		var apikey ApiKey
		err = result.Scan(&apikey.ApiKey, &apikey.LastGenerated, &apikey.Created, &apikey.Uses, &apikey.Disabled, &apikey.Owner, &apikey.Notes)
		if err != nil {
			return nil, err
		}
		return &apikey, nil
	} else {
		return nil, nil
	}
}

/*
IsKeyDisabled ~ Used to check if a key is disabled or not
*/
func IsKeyDisabled(apikey string) (*ApiKey, error) {
	result, err := database.Query("SELECT * FROM apikeys WHERE apikey = ?", apikey)
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, errors.New("result is nil")
	}
	if result.Next() {
		var apikey ApiKey
		err = result.Scan(&apikey.ApiKey, &apikey.LastGenerated, &apikey.Created, &apikey.Uses, &apikey.Disabled, &apikey.Owner, &apikey.Notes)
		if err != nil {
			return nil, err
		}
		return &apikey, nil
	} else {
		return nil, errors.New("key does not exist")
	}
}

/*
UpdateCooldown ~ Used to set the lastgenerated column for the api key
*/
func UpdateCooldown(apikey string) error {
	_, err := database.Exec("UPDATE apikeys SET lastgenerated = ? WHERE apikey = ?", time.Now().Unix(), apikey)
	return err
}
