package core

import (
	"agent301/helper"
	"flag"
	"fmt"
	"os"

	"github.com/gookit/config/v2"
)

var maxThread, minRandomSleep, maxRandomSleep, selectedTools int
var queryDataPath, qrTokenPath, proxyPath, outputPath, refId string
var isSpinWheel, isUseProxy bool

func initConfig() {
	maxThread = config.Int("MAX_THREAD")
	isUseProxy = config.Bool("USE_PROXY")
	minRandomSleep = config.Int("RANDOM_SLEEP.MIN")
	maxRandomSleep = config.Int("RANDOM_SLEEP.MAX")
	refId = config.String("REF_ID")
	isSpinWheel = config.Bool("AUTO_SPIN")

	queryDataPath = "./query.txt"
	qrTokenPath = "./qr_token.txt"
	proxyPath = "./proxy.txt"
	outputPath = "./output"

	if !helper.CheckFileOrFolder(queryDataPath) {
		file, _ := os.Create(queryDataPath)
		defer file.Close()
	}

	if !helper.CheckFileOrFolder(qrTokenPath) {
		file, _ := os.Create(qrTokenPath)
		defer file.Close()
	}

	if !helper.CheckFileOrFolder(proxyPath) {
		file, _ := os.Create(proxyPath)
		defer file.Close()
	}

	if !helper.CheckFileOrFolder(outputPath) {
		os.MkdirAll(outputPath, 0755)
	}
}

func LaunchBot() {
	initConfig()

	queryData := helper.ReadFileTxt(queryDataPath)
	if queryData == nil {
		helper.PrettyLog("error", "Query data not found")
		return
	}

	helper.PrettyLog("info", fmt.Sprintf("%v Query Data Detected", len(queryData)))

	flagArg := flag.Int("c", 0, "Input Choice With Flag -c, 1 = Auto Completing All Task (Unlimited Loop Without Proxy),  2 = Auto Completing All Task (Unlimited Loop With Proxy), 3 = Get Qr Token, 4 = Qr Farming, 5 = Connect Wallet (Development Stage)")

	flag.Parse()

	if *flagArg > 5 {
		helper.PrettyLog("error", "Invalid Flag Choice")
	} else if *flagArg != 0 {
		selectedTools = *flagArg
	}

	if selectedTools == 0 {
		helper.PrettyLog("1", "Auto Completing All Task")
		helper.PrettyLog("2", "Generate Qr Token")
		helper.PrettyLog("3", "Merge Qr Token")
		helper.PrettyLog("4", "Qr Farming")
		helper.PrettyLog("5", "Connect Wallet (Upcoming)")

		helper.PrettyLog("input", "Select Your Choice: ")

		_, err := fmt.Scan(&selectedTools)
		if err != nil {
			helper.PrettyLog("error", "Selection Invalid")
			return
		}
	}

	processSelectedTools(queryData)
}
