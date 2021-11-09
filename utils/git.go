package utils

import (
	"fmt"
	"github.com/go-git/go-git"
)

func getCurrentBranchName() (string, error) {
	open, err := git.PlainOpen(".") // Get working directory, is "." enough?
	if err != nil {
		return "", err
	}

	head, err := open.Head()
	if err != nil {
		return "", err
	}
	referenceName := head.Name()
	if !referenceName.IsBranch() {
		return "", fmt.Errorf("not branch")
	}

	return referenceName.String(), nil
}
