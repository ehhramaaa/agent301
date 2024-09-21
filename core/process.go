package core

import (
	"agent301/helper"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func getAccountFromQuery(account *Account) {
	// Parsing Query To Get Username
	value, err := url.ParseQuery(account.QueryData)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to parse query: %v", err.Error()))
		return
	}

	if len(value.Get("query_id")) > 0 {
		account.QueryId = value.Get("query_id")
	}

	if len(value.Get("auth_date")) > 0 {
		account.AuthDate = value.Get("auth_date")
	}

	if len(value.Get("hash")) > 0 {
		account.Hash = value.Get("hash")
	}

	userParam := value.Get("user")

	// Mendekode string JSON
	var userData map[string]interface{}
	err = json.Unmarshal([]byte(userParam), &userData)
	if err != nil {
		panic(err)
	}

	// Mengambil ID dan username dari hasil decode
	userIDFloat, ok := userData["id"].(float64)
	if !ok {
		helper.PrettyLog("error", "Failed to convert ID to float64")
		return
	}

	account.UserId = int(userIDFloat)

	// Ambil username
	username, ok := userData["username"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get username")
		return
	}
	account.Username = username

	// Ambil first name
	firstName, ok := userData["first_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get first_name")
		return
	}
	account.FirstName = firstName

	// Ambil first name
	lastName, ok := userData["last_name"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get last_name")
		return
	}
	account.LastName = lastName

	// Ambil language code
	languageCode, ok := userData["language_code"].(string)
	if !ok {
		helper.PrettyLog("error", "Failed to get language_code")
		return
	}
	account.LanguageCode = languageCode

	// Ambil allowWriteToPm
	allowWriteToPm, ok := userData["allows_write_to_pm"].(bool)
	if !ok {
		helper.PrettyLog("error", "Failed to get allows_write_to_pm")
		return
	}
	account.AllowWriteToPm = allowWriteToPm
}

func processSelectedTools(queryData []string) {
	switch selectedTools {
	case 1:
		fmt.Println("<=====================[Auto Complete Task & Auto Play Game]=====================>")
	case 2:
		fmt.Println("<=====================[Get Qr Token]=====================>")
	case 3:
		fmt.Println("<=====================[Merge Qr Token]=====================>")
	case 4:
		fmt.Println("<=====================[Qr Farming]=====================>")
	case 5:
		helper.PrettyLog("info", "Coming Soon...")
		os.Exit(0)
	}

	if len(queryData) < maxThread {
		maxThread = len(queryData)
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, maxThread)

	if selectedTools == 1 {
		for {
			for index, query := range queryData {
				wg.Add(1)
				go processAccount(&semaphore, &wg, index, query)
			}
			wg.Wait()
		}
	} else {
		for index, query := range queryData {
			wg.Add(1)
			go processAccount(&semaphore, &wg, index, query)
		}

		wg.Wait()
	}
}

func processAccount(semaphore *chan struct{}, wg *sync.WaitGroup, index int, query string) {
	defer wg.Done()
	*semaphore <- struct{}{}

	account := &Account{QueryData: query}
	getAccountFromQuery(account)

	client := &Client{
		account:    account,
		httpClient: &http.Client{},
	}

	helper.PrettyLog("info", fmt.Sprintf("| %s | Start Launch Bot...", client.account.Username))

	client.authToken = query

	if isUseProxy {
		proxyList := helper.ReadFileTxt(proxyPath)

		if proxyList == nil {
			helper.PrettyLog("error", "Proxy Data Not Found")
			return
		}

		helper.PrettyLog("info", fmt.Sprintf("%v Proxy Detected", len(proxyList)))

		proxy := proxyList[index%len(proxyList)]

		client.proxy = proxy
	}

	client.selectProcess()

	if selectedTools == 1 {
		<-*semaphore

		randomSleep := helper.RandomNumber(minRandomSleep, maxRandomSleep)

		helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished, Sleep %v Before Next Lap...", client.account.Username, randomSleep))

		time.Sleep(time.Duration(randomSleep) * time.Second)
	} else {
		helper.PrettyLog("info", fmt.Sprintf("| %s | Launch Bot Finished, Sleep 5s Before Next Account...", client.account.Username))
		time.Sleep(5 * time.Second)
		<-*semaphore
	}
}

func (c *Client) selectProcess() {
	switch selectedTools {
	case 1:
		c.processMainTools()
	case 2:
		c.processGetQrToken()
	case 3:
		processMergeQrToken()
	case 4:
		c.processQrFarming()
	}
}

func (c *Client) processMainTools() {
	userData := c.getMe()
	if userData == nil {
		return
	}

	if data, exists := userData["result"].(map[string]interface{}); exists && len(data) > 0 {
		userData = data
	} else {
		helper.PrettyLog("error", "User Data Nil")
	}

	wheel := c.getWheel()
	if wheel == nil {
		return
	}

	if data, exists := wheel["result"].(map[string]interface{}); exists && len(data) > 0 {
		wheel = data
	}

	for key, value := range wheel {
		if key == "notcoin" || key == "toncoin" {
			userData[key] = value
		}
	}

	helper.PrettyLog("success", fmt.Sprintf("%s | Balance: %.0f | Tickets: %.0f | Toncoin: %.0f | Notcoin: %.0f", c.account.Username, userData["balance"].(float64), userData["tickets"].(float64), userData["toncoin"].(float64), userData["notcoin"].(float64)))

	mainTask := c.getMainTask()
	if mainTask == nil {
		return
	}

	if data, exists := mainTask["result"].(map[string]interface{}); exists && len(data) > 0 {
		mainTask = data
	}

	if tasks, exists := mainTask["data"].([]interface{}); exists {
		for _, task := range tasks {
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}

			if !taskMap["is_claimed"].(bool) {
				if taskMap["type"].(string) == "video" {

					var count, maxCount int = int(taskMap["count"].(float64)), int(taskMap["max_count"].(float64))

					for count < maxCount {
						result := c.mainTask(taskMap["type"].(string))
						if result == nil {
							count++
							continue
						}

						if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
							result = data
						}

						if result["is_completed"].(bool) {
							helper.PrettyLog("success", fmt.Sprintf("%s | Completed Task : %s | Reward: %.0f | Current Balance: %.0f", c.account.Username, taskMap["type"].(string), result["reward"].(float64), result["balance"].(float64)))
						} else {
							helper.PrettyLog("error", fmt.Sprintf("%s | Failed Completed Task : %s | Try Again Letter...", c.account.Username, taskMap["type"].(string)))
						}

						count++

						helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 30s Before Completing Another Video Task...", c.account.Username))

						time.Sleep(30 * time.Second)
					}

					helper.PrettyLog("error", fmt.Sprintf("%s | Completing Video Task Limit Reached...", c.account.Username))

				} else {
					result := c.mainTask(taskMap["type"].(string))
					if result == nil {
						continue
					}

					if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
						result = data
					}

					if result["is_completed"].(bool) {
						helper.PrettyLog("success", fmt.Sprintf("%s | Completed Task : %s | Reward: %.0f | Current Balance: %.0f", c.account.Username, taskMap["type"].(string), result["reward"].(float64), result["balance"].(float64)))
					} else {
						helper.PrettyLog("error", fmt.Sprintf("%s | Failed Completed Task : %s | Try Again Letter...", c.account.Username, taskMap["type"].(string)))
					}

					helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", c.account.Username))

					time.Sleep(5 * time.Second)
				}
			}
		}
	} else {
		helper.PrettyLog("warning", fmt.Sprintf("%s | No tasks found...", c.account.Username))
	}

	// Completing Daily And Wheel Task
	if wheelTask, exits := wheel["tasks"].(map[string]interface{}); exits {
		for key, task := range wheelTask {
			// Completing Daily Task
			if key == "daily" {
				if (int64(task.(float64) + 86400)) < time.Now().Unix() {
					result := c.wheelTask("daily")
					if result == nil {
						continue
					}

					if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
						result = data
					}

					helper.PrettyLog("success", fmt.Sprintf("%s | Completed Daily Task | Current Ticket: %.0f", c.account.Username, result["tickets"].(float64)))
					helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", c.account.Username))

					time.Sleep(5 * time.Second)
				}
			}

			// Completing Video Task
			if key == "hour" {
				for key, value := range task.(map[string]interface{}) {
					var count int
					if key != "timestamp" {
						count = int(value.(float64))
					}

					if key == "timestamp" && (float64(time.Now().Unix()) >= value.(float64)) {
						var taskData map[string]interface{}

						for count <= 5 {
							result := c.wheelTask("hour")
							if result == nil {
								count++
								continue
							}

							if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
								taskData = data
							}

							helper.PrettyLog("success", fmt.Sprintf("%s | Completed Video Task | Sleep 30 Second...", c.account.Username))

							count++

							time.Sleep(30 * time.Second)
						}

						helper.PrettyLog("success", fmt.Sprintf("%s | Video Task Limit | Current Ticket: %.0f", c.account.Username, taskData["tickets"].(float64)))
					}
				}
			}

			if key != "daily" && key != "hour" && !task.(bool) {
				result := c.wheelTask(key)
				if result == nil {
					continue
				}

				if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
					result = data
				}

				helper.PrettyLog("success", fmt.Sprintf("%s | Completed Wheel Task : %s | Current Ticket: %.0f", c.account.Username, key, result["tickets"].(float64)))
				helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", c.account.Username))

				time.Sleep(5 * time.Second)
			}
		}
	} else {
		helper.PrettyLog("error", fmt.Sprintf("%s | No tasks wheel found...", c.account.Username))
	}

	if isSpinWheel {
		userData = c.getMe()
		if userData == nil {
			return
		}

		if data, exists := userData["result"].(map[string]interface{}); exists && len(data) > 0 {
			userData = data
		}

		isLimit := false

		for !isLimit && (int(userData["tickets"].(float64)) > 0) {
			result := c.spinWheel()
			if result == nil {
				continue
			}

			if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
				result = data
			}

			helper.PrettyLog("success", fmt.Sprintf("%s | Spinning Wheel | Reward: %s | Balance: %.0f | Toncoin: %.0f | Notcoin: %.0f | Ticket: %.0f", c.account.Username, result["reward"].(string), result["balance"].(float64), result["toncoin"].(float64), result["notcoin"].(float64), result["tickets"].(float64)))

			helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 15 Second Before Spinning Wheel Again...", c.account.Username))

			time.Sleep(15 * time.Second)
		}
	}
}

func (c *Client) processGetQrToken() {
	filePath := fmt.Sprintf("%s/qr_token_%s.txt", outputPath, c.account.Username)

	if helper.CheckFileOrFolder(filePath) {
		helper.PrettyLog("success", fmt.Sprintf("| %s | QrToken Found %s...", c.account.Username, filePath))
		return
	}

	userData := c.getMe()
	if userData == nil {
		return
	}

	if data, exists := userData["result"].(map[string]interface{}); exists && len(data) > 0 {
		userData = data
	}

	token := userData["qr_token"].(string)

	time.Sleep(1 * time.Second)

	helper.SaveFileTxt(filePath, token)

	time.Sleep(1 * time.Second)

	if helper.CheckFileOrFolder(filePath) {
		helper.PrettyLog("success", fmt.Sprintf("| %s | QrToken Saved In %s...", c.account.Username, filePath))
	} else {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Save QrToken Failed...", c.account.Username))
	}
}

func (c *Client) processQrFarming() {
	userData := c.getMe()
	if userData == nil {
		return
	}

	if data, exists := userData["result"].(map[string]interface{}); exists && len(data) > 0 {
		userData = data
	}

	tokens := helper.ReadFileTxt(qrTokenPath)
	if tokens == nil {
		helper.PrettyLog("error", "Qr Token data not found, Please Insert Your Qr Token In qr_token.txt")
		return
	}

	for _, token := range tokens {
		if token != userData["qr_token"].(string) {
			result := c.qrFarming(token)
			if result == nil {
				continue
			}

			if data, exists := result["result"].(map[string]interface{}); exists && len(data) > 0 {
				result = data
			}

			if reward, exits := result["reward"].(float64); exits {
				helper.PrettyLog("success", fmt.Sprintf("%s | Scan Qr Farming Successfully | Reward: %.0f | Sleep 15s Before Scan Another Qr...", c.account.Username, reward))
				time.Sleep(15 * time.Second)
			} else {
				helper.PrettyLog("error", fmt.Sprintf("%s | Scan Qr Farming Failed | Sleep 15s Before Scan Another Qr...", c.account.Username))
				time.Sleep(15 * time.Second)
			}

		}
		helper.PrettyLog("success", fmt.Sprintf("%s | Scan Qr All Account Successfully...", c.account.Username))
	}
}

func processMergeQrToken() {
	files, err := os.ReadDir(outputPath)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to read directory %s: %v", outputPath, err))
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Qr Token File Detected", len(files)))

	var mergedData []string

	for _, file := range files {
		// Baca data dari file
		account := helper.ReadFileTxt(fmt.Sprintf("%s/%s", outputPath, file.Name()))

		for _, value := range account {
			mergedData = append(mergedData, value)
		}
	}

	// Check Folder
	mergePath := fmt.Sprintf("%s/merge/", outputPath)
	fileName := fmt.Sprintf("merge_qr_token_%s.txt", time.Now().Format("20060102150405"))

	if !helper.CheckFileOrFolder(mergePath) {
		os.Mkdir(fmt.Sprintf(mergePath), 0755)
	}

	// Save Query Data To Txt
	for _, value := range mergedData {
		err := helper.SaveFileTxt(mergePath+fileName, value)
		if err != nil {
			helper.PrettyLog("error", fmt.Sprintf("Error saving file: %v", err))
		}
	}

	helper.PrettyLog("success", fmt.Sprintf("Merge Query Data Successfully Saved In: %s", mergePath+fileName))
}
