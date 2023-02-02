package git_test

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	r "math/rand"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

var worktreeFiles = [...]string{
	"file1",
	"file2",
	"file3",
	"file4",
	"file5",
	"file6",
	"file7",
	"file8",
	"file9",
}

func generateGitRepository(commitMessages ...string) string {
	repoPath, err := ioutil.TempDir("/tmp", "test-git-repo-*")
	if err != nil {
		panic(fmt.Errorf("failed to create temp dir: %w", err))
	}

	fmt.Printf("Repo path: %s\n", repoPath)
	repo, err := gogit.PlainInit(repoPath, false)
	if err != nil {
		panic(fmt.Errorf("failed to create git repo: %w", err))
	}

	worktree, err := repo.Worktree()
	if err != nil {
		panic(fmt.Errorf("failed to get worktree: %w", err))
	}

	for _, message := range commitMessages {
		createCommit(worktree, message)
	}

	return repoPath
}

func createCommit(wt *gogit.Worktree, commitMessage string) plumbing.Hash {
	file := worktreeFiles[r.Intn(len(worktreeFiles))]

	writeToFileInWorktree(wt, file)

	_, err := wt.Add(file)
	if err != nil {
		panic(fmt.Errorf("failed to stage file: %w", err))
	}

	commit, err := wt.Commit(commitMessage, &gogit.CommitOptions{})
	if err != nil {
		panic(fmt.Errorf("commit failed: %w", err))
	}

	return commit
}

func writeToFileInWorktree(wt *gogit.Worktree, fileName string) {
	filePath := addFileToWorktree(wt, fileName)

	file := openFile(filePath)
	file.WriteString(generateRandomStringURLSafe(16))
	file.Close()
}

func addFileToWorktree(wt *gogit.Worktree, fileName string) string {
	filePath := wt.Filesystem.Join(wt.Filesystem.Root(), fileName)

	if !fileExists(filePath) {
		file := createFile(filePath)
		file.Close()
	}

	return filePath
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	if err == nil {
		return true
	} else if err != nil && errors.Is(err, os.ErrNotExist) {
		return false
	} else {
		panic("stat call failed")
	}
}

func openFile(filePath string) *os.File {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("unable to open file %s; %w", filePath, err))
	}

	return file
}

func createFile(filePath string) *os.File {
	file, err := os.Create(filePath)
	if err != nil {
		panic(fmt.Errorf("error creating file (%s): %w", filePath, err))
	}

	return file
}

func assertAvailablePRNG() {
	buf := make([]byte, 1)

	_, err := io.ReadFull(rand.Reader, buf)
	if err != nil {
		panic(fmt.Sprintf("crypto/rand is unavailable: Read() failed with %#v", err))
	}
}

func generateRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)

	if err != nil {
		panic(fmt.Errorf("failed to generate random bytes: %w", err))
	}

	return b
}

func generateRandomStringURLSafe(n int) string {
	return base64.URLEncoding.EncodeToString(generateRandomBytes(n))
}
