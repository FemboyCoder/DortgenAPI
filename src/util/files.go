package util

import "os"

/*
CreateFileIfNotExists ~ creates a file at the specified path if it does not exist.
*/
func CreateFileIfNotExists(path string) error {
	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if not, create it
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
	}
	return nil
}

/*
CreateFolderIfNotExists ~ creates a folder at the specified path if it does not exist.
*/
func CreateFolderIfNotExists(path string) error {
	// check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if not, create it
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}
