package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

)

type Client struct {
	apiURL     string
	referURL   string
	authToken  string
	httpClient *http.Client
}

func handleResponse(respBody []byte) (map[string]interface{}, error) {
	// Mengurai JSON ke dalam map[string]interface{}
	var result map[string]interface{}
	err := json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return result, nil
}

func (c *Client) makeRequest(method string, endpoint string, jsonBody interface{}) ([]byte, error) {
	fullURL := c.apiURL + endpoint

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
	setHeader(req, c.referURL, c.authToken)

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
func (c *Client) getMe(refId string) ([]byte, error) {
	payload := map[string]string{
		"referrer_id": refId,
	}

	return c.makeRequest("POST", "/getMe", payload)
}

func (c *Client) getWheel() ([]byte, error) {
	payload := map[string]string{}

	return c.makeRequest("POST", "/wheel/load", payload)
}

// Get Main Task
func (c *Client) getMainTask() ([]byte, error){
	payload := map[string]string{}

	return c.makeRequest("POST", "/getTasks", payload)
}

// Completing Main Task
func (c *Client) mainTask(taskType string) ([]byte, error) {
	payload := map[string]string{
		"type": taskType,
	}

	return c.makeRequest("POST", "/completeTask", payload)
}

// Completing Wheel Task
func (c *Client) wheelTask(taskType string) ([]byte, error) {
	payload := map[string]string{
		"type": taskType,
	}

	return c.makeRequest("POST", "/wheel/task", payload)
}

// Completing Video Task
func (c *Client) videoTask() ([]byte, error) {
	payload := map[string]string{
		"type": "hour",
	}

	return c.makeRequest("POST", "/wheel/task", payload)
}

// Completing Daily Task
func (c *Client) dailyTask() ([]byte, error) {
	payload := map[string]string{
		"type": "daily",
	}

	return c.makeRequest("POST", "/wheel/task", payload)
}

// Spin Wheel
func (c *Client) spinWheel() ([]byte, error) {
	payload := map[string]string{}

	return c.makeRequest("POST", "/wheel/spin", payload)
}
