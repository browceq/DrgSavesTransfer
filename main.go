package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if copySaves() && insertSaves() && copyConfig() {
		fmt.Println("Transfer completed successfully")
	}
	fmt.Scanln()
}

func copySaves() bool {

	savesPath := scanRow(0, "save")
	files, err := os.ReadDir(savesPath)
	if err != nil {
		fmt.Printf("Failed to read files in %s: %s\n", savesPath, err)
		return false
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fileName := file.Name()
		if !strings.HasSuffix(fileName, "_Player.sav") && !strings.HasSuffix(fileName, "_Backup_3.sav") {
			continue
		}

		filePath := filepath.Join(savesPath, fileName)
		copyPath := filepath.Join(".", fileName)

		if !copyFile(filePath, copyPath) {
			fmt.Printf("Failed to copy %s to %s\n", fileName, copyPath)
			return false
		}

	}
	return true
}

func insertSaves() bool {

	copyPath := scanRow(1, "save")
	newFiles, err := os.ReadDir(copyPath)
	if err != nil {
		fmt.Printf("Failed to read files in %s: %s\n", copyPath, err)
		return false
	}

	var playerId string
	for _, file := range newFiles {
		if file.IsDir() {
			continue
		}
		playerId = file.Name()[:17]
	}

	oldFiles, err := os.ReadDir(".")
	if err != nil {
		fmt.Printf("Failed to read files in %s: %s\n", copyPath, err)
	}
	for _, file := range oldFiles {

		fileName := file.Name()
		if !strings.HasSuffix(fileName, "_Player.sav") && !strings.HasSuffix(fileName, "_Backup_3.sav") {
			continue
		}

		newName := replaceId(fileName, playerId)
		oldPath := filepath.Join(".", fileName)
		newPath := filepath.Join(copyPath, newName)

		if !copyFile(oldPath, newPath) {
			fmt.Printf("Failed to copy %s to %s\n", oldPath, newPath)
			return false
		}

	}
	return true
}

func copyConfig() bool {
	if !copyFile(scanRow(0, "config"), scanRow(1, "config")) {
		return false
	}
	return true
}

func scanRow(row int, mode string) string {
	path := "paths.txt"

	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", path)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	currentRow := 0
	for scanner.Scan() {
		if currentRow != row {
			currentRow++
			continue
		}
		line := scanner.Text()
		if mode == "save" {
			return line + "\\FSD\\Saved\\SaveGames"
		}
		if mode == "config" {
			return line + "\\FSD\\Saved\\Config\\WindowsNoEditor\\GameUserSettings.ini"
		}
	}
	return ""
}

func copyFile(filePath, copyPath string) bool {
	sourceFile, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("Failed to open file: %s\n", filePath)
		return false
	}
	defer sourceFile.Close()

	copyFile, err := os.Create(copyPath)
	if err != nil {
		fmt.Printf("Failed to create copy of the file: %s\n", filePath)
		return false
	}
	defer copyFile.Close()

	_, err = io.Copy(copyFile, sourceFile)
	if err != nil {
		fmt.Printf("Failed to copy the file: %s\n", filePath)
		return false
	}
	return true
}

func replaceId(filePath string, newId string) string {
	filePath = filePath[17:]
	filePath = newId + filePath
	return filePath
}
