package main

import (
	"embed"
	"errors"
	"os"
)

//go:embed templates
var templateFS embed.FS

func copyFileFromTemplate(src, dest string) error {
	// check to ensure destination file does not already exist
	if fileExists(dest) {
		return errors.New(dest + " already exists!")
	}

	//read data of src files from templateFS
	data, err := templateFS.ReadFile(src)
	if err != nil {
		exitGracefully(err)
	}
	//copy data to dist file
	err = copyDataToFile(data, dest)
	if err != nil {
		exitGracefully(err)
	}

	return nil
}

func copyDataToFile(data []byte, dest string) error {
	err := os.WriteFile(dest, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func fileExists(fileToCheck string) bool {
	if _, err := os.Stat(fileToCheck); os.IsNotExist(err) {
		return false
	}
	return true
}
