package command

import "fmt"

// TODO: add colors and emojis!!!!!

func logInfo(operationName, message string) {
	fmt.Println(fmt.Sprintf("[%s] INFO: %s", operationName, message))
}

func logWarning(operationName, message string) {
	fmt.Println(fmt.Sprintf("[%s] WARN: %s", operationName, message))
}

func logError(operationName, message string) {
	fmt.Println(fmt.Sprintf("[%s] ERROR: %s", operationName, message))
}
