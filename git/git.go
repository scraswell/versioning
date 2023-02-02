package git

import (
	"fmt"

	git "github.com/go-git/go-git/v5"
)

func Open(repositoryPath string) *git.Repository {
	repo, err := git.PlainOpen(repositoryPath)
	if err != nil {
		panic(fmt.Errorf("unable to read repository at path: %s\n%w", repositoryPath, err))
	}

	return repo
}
