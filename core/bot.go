package core

import (
	"agent301/helper"
	"fmt"
	"net/http"
	"os"
	"time"
)

var qrToken []string

func processResponse(res map[string]interface{}) map[string]interface{} {
	var result map[string]interface{}
	// Mengakses data dari map
	if ok, exists := res["ok"].(bool); exists && ok {
		// Akses data "result" dari map
		if data, exists := res["result"].(map[string]interface{}); exists {
			result = data
		}
	} else {
		fmt.Println("Request failed or 'ok' is false")
	}

	return result
}

func getUserData(client *Client, refId string) map[string]interface{} {
	req, err := client.getMe(refId)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get user data: %v", err))
		return nil
	}

	res, err := handleResponse(req)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return processResponse(res)
}

func generateQrToken(token string, queryData []string) {
	// Check Is Qr Token Exits
	for _, existingToken := range qrToken {
		if existingToken == token {
			return
		}
	}

	qrToken = append(qrToken, token)

	// Save QrToken To Txt
	if qrToken != nil && len(qrToken) == len(queryData) {
		if helper.CheckFileOrFolder("./qr-token.txt") {
			os.Remove("./qr-token.txt")
			helper.PrettyLog("success", "QrToken Removed...")
		}

		for _, token := range qrToken {
			helper.SaveFileTxt("./qr-token.txt", token)
		}

		helper.PrettyLog("success", "QrToken Saved...")
	}
}

func completingMainTask(client *Client, username string, taskType string) {
	req, err := client.mainTask(taskType)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to completing main task: %v", err))
		return
	}

	res, err := handleResponse(req)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return
	}

	taskData := processResponse(res)

	if taskData["is_completed"].(bool) {
		helper.PrettyLog("success", fmt.Sprintf("%s | Completed Task : %s | Reward: %.0f | Current Balance: %.0f", username, taskType, taskData["reward"].(float64), taskData["balance"].(float64)))
	} else {
		helper.PrettyLog("error", fmt.Sprintf("%s | Failed Completed Task : %s | Try Again Letter...", username, taskType))
	}

	if taskType == "video" {
		helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 30s Before Completing Another Video Task...", username))

		time.Sleep(30 * time.Second)
	} else {
		helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", username))

		time.Sleep(5 * time.Second)
	}
}

func launchBot(username string, queryData []string, query string, apiUrl string, referUrl string, refId string, isSpinWheel bool, isQrFarming bool) {
	client := &Client{
		apiURL:     apiUrl,
		referURL:   referUrl,
		authToken:  query,
		httpClient: &http.Client{},
	}

	// Get User Data
	userData := getUserData(client, refId)

	if userData == nil {
		return
	}

	if isQrFarming {
		generateQrToken(userData["qr_token"].(string), queryData)

		tokens := helper.ReadFileTxt("./qr-token.txt")
		if tokens == nil {
			helper.PrettyLog("error", "Qr Token data not found")
			return
		}

		for _, token := range tokens {
			if token != userData["qr_token"].(string) {
				req, err := client.qrFarming(token)
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("%s | Failed to qr farming: %v", username, err))
					continue
				}
				res, err := handleResponse(req)

				qrFarming := processResponse(res)

				if reward, exits := qrFarming["reward"].(float64); exits {
					helper.PrettyLog("success", fmt.Sprintf("%s | Scan Qr Farming Successfully | Reward: %.0f | Sleep 15s Before Scan Another Qr...", username, reward))
				} else {
					helper.PrettyLog("error", fmt.Sprintf("%s | Scan Qr Farming Failed | Sleep 15s Before Scan Another Qr...", username))
				}

				time.Sleep(15 * time.Second)
			}
		}

		return
	}

	time.Sleep(5 * time.Hour)

	req, err := client.getWheel()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get user data: %v", err))
		return
	}

	res, err := handleResponse(req)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return
	}

	wheel := processResponse(res)

	// Merge Wheel To User Data
	for key, value := range wheel {
		if key == "notcoin" || key == "toncoin" {
			userData[key] = value
		}
	}

	helper.PrettyLog("success", fmt.Sprintf("%s | Balance: %.0f | Tickets: %.0f | Toncoin: %.0f | Notcoin: %.0f", username, userData["balance"].(float64), userData["tickets"].(float64), userData["toncoin"].(float64), userData["notcoin"].(float64)))

	// Completing Main Task
	req, err = client.getMainTask()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get main task: %v", err))
		return
	}

	res, err = handleResponse(req)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return
	}

	mainTask := processResponse(res)

	if tasks, exists := mainTask["data"].([]interface{}); exists {
		for _, task := range tasks {
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}

			if !taskMap["is_claimed"].(bool) {
				if taskMap["type"].(string) == "video" {
					limit := 0
					for int(taskMap["count"].(float64)) <= int(taskMap["max_count"].(float64)) {
						if limit == 40 {
							helper.PrettyLog("error", fmt.Sprintf("%s | Completing Video Task Limit Reached...", username))
							break
						}
						completingMainTask(client, username, taskMap["type"].(string))
						limit++
					}
				} else {
					completingMainTask(client, username, taskMap["type"].(string))
				}
			}
		}
	} else {
		fmt.Println("No tasks found")
	}

	// Completing Daily And Wheel Task
	if wheelTask, exits := wheel["tasks"].(map[string]interface{}); exits {
		for key, task := range wheelTask {
			// Completing Daily Task
			if key == "daily" {
				if (int64(task.(float64) + 86400)) < time.Now().Unix() {
					req, err = client.dailyTask()
					if err != nil {
						helper.PrettyLog("error", fmt.Sprintf("Failed to completing daily task: %v", err))
					}

					res, err = handleResponse(req)
					if err != nil {
						fmt.Println("Error handling response:", err)
					}

					taskData := processResponse(res)

					helper.PrettyLog("success", fmt.Sprintf("%s | Completed Daily Task | Current Ticket: %.0f", username, taskData["tickets"].(float64)))
					helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", username))

					time.Sleep(5 * time.Second)
				}
			}

			// Completing Video Task
			if key == "hour" {
				for key, value := range task.(map[string]interface{}) {
					if key == "timestamp" && (float64(time.Now().Unix()) >= value.(float64)) {
						var taskData map[string]interface{}

						for i := 1; i <= 7; i++ {
							req, err = client.videoTask()
							if err != nil {
								helper.PrettyLog("error", fmt.Sprintf("Failed to completing video task: %v", err))
								break
							}

							res, err = handleResponse(req)
							if err != nil {
								fmt.Println("Error handling response:", err)
							}

							taskData = processResponse(res)

							helper.PrettyLog("success", fmt.Sprintf("%s | Completed Video Task | Sleep 30 Second...", username))

							time.Sleep(30 * time.Second)
						}

						helper.PrettyLog("success", fmt.Sprintf("%s | Video Task Limit | Current Ticket: %.0f", username, taskData["tickets"].(float64)))
					}
				}
			}

			if key != "daily" && key != "hour" && !task.(bool) {
				req, err = client.wheelTask(key)
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("Failed to completing wheel task: %v", err))
				}

				res, err = handleResponse(req)
				if err != nil {
					fmt.Println("Error handling response:", err)
				}

				taskData := processResponse(res)

				helper.PrettyLog("success", fmt.Sprintf("%s | Completed Wheel Task : %s | Current Ticket: %.0f", username, key, taskData["tickets"].(float64)))
				helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 5s Before Completing Another Main Task...", username))

				time.Sleep(5 * time.Second)
			}
		}
	} else {
		helper.PrettyLog("error", fmt.Sprintf("%s | No tasks wheel found...", username))
	}

	if isSpinWheel {
		userData = getUserData(client, refId)
		if userData == nil {
			return
		}

		isLimit := false

		for !isLimit && (int(userData["tickets"].(float64)) > 0) {
			req, err = client.spinWheel()
			if err != nil {
				helper.PrettyLog("error", fmt.Sprintf("Failed to spin wheel: %v", err))
				isLimit = true
				break
			}

			res, err = handleResponse(req)
			if err != nil {
				fmt.Println("Error handling response:", err)
			}

			userData = processResponse(res)

			helper.PrettyLog("success", fmt.Sprintf("%s | Spinning Wheel | Reward: %s | Balance: %.0f | Toncoin: %.0f | Notcoin: %.0f | Ticket: %.0f", username, userData["reward"].(string), userData["balance"].(float64), userData["toncoin"].(float64), userData["notcoin"].(float64), userData["tickets"].(float64)))

			helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 15 Second Before Spinning Wheel Again...", username))

			time.Sleep(15 * time.Second)
		}
	}
}
