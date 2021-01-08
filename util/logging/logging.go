package logging

import (
	"fmt"
	"github.com/logrusorgru/aurora"
)

func Error(message ...interface{}) {
	fmt.Printf("%s%s%s\n", aurora.Red(":: "), aurora.Bold("An error occurred: "), message)
}

func Warning(message ...interface{}) {
	fmt.Printf("%s%s%s\n", aurora.Yellow(":: "), aurora.Bold("Warning: "), fmt.Sprint(message...))
}

func Note(message ...interface{}) {
	fmt.Printf("%s%s%s\n", aurora.Magenta(":: "), aurora.Bold("Note: "), fmt.Sprint(message...))
}

func Success(message ...interface{}) {
	fmt.Printf("%s%s\n", aurora.Green(":: "), fmt.Sprint(message...))
}

func SuccessExclamation(message ...interface{}) {
	fmt.Printf("%s%s", aurora.Bold(aurora.Green("! ")), aurora.Bold(fmt.Sprint(message...)))
}