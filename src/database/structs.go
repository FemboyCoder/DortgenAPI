package database

import (
	"database/sql"
	"errors"
)

type DatabaseConnection struct {
	Database *sql.DB
}

func (databaseConnection *DatabaseConnection) GetKeyFromOwner(s string) (string, error) {
	// get api key from owner
	result, err := databaseConnection.Database.Query("SELECT apikey FROM apikeys WHERE owner = ?", s)
	if err != nil {
		return "", err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	var key string
	if result.Next() {
		err = result.Scan(&key)
		if err != nil {
			return "", err
		}
		return key, nil
	}
	return "", errors.New("key not found")
}

func (databaseConnection *DatabaseConnection) GetOwnerFromKey(key string) (string, error) {
	// get owner from api key
	result, err := databaseConnection.Database.Query("SELECT owner FROM apikeys WHERE apikey = ?", key)
	if err != nil {
		return "", err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	var owner string
	if result.Next() {
		err = result.Scan(&owner)
		if err != nil {
			return "", err
		}
		return owner, nil
	}
	return "", errors.New("key not found")
}

func (databaseConnection *DatabaseConnection) DoesOwnerExist(owner string) (bool, error) {
	// check if owner exists
	result, err := databaseConnection.Database.Query("SELECT owner FROM apikeys WHERE owner = ?", owner)
	if err != nil {
		return true, err
	}
	defer func(result *sql.Rows) {
		_ = result.Close()
	}(result)

	if result.Next() {
		return true, nil
	}
	return false, nil
}

type ApiKey struct {
	ApiKey        string
	LastGenerated int64
	Created       int64
	Uses          int `json:"uses,omitempty"`
	Disabled      bool
	Owner         string `json:"owner,omitempty"`
	Notes         string
}

type KeyCreator struct {
	keyLength int
}

type Alt struct {
	Id       int
	Email    string
	Password string
}
