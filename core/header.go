package core

import (
	"agent301/helper"
	"fmt"
	"net/http"
)

func generateChromeVersion() int {
	var chromeVersion []int

	for i := 110; i <= 127; i++ {
		chromeVersion = append(chromeVersion, i)
	}

	return chromeVersion[helper.RandomNumber(0, len(chromeVersion))]
}

func setHeader(http *http.Request, referUrl string, authToken string) {
	device := []string{
		"ios",
		"android",
	}

	browserVersion := generateChromeVersion()

	header := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9,id;q=0.8",
		"content-type":       "application/json",
		"priority":           "u=1, i",
		"sec-ch-ua":          fmt.Sprintf("\"Chromium\";v=\"%v\", \"Not;A=Brand\";v=\"24\", \"Google Chrome\";v=\"%v\"", browserVersion, browserVersion),
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": fmt.Sprintf("\"%v\"", device[helper.RandomNumber(0, (len(device)-1))]),
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"Referer":            referUrl,
		"Referrer-Policy":    "strict-origin-when-cross-origin",
	}

	if authToken != "" {
		header["authorization"] = authToken
	}

	for key, value := range header {
		http.Header.Set(key, value)
	}
}
