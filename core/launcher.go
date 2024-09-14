package core

import (
	"agent301/helper"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/gookit/config/v2"
)

func readQueryData(queryPath string) []string {
	query := helper.ReadFileTxt(queryPath)

	if query == nil {
		helper.PrettyLog("error", "Query data not found")
		return nil
	}

	return query
}

func getUsernameFromQuery(account string) string {
	// Parsing Query To Get Username
	value, err := url.ParseQuery(account)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to parse query: %v", err.Error()))
		return ""
	}

	userParam := value.Get("user")

	// Mendekode string JSON
	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		panic(err)
	}

	// Mengambil username dari hasil decode
	return userData["username"].(string)
}

func ProcessBot(config *config.Config) {
	queryPath := config.String("query-file")
	apiUrl := config.String("bot.api-url")
	referUrl := config.String("bot.refer-url")
	refId := config.String("bot.ref-Id")
	isSpinWheel := config.Bool("auto-spin")
	maxThread := config.Int("max-thread")

	queryData := readQueryData(queryPath)
	if queryData == nil {
		helper.PrettyLog("error", "Query data not found")
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	helper.PrettyLog("1", "Generate Qr Token")
	helper.PrettyLog("2", "Qr Farming")
	helper.PrettyLog("3", "Auto Completing Task And Auto Spin Wheel (Unlimited Loop)")

	cFlag := flag.String("c", "", "Choice of action (1, 2 or 3)")

	flag.Parse()

	var choice string

	// Jika flag -c disediakan
	if *cFlag != "" {
		// Validasi apakah flag -c berada dalam rentang yang benar (1 atau 2)
		if *cFlag == "1" || *cFlag == "2" || *cFlag == "3" {
			choice = *cFlag
		} else {
			fmt.Println("Invalid flag value for -c. Only 1 or 3 are allowed.")
			return
		}
	} else {
		// TODO
		choice = helper.InputTerminal("Select Choice: ")
	}

	helper.PrettyLog("info", "Start Processing Account...")

	time.Sleep(3 * time.Second)

	var isGenerateToken bool
	var isQrFarming bool

	var wg sync.WaitGroup

	// Membuat semaphore dengan buffered channel
	// TODO Deadlock
	semaphore := make(chan struct{}, maxThread)

	for j, query := range queryData {
		wg.Add(1)

		// Goroutine untuk setiap job
		go func(index int, query string) {
			defer wg.Done()

			// Mengambil token dari semaphore sebelum menjalankan job
			semaphore <- struct{}{}

			username := getUsernameFromQuery(query)
			helper.PrettyLog("info", fmt.Sprintf("%s | Started Bot...", username))

			switch choice {
			case "1":
				isQrFarming = false
				isGenerateToken = true

				launchBot(username, queryData, query, apiUrl, referUrl, refId, isSpinWheel, isQrFarming, isGenerateToken)
			case "2":
				isQrFarming = true

				launchBot(username, queryData, query, apiUrl, referUrl, refId, isSpinWheel, isQrFarming, isGenerateToken)
			case "3":
				launchBot(username, queryData, query, apiUrl, referUrl, refId, isSpinWheel, isQrFarming, isGenerateToken)
			}

			// Melepaskan token dari semaphore
			<-semaphore

			if choice == "3" {
				// Sleep setelah job selesai
				randomSleep := helper.RandomNumber(config.Int("random-sleep.min"), config.Int("random-sleep.max"))

				helper.PrettyLog("info", fmt.Sprintf("%s | Launch Bot Finished, Sleeping for %v seconds..", username, randomSleep))

				time.Sleep(time.Duration(randomSleep) * time.Second)
			}
		}(j, query)
	}

	// Tunggu sampai semua worker selesai memproses pekerjaan
	wg.Wait()

	if choice == "2" {
		// Program utama berjalan terus menerus
		select {} // Block forever to keep the program running
	}
}
