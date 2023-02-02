package git

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spf13/viper"
)

const GoModule = "go.mod"

var v *viper.Viper
var c *Config

func init() {
	v = viperInit()
	c = getConfig(v)
}

func FindModules(startingDirectory string) []string {
	if !directoryExists(startingDirectory) {
		panic(fmt.Errorf("%s is not a directory or it doesn't exist", startingDirectory))
	}

	var modules []string

	filepath.Walk(startingDirectory, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.Name() == GoModule {
			moduleRelativePath := strings.ReplaceAll(filepath.Dir(path), startingDirectory, "")
			moduleRelativePath = strings.Trim(moduleRelativePath, "/")

			modules = append(modules, moduleRelativePath)
		}

		return nil
	})

	return modules
}

func FindBumps(repositoryPath string, branch string) map[string][]*Bump {
	re := GetVersionBumpRegexp()

	repo := Open(repositoryPath)
	hash := getHash(repo, branch)
	commitIter := getLog(repo, hash)
	moduleBumps := make(map[string][]*Bump)

	// TODO: this will need to be re-worked such that it only
	// traverses the commit chain until the most recent tag for
	// each module is found.  Obviously, it'll traverse the whole
	// thing if it doesn't find all the modules.

	commitIter.ForEach(func(commit *object.Commit) error {
		matches := re.FindAllStringSubmatch(commit.Message, -1)

		if len(matches) > 0 {
			recordBumps(moduleBumps, matches, &commit.Hash)
		}

		return nil
	})
	commitIter.Close()

	for _, v := range moduleBumps {
		Reverse(v)
	}

	resolveModuleVersions(moduleBumps)

	return moduleBumps
}

func resolveModuleVersions(moduleBumps map[string][]*Bump) {
	for k := range moduleBumps {
		previousVersion := getDefaultVersion()

		for _, bump := range moduleBumps[k] {
			if !previousVersion.Equals(getDefaultVersion()) {
				bump.Version = previousVersion
			}

			switch bump.Type {
			case BumpTypeMajor:
				bump.Version = *IncrementMajor(&bump.Version)
			case BumpTypeMinor:
				bump.Version = *IncrementMinor(&bump.Version)
			default:
				bump.Version = *IncrementPatch(&bump.Version)
			}

			previousVersion = bump.Version
		}
	}
}

func recordBumps(moduleBumps map[string][]*Bump, matches [][]string, hash *plumbing.Hash) {
	for _, match := range matches {
		moduleName := match[ModuleSubmatchIndex]

		if moduleBumps[moduleName] == nil {
			moduleBumps[moduleName] = []*Bump{}
		}

		moduleBumps[moduleName] = append(moduleBumps[moduleName], &Bump{
			Type:    getBumpType(match[BumpTypeSubmatchIndex]),
			Commit:  hash.String(),
			Version: getDefaultVersion(),
		})
	}
}

func getLog(repo *git.Repository, hash *plumbing.Hash) object.CommitIter {
	commitIter, err := repo.Log(&git.LogOptions{
		From:  *hash,
		All:   true,
		Order: git.LogOrderCommitterTime,
	})

	if err != nil {
		panic(fmt.Errorf("unable to read commits %w", err))
	}

	return commitIter
}

func getHash(repo *git.Repository, branch string) *plumbing.Hash {
	hash, err := repo.ResolveRevision(plumbing.Revision(branch))
	if err != nil {
		panic(fmt.Errorf("unable to resolve revision %w", err))
	}

	return hash
}
