package terminal

import (
	"fmt"
	"github.com/buger/goterm"
)

func UpdateStatus(operationName, status, message, link string) error {
	msg := ""
	if message != "" {
		msg = fmt.Sprintf("%s: %s (%s) %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE))
	} else {
		msg = fmt.Sprintf("%s: %s %s", goterm.Bold(operationName), status, goterm.Color(link, goterm.BLUE))
	}
	return replaceLine(msg)
}

func UpdateFail(operationName, status, message, link string) error {
	msg := ""
	if message != "" {
		msg = fmt.Sprintf("%s: %s (%s) %s", goterm.Bold(operationName), goterm.Color(status, goterm.RED), message, link)
	} else {
		msg = fmt.Sprintf("%s: %s %s", goterm.Bold(operationName), goterm.Color(status, goterm.RED), link)
	}
	return replaceLine(msg)
}

func UpdateUnrecoverable(operationName, message, link string) error {
	return replaceLine("%s\n%s\nLink: %s", goterm.Bold(operationName), message, link)
}

func DoneMessage(operationName, message, link string) error {
	if message == "" {
		return nil
	}
	return replaceLine("%s\n%s\nLink: %s", goterm.Bold(operationName), message, link)
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
