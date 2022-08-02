package utils

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"os"
)

func GetCurrentBranchName(c *components.Context) (string, error) {
	if value := c.GetStringFlagValue("branch"); value != "" {
		return value, nil
	}

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

func GetCommitHash(branch string, remote bool) (string, error) {
	repository, err := getLocalRepo()
	if err != nil {
		return "", err
	}

	revision := branch
	if remote {
		revision = "refs/remotes/origin/" + revision
	}
	resolvedRev, err := repository.ResolveRevision(plumbing.Revision(revision))
	if resolvedRev.IsZero() {
		return "", fmt.Errorf("requested revision %s does not exist in the remote git", revision)
	}
	return resolvedRev.String(), nil
}

func getLocalRepo() (*git.Repository, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repository, err := git.PlainOpenWithOptions(wd, &git.PlainOpenOptions{
		EnableDotGitCommonDir: true,
		DetectDotGit:          true,
	})
	return repository, nil
}

func GetRootDir() (string, error) {
	repository, err := getLocalRepo()
	if err != nil {
		return "", err
	}

	worktree, err := repository.Worktree()
	if err != nil {
		return "", err
	}
	return worktree.Filesystem.Root(), nil
}
