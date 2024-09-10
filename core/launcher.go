package core

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
	"time"

	"agent301/helper"
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

func ProcessAccount(thread int, queryPath string, apiUrl string, referUrl string, refId string) {
	queryData := readQueryData(queryPath)
	if queryData == nil {
		helper.PrettyLog("error", "Query data not found")
		return
	}

	jobs := make(chan int)
	var wg sync.WaitGroup

	if len(queryData) < thread {
		thread = len(queryData)
	}

	// Memulai beberapa worker goroutine
	for i := 1; i <= thread; i++ {
		wg.Add(1)
		go worker(jobs, queryData, apiUrl, referUrl, refId, &wg)
	}

	// Mengirim pekerjaan baru secara terus menerus
	go func() {
		for {
			for index := range queryData {
				jobs <- index                      // Kirimkan pekerjaan ke worker
				time.Sleep(time.Millisecond * 500) // Simulasi penundaan antara pengiriman pekerjaan
			}
			time.Sleep(time.Second * 3) // Simulasi delay sebelum mengulang
		}
	}()

	// Menutup channel jobs dan menunggu semua worker selesai
	go func() {
		wg.Wait()   // Tunggu sampai semua worker selesai
		close(jobs) // Tutup channel jobs setelah semua pekerjaan selesai
	}()

	// Program utama berjalan terus menerus
	select {} // Block forever to keep the program running
}
