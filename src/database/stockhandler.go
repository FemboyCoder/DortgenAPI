package database

import (
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func (database *DatabaseConnection) GetStockAmount() (int, error) {
	// get stock amount
	var stock int
	err := database.Database.QueryRow("SELECT COUNT(*) FROM altlist").Scan(&stock)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func (database *DatabaseConnection) GetAltAndRemoveFromStock() (*Alt, error) {
	// get stock amount
	var alt Alt
	err := database.Database.QueryRow("SELECT * FROM altlist LIMIT 1").Scan(&alt.Id, &alt.Email, &alt.Password)
	if err != nil {
		return nil, err
	}
	_, err = database.Database.Exec("DELETE FROM altlist WHERE id = ?", alt.Id)
	if err != nil {
		return nil, err
	}
	return &alt, nil
}

func (database *DatabaseConnection) AddAltToStock(email string, password string) error {
	_, err := database.Database.Exec("INSERT INTO altlist (email, password) VALUES (?, ?)", email, password)
	return err
}

func (database *DatabaseConnection) AddAccountsFromFile(file io.Reader, fileSize int64) (string, error) {
	fileBuffer := make([]byte, fileSize)
	_, err := file.Read(fileBuffer)
	if err != nil {
		return "", err
	}
	fileData := string(fileBuffer)
	combos := strings.Split(fileData, "\n")
	total := 0
	success := 0
	for _, combo := range combos {
		combo = strings.ReplaceAll(combo, "\n", "")
		combo = strings.ReplaceAll(combo, "\r", "")
		combo = strings.ReplaceAll(combo, " ", "")
		alt := strings.Split(combo, ":")
		if len(alt) != 2 {
			continue
		}
		err = database.AddAltToStock(alt[0], alt[1])
		total++
		if err != nil {
			if strings.Contains(err.Error(), "UNIQUE constraint failed") {
				log.Println(" [!] Duplicate account:", alt[0], alt[1])
				continue
			}
			return "", err
		}
		log.Println(" [+] Added alt: " + alt[0] + ":" + alt[1])
		success++
	}
	return "Added " + strconv.Itoa(success) + " accounts out of " + strconv.Itoa(total), nil
}

func (database *DatabaseConnection) GetCooldown(key string) (int, error) {
	// returns time until cooldown is over
	var cooldown int
	err := database.Database.QueryRow("SELECT lastgenerated FROM apikeys WHERE apikey = ?", key).Scan(&cooldown)
	if err != nil {
		return 0, err
	}

	nextGen := cooldown + int(GenerateCooldown)
	if nextGen < int(time.Now().Unix()) {
		return 0, nil
	}
	return nextGen - int(time.Now().Unix()), nil
}

func (database *DatabaseConnection) SetCooldown(key string) error {
	_, err := database.Database.Exec("UPDATE apikeys SET lastgenerated = ? WHERE apikey = ?", time.Now().Unix(), key)
	return err
}
