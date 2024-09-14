package helper

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"
)

func ReadFileTxt(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		PrettyLog("error", fmt.Sprintf("Failed to read file txt: %v", err))
		return nil
	}
	defer file.Close()

	var value []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		value = append(value, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		PrettyLog("error", fmt.Sprintf("Error reading file: %v", err))
	}

	return value
}

func SaveFileTxt(filePath string, data string) error {
	// Cek apakah file sudah ada
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// Jika file tidak ada, tulis data baru
		err = os.WriteFile(filePath, []byte(data+"\n"), 0644)
	} else {
		// Jika file sudah ada, tambahkan data ke akhir file
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = f.WriteString(data + "\n")
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func CheckFileOrFolder(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func RandomNumber(min int, max int) int {
	return rand.Intn(max-min) + min
}

func FindKeyByValue(m map[string]interface{}, value interface{}) []string {
	var key []string
	for k, v := range m {
		if v == value {
			key = append(key, k)
		}
	}
	return key
}

func InputTerminal(prompt string) string {
	PrettyLog("input", prompt)

	reader := bufio.NewReader(os.Stdin)

	value, _ := reader.ReadString('\n')

	return strings.TrimSpace(value)
}
