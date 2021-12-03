package terminal

import "github.com/buger/goterm"

func UpdateStatus(operationName, status, message, link string) error {
	_, err := goterm.Println("%s: %s (%s) - %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE))
	if err != nil {
		return err
	}
	goterm.MoveCursorUp(1)

	return nil
}
