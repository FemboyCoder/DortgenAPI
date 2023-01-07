package database

/*
CreateApiKeyTable ~ Used to create the apikey table if it doesn't exist
*/
func (databaseConnection *DatabaseConnection) CreateApiKeyTable() error {
	database := databaseConnection.Database
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS apikeys(
    								apikey TEXT NOT NULL PRIMARY KEY UNIQUE, -- api key for access
    								lastgenerated INTEGER NOT NULL DEFAULT 0, -- last time the key was used to generate a combo in unix millis
    								created INTEGER NOT NULL DEFAULT (strftime('%s', 'now')), -- when the key was created in unix millis
    								uses INTEGER NOT NULL DEFAULT 0, -- how many times the key has been used to generate a combo
    								disabled INTEGER NOT NULL DEFAULT 0, -- if the key is allowed to generate combos
    								owner TEXT NOT NULL UNIQUE, -- who owns the key
    								notes TEXT NOT NULL DEFAULT '' -- notes about the key
    								);`)
	return err
}

func (databaseConnection *DatabaseConnection) CreateAdminUser() (string, error) {

	// check if admin user already exists and return their key if it does
	exists, err := databaseConnection.DoesUserExist("admin")
	if err != nil {
		return "", err
	}
	if exists {
		key, err := databaseConnection.GetKeyFromOwner("admin")
		if err != nil {
			return "", err
		}
		return key, nil
	}

	keyCreator := &KeyCreator{
		keyLength: 32,
	}

	key, err := keyCreator.generateUniqueKey()
	if err != nil {
		return "", err
	}

	database := databaseConnection.Database
	_, err = database.Exec(`INSERT INTO apikeys(owner, apikey) VALUES('admin', ?)`, key)
	return key, err
}

func (databaseConnection *DatabaseConnection) CreateAltListTable() error {
	database := databaseConnection.Database
	_, err := database.Exec(`CREATE TABLE IF NOT EXISTS altlist(
									id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE, -- id of the alt
									email TEXT NOT NULL UNIQUE, -- email of the alt
									password TEXT NOT NULL); -- password of the alt`,
	)
	return err
}
