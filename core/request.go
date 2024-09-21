package core

import (
	"agent301/helper"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c *Client) makeRequest(method string, endpoint string, jsonBody interface{}) ([]byte, error) {
	fullURL := "https://api.agent301.org" + endpoint

	// Convert body to JSON
	var reqBody []byte
	var err error
	if jsonBody != nil {
		reqBody, err = json.Marshal(jsonBody)
		if err != nil {
			return nil, err
		}
	}

	// Create new request
	req, err := http.NewRequest(method, fullURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, err
	}

	// Set header
	setHeader(req, c.authToken)

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Handle non-200 status code
	if resp.StatusCode >= 400 {
		// Read the response body to include in the error message
		bodyBytes, bodyErr := io.ReadAll(resp.Body)
		if bodyErr != nil {
			return nil, fmt.Errorf("error status: %v, and failed to read body: %v", resp.StatusCode, bodyErr)
		}
		return nil, fmt.Errorf("error status: %v, error message: %s", resp.StatusCode, string(bodyBytes))
	}

	return io.ReadAll(resp.Body)
}

// Get Account Detail
func (c *Client) getMe() map[string]interface{} {
	payload := map[string]string{
		"referrer_id": refId,
	}

	res, err := c.makeRequest("POST", "/getMe", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get me: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

func (c *Client) getWheel() map[string]interface{} {
	payload := map[string]string{}

	res, err := c.makeRequest("POST", "/wheel/load", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get wheel: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

// Get Main Task
func (c *Client) getMainTask() map[string]interface{} {
	payload := map[string]string{}

	res, err := c.makeRequest("POST", "/getTasks", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to get main task: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

// Completing Main Task
func (c *Client) mainTask(taskType string) map[string]interface{} {
	payload := map[string]string{
		"type": taskType,
	}

	res, err := c.makeRequest("POST", "/completeTask", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to completing main task: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

// Completing Wheel Task
func (c *Client) wheelTask(taskType string) map[string]interface{} {
	payload := map[string]string{
		"type": taskType,
	}

	res, err := c.makeRequest("POST", "/wheel/task", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to completing wheel task: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

// Spin Wheel
func (c *Client) spinWheel() map[string]interface{} {
	payload := map[string]string{}

	res, err := c.makeRequest("POST", "/wheel/spin", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to spin wheel: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}

// Qr Farming
func (c *Client) qrFarming(qrToken string) map[string]interface{} {
	payload := map[string]string{
		"token": qrToken,
	}

	res, err := c.makeRequest("POST", "/passQrToken", payload)
	if err != nil {
		helper.PrettyLog("error", fmt.Sprintf("| %s | Failed to scan qr: %v", c.account.Username, err))
		return nil
	}

	result, err := handleResponseMap(res)
	if err != nil {
		fmt.Println("Error handling response:", err)
		return nil
	}

	return result
}
