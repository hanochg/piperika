package terminal

import (
	"fmt"
	"github.com/buger/goterm"
	"strings"
	"time"
)

const (
	breakLine         = "\n↳ "
	animationInterval = time.Millisecond * 100
)

var (
	progressChan chan struct{}
	animation    = [...]string{"🕛", "🕐", "🕑", "🕒", "🕓", "🕔", "🕕", "🕖", "🕗", "🕘", "🕙", "🕚"}
)

func UpdateStatus(operationName, status, message, link string) {
	msg := ""
	if message != "" {
		msg = fmt.Sprintf("%s: %s (%s) %s", goterm.Bold(operationName), status, message, goterm.Color(link, goterm.BLUE))
	} else {
		msg = fmt.Sprintf("%s: %s %s", goterm.Bold(operationName), status, goterm.Color(link, goterm.BLUE))
	}
	progressLine(msg)
}

func UpdateFail(operationName, status, message, link string) {
	msg := ""
	if message != "" {
		msg = fmt.Sprintf("%s: %s (%s) %s", goterm.Bold(operationName), goterm.Color(status, goterm.RED), message, link)
	} else {
		msg = fmt.Sprintf("%s: %s %s", goterm.Bold(operationName), goterm.Color(status, goterm.RED), link)
	}
	progressLine(msg)
}

func UpdateUnrecoverable(operationName, message, link string) error {
	stopProcess()
	message = strings.ReplaceAll(message, "\n", breakLine)
	return replaceLine("💩 %s%s%s%s%s", goterm.Bold(operationName), breakLine, message, getOptionalLinkText(link))
}

func DoneMessage(operationName, message, link string) error {
	stopProcess()
	message = strings.ReplaceAll(message, "\n", breakLine)
	return replaceLine("✅ %s%s%s%s", goterm.Bold(operationName), breakLine, message, getOptionalLinkText(link))
}

func StartingRun(operationName string) error {
	_, err := goterm.Println("")
	if err != nil {
		return err
	}

	progressLine("%s...", operationName)
	return nil
}

func getOptionalLinkText(link string) string {
	linkText := ""
	if link != "" {
		linkText = fmt.Sprintf("%sLink: %s", breakLine, goterm.Color(link, goterm.BLUE))
	}
	return linkText
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

func progressLine(format string, a ...interface{}) {
	stopProcess()

	progressChan = make(chan struct{})
	go func() {
		for i := 0; ; i = (i + 1) % len(animation) {
			select {
			case <-time.After(animationInterval):
				err := replaceLine(animation[i]+format, a...)
				if err != nil {
					fmt.Printf("Error printing %v", err)
					return
				}
			case <-progressChan:
				return
			}
		}
	}()
}

func stopProcess() {
	if progressChan == nil {
		return
	}

	progressChan <- struct{}{}
	progressChan = nil
}
