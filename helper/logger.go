package helper

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

func PrettyLog(level, message string) {
	level = strings.ToUpper(level)

	var levelColor *color.Color
	switch level {
	case "INFO":
		levelColor = color.New(color.FgWhite) // Blue for INFO
	case "ERROR":
		levelColor = color.New(color.FgRed) // Red for ERROR
	case "WARNING":
		levelColor = color.New(color.FgYellow) // Yellow for WARNING
	case "INPUT":
		levelColor = color.New(color.FgCyan) // Cyan for INPUT
	case "SUCCESS":
		levelColor = color.New(color.FgGreen) // Cyan for INPUT
	default:
		levelColor = color.New(color.FgWhite) // White for default
	}

	// Print the log message with color
	if level == "INPUT" {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s", message)
	} else {
		levelColor.Printf("[%s] ", level)
		fmt.Printf("%s\n", message)
	}
}

func PrintLogo() {
	levelColor := color.New(color.FgYellow)

	levelColor.Println(`
  /$$$$$$                                  /$$      /$$$$$$   /$$$$$$    /$$  
 /$$__  $$                                | $$     /$$__  $$ /$$$_  $$ /$$$$  
| $$  \ $$  /$$$$$$   /$$$$$$  /$$$$$$$  /$$$$$$  |__/  \ $$| $$$$\ $$|_  $$  
| $$$$$$$$ /$$__  $$ /$$__  $$| $$__  $$|_  $$_/     /$$$$$/| $$ $$ $$  | $$  
| $$__  $$| $$  \ $$| $$$$$$$$| $$  \ $$  | $$      |___  $$| $$\ $$$$  | $$  
| $$  | $$| $$  | $$| $$_____/| $$  | $$  | $$ /$$ /$$  \ $$| $$ \ $$$  | $$  
| $$  | $$|  $$$$$$$|  $$$$$$$| $$  | $$  |  $$$$/|  $$$$$$/|  $$$$$$/ /$$$$$$
|__/  |__/ \____  $$ \_______/|__/  |__/   \___/   \______/  \______/ |______/
           /$$  \ $$                                                          
          |  $$$$$$/                                                          
           \______/                                                           
`)
	levelColor.Println("ρσωєяє∂ ву: ѕкιвι∂ι ѕιgмα ¢σ∂є")
}
