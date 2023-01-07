package database

import (
	"database/sql"
	"errors"
	"math/rand"
)

func (database *DatabaseConnection) CreateApiKey(user string, keyLength int) error {
	keyCreator := KeyCreator{
		keyLength: keyLength,
	}
	key, err := keyCreator.generateUniqueKey()
	if err != nil {
		return err
	}

	// insert the key into the database
	_, err = database.Database.Exec("INSERT INTO apikeys (apikey, owner) VALUES (?, ?)", key, user)
	return err
}

/*
generateRandomKey ~ Used to generate a new random api key to put into the database
*/
func (keyCreator *KeyCreator) generateRandomKey() string {
	prefix := "dortgen-"
	randomString := ""
	for i := 0; i < keyCreator.keyLength; i++ {
		randomString += string(rune(65 + rand.Intn(25)))
	}
	return prefix + randomString
}

/*
generateRandomKey ~ Used to generate a new random api key to put into the database
*/
func (keyCreator *KeyCreator) generateUniqueKey() (string, error) {
	// generate a random key
	key := keyCreator.generateRandomKey()

	// check if the key is unique
	exists, err := Connection.DoesKeyExist(key)
	if err != nil {
		return "", err
	}
	for exists {
		key = keyCreator.generateRandomKey()
		exists, err = Connection.DoesKeyExist(key)
		if err != nil {
			return "", err
		}
	}
	return key, nil
}

func (databaseConnection *DatabaseConnection) DoesKeyExist(key string) (bool, error) {
	// check if key is in database
	result, err := databaseConnection.Database.Query("SELECT * FROM apikeys WHERE apikey = ?", key)
	if err != nil {
		return true, err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	return result.Next(), nil
}

func (databaseConnection *DatabaseConnection) DoesUserExist(user string) (bool, error) {
	// check if key is in database
	result, err := databaseConnection.Database.Query("SELECT * FROM apikeys WHERE owner = ?", user)
	if err != nil {
		return true, err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	return result.Next(), nil
}

func (databaseConnection *DatabaseConnection) IsKeyDisabled(key string) (bool, error) {
	// check if key is in database
	result, err := databaseConnection.Database.Query("SELECT disabled FROM apikeys WHERE apikey = ?", key)
	if err != nil {
		return true, err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	var disabled bool
	if result.Next() {
		err = result.Scan(&disabled)
		if err != nil {
			return true, err
		}
		return disabled, nil
	}
	return true, errors.New("key not found")
}
