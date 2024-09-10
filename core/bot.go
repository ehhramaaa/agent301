package core

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"agent301/helper"

)

// TODO Add Sleep
func worker(jobs <-chan int, query []string, apiUrl string, referUrl string, refId string, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := range jobs {
		username := getUsernameFromQuery(query[j])
		helper.PrettyLog("info", fmt.Sprintf("User: %s Started Bot...", username))
		processBot(username, query[j], apiUrl, referUrl, refId)
	}
}

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

func processBot(username string, query string, apiUrl string, referUrl string, refId string) {
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
	if wheelTask, exits := wheel["tasks"].([]interface{}); exits {
		for _, task := range wheelTask {
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}

			// Completing Daily Task
			if (taskMap["daily"].(int64) + 86400) < time.Now().Unix() {
				req, err = client.dailyTask()
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("Failed to completing daily task: %v", err))
					return
				}

				res, err = handleResponse(req)
				if err != nil {
					fmt.Println("Error handling response:", err)
					return
				}

				taskData := processResponse(res)

				helper.PrettyLog("success", fmt.Sprintf("%s | Completed Daily Task | Current Ticket: %.0f", username, taskData["tickets"].(float64)))
			}

			taskNotCompleted := helper.FindKeyByValue(taskMap, false)

			for _, task := range taskNotCompleted {
				req, err = client.wheelTask(task)
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("Failed to completing wheel task: %v", err))
					return
				}

				res, err = handleResponse(req)
				if err != nil {
					fmt.Println("Error handling response:", err)
					return
				}

				taskData := processResponse(res)

				helper.PrettyLog("success", fmt.Sprintf("%s | Completed Wheel Task : %s | Current Ticket: %.0f", username, taskData[fmt.Sprintf("tasks.%v", task)].(string), taskData["tickets"].(float64)))
			}
		}
	}

	// Completing Video Task
	if videoTask, exits := wheel["tasks.hour"].([]interface{}); exits {
		for _, task := range videoTask { 
			var taskData map[string]interface{}
			taskMap, ok := task.(map[string]interface{})
			if !ok {
				continue
			}
			
			for taskMap["count"].(int) != 5 {
				req, err = client.videoTask()
				if err != nil {
					helper.PrettyLog("error", fmt.Sprintf("Failed to completing video task: %v", err))
					return
				}
		
				res, err = handleResponse(req)
				if err != nil {
					fmt.Println("Error handling response:", err)
					return
				}
		
				taskData = processResponse(res)
		
				if taskData["hour.count"] == 5 {
					break
				}
		
				sleep := 15
		
				helper.PrettyLog("success", fmt.Sprintf("%s | Completed Video Task | Sleep: %v Second...", username, sleep))
		
				time.Sleep(time.Duration(sleep) * time.Second)
			}
		
			helper.PrettyLog("success", fmt.Sprintf("%s | Video Task Limit | Current Ticket: %.0f", username, taskData["tickets"].(float64)))
		}
	
	}
}
