package terminal

import (
	"github.com/buger/goterm"
)

func UpdateStatus(operationName, status, message, link string) error {
	return replaceLine("%s: %s (%s) %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE))
}

func UpdateFail(operationName, status, message, link string) error {
	return replaceLine("%s: %s (%s) %s", goterm.Bold(operationName), goterm.Color(status, goterm.RED), message, goterm.Color(link, goterm.BLUE))
}

func UpdateUnrecoverable(operationName, message, link string) error {
	return replaceLine("%s\n%s\nLink: %s", goterm.Bold(operationName), message, goterm.Color(link, goterm.BLUE))
}

func DoneMessage(operationName, message, link string) error {
	if message == "" {
		return nil
	}
	return replaceLine("%s\n%s\nLink: %s", goterm.Bold(operationName), message, goterm.Color(link, goterm.BLUE))
}

func StartingRun(operationName string) error {
	return replaceLine("%s...", operationName)
}

func replaceLine(format string, a ...interface{}) error {
	_, err := goterm.Print(goterm.ResetLine(""))
	if err != nil {
		return err
	}
	_, err = goterm.Printf(format, a...)
	if err != nil {
		return err
	}

	goterm.Flush()

	return nil
}
