package core

import (
	"fmt"
	"net/http"
	"time"

	"agent301/helper"
)

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

func launchBot(username string, query string, apiUrl string, referUrl string, refId string, isSpinWheel bool) {
	client := &Client{
		apiURL:     apiUrl,
		referURL:   referUrl,
		authToken:  query,
		httpClient: &http.Client{},
	}

	// Get User Data
	req, err := client.getMe(refId)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get user data: %v", err))
		return
	}

	res, err := handleResponse(req)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return
	}

	userData := processResponse(res)

	req, err = client.getWheel()
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("Failed to get user data: %v", err))
		return
	}

	res, err = handleResponse(req)
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
	if tasks, exists := userData["tasks"].([]interface{}); exists {
		for _, task := range tasks {
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}

			if !taskMap["is_claimed"].(bool) {
				req, err = client.mainTask(taskMap["type"].(string))
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("Failed to completing main task: %v", err))
					return
				}

				res, err = handleResponse(req)
				if err != nil {
					fmt.Println("Error handling response:", err)
					return
				}

				taskData := processResponse(res)

				if taskData["is_completed"].(bool) {
					helper.PrettyLog("success", fmt.Sprintf("%s | Completed Task : %s | Reward: %.0f | Current Balance: %.0f", username, taskMap["type"].(string), taskData["reward"].(float64), taskData["balance"].(float64)))
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
				}
			}

			// Completing Video Task
			if key == "hour" {
				for key, value := range task.(map[string]interface{}) {
					if key == "timestamp" && (float64(time.Now().Unix()) >= value.(float64)) {
						var taskData map[string]interface{}

						for i := 1; i <= 5; i++ {
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

							helper.PrettyLog("success", fmt.Sprintf("%s | Completed Video Task | Sleep 15 Second...", username))

							time.Sleep(15 * time.Second)
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
			}
		}
	} else {
		helper.PrettyLog("error", fmt.Sprintf("%s | No tasks wheel found...", username))
	}

	if isSpinWheel {
		isLimit := false

		for !isLimit {
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

			spinner := processResponse(res)

			helper.PrettyLog("success", fmt.Sprintf("%s | Spinning Wheel | Reward: %s | Balance: %.0f | Toncoin: %.0f | Notcoin: %.0f | Ticket: %.0f", username, spinner["reward"].(string), spinner["balance"].(float64), spinner["toncoin"].(float64), spinner["notcoin"].(float64), spinner["tickets"].(float64)))

			helper.PrettyLog("info", fmt.Sprintf("%s | Sleep 15 Second Before Spinning Wheel Again...", username))

			time.Sleep(15 * time.Second)
		}
	}
}
