package terminal

import (
	"github.com/buger/goterm"
)

func UpdateStatus(operationName, status, message, link string) error {
	_, err := goterm.Println(goterm.Bold(operationName), ": ", status, " (", message, ") ", goterm.Color(link, goterm.BLUE))
	if err != nil {
		return err
	}

	goterm.Flush()
	goterm.MoveCursorUp(1)

	return nil
}

func UpdateFail(operationName, status, message, link string) error {
	_, err := goterm.Println(goterm.Bold(operationName), ": ", goterm.Color(status, goterm.RED), " (", message, ") ", goterm.Color(link, goterm.BLUE))
	if err != nil {
		return err
	}

	goterm.Flush()
	goterm.MoveCursorUp(1)
	return nil
}

func UpdateUnrecoverable(operationName, message, link string) error {
	_, err := goterm.Println(goterm.Bold(operationName), ": ", message, " - ", goterm.Color(link, goterm.BLUE))
	if err != nil {
		return err
	}
	goterm.Flush()
	goterm.MoveCursorUp(1)
	return nil
}

func StartingRun(operationName string) error {
	_, err := goterm.Println(operationName, "...")
	if err != nil {
		return err
	}
	goterm.Flush()
	goterm.MoveCursorUp(1)

	return nil
}
