package terminal

import (
	"fmt"
	"github.com/buger/goterm"
)

func UpdateStatus(operationName, status, message, link string, isTempLine bool) error {
	// TODO: remove this when goterm is working. doesn't work for me... not sure why
	fmt.Println(fmt.Printf("%s: %s (%s) - %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE)))

	_, err := goterm.Println("%s: %s (%s) - %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE))
	if err != nil {
		return err
	}

	if isTempLine {
		goterm.MoveCursorUp(1)
	}

	return nil
}
