package core

import (
	"agent301/helper"
	"fmt"
	"net/http"
)

func setHeader(http *http.Request, authToken string) {
	userAgent, os := generateRandomUserAgent()
	if userAgent == "" || os == "" {
		helper.PrettyLog("error", "Failed Generate Random User Agent")
		return
	}

	header := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "en-US,en;q=0.9,id;q=0.8",
		"content-type":       "application/json",
		"priority":           "u=1, i",
		"sec-ch-ua":          `"Android WebView";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`,
		"sec-ch-ua-platform": fmt.Sprintf("\"%s\"", os),
		"authorization":      authToken,
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-site",
		"Referrer":           "https://telegram.agent301.org/",
		"Referrer-Policy":    "strict-origin-when-cross-origin",
		"X-Requested-With":   "org.telegram.messenger.web",
		"User-Agent":         userAgent,
	}

	for key, value := range header {
		http.Header.Set(key, value)
	}
}
