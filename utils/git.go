package utils

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"os"
)

func GetCurrentBranchName() (string, error) {
	repository, err := getLocalRepo()
	if err != nil {
		return "", err
	}

	head, err := repository.Head()
	if err != nil {
		return "", err
	}
	referenceName := head.Name()
	if !referenceName.IsBranch() {
		return "", fmt.Errorf("not branch")
	}

	return referenceName.Short(), nil
}

func GetLocalBranches() ([]string, error) {
	repository, err := getLocalRepo()
	if err != nil {
		return nil, err
	}
	branches, err := repository.Branches()

	res := make([]string, 0)
	err = branches.ForEach(func(reference *plumbing.Reference) error {
		res = append(res, reference.Name().Short())
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func getLocalRepo() (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repository, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	return repository, nil
}
