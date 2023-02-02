package git_test

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	git "github.com/scraswell/versioning/git"
)

func init() {
	rand.Seed(time.Now().Unix())
	assertAvailablePRNG()
}

func TestMatcherMatchesBumpStrings(t *testing.T) {
	re := git.GetVersionBumpRegexp()

	expectedMatches := []string{
		"module1+patch",
		"path/to/module2+minor",
	}

	matches := re.FindAllString("fixed bug module1+patch path/to/module2+minor", -1)

	if len(matches) != 2 {
		t.Fail()
	}

	for _, match := range matches {
		if !git.Contains(expectedMatches, match) {
			t.Fail()
		}
	}
}

func TestVersionResolutionFromCommitMessages(t *testing.T) {
	var commitMessages = []string{
		"Initial commit",
		"fixed bug module1+patch path/to/module2+minor", // module1@0.0.1 & path/to/module2@0.1.0
		"refactoring path/to/module2+minor",             // path/to/module2@0.2.0
		"refactoring module1+minor",                     // module1@0.1.0
		"release prep",
		"release module1+major path/to/module2+minor", // module1@1.0.0 & path/to/module2@0.3.0
		"bugfix path/to/module2+patch",                // path/to/module2@0.3.1
	}

	gitRepo := generateGitRepository(commitMessages...)
	bumps := git.FindBumps(gitRepo, "master")

	performChecks(bumps, t)

	err := os.RemoveAll(gitRepo)
	if err != nil {
		panic(fmt.Errorf("failed to remove temporary git repo @ %s : %w", gitRepo, err))
	}
}

func TestVersionResolutionFromCommitMessagesWithShorthandBumps(t *testing.T) {
	var commitMessages = []string{
		"Initial commit",
		"fixed bug module1+p path/to/module2+n", // module1@0.0.1 & path/to/module2@0.1.0
		"refactoring path/to/module2+n",         // path/to/module2@0.2.0
		"refactoring module1+n",                 // module1@0.1.0
		"release prep",
		"release module1+j path/to/module2+n", // module1@1.0.0 & path/to/module2@0.3.0
		"bugfix path/to/module2+p",            // path/to/module2@0.3.1
	}

	gitRepo := generateGitRepository(commitMessages...)
	bumps := git.FindBumps(gitRepo, "master")

	performChecks(bumps, t)

	err := os.RemoveAll(gitRepo)
	if err != nil {
		panic(fmt.Errorf("failed to remove temporary git repo @ %s : %w", gitRepo, err))
	}
}

func performChecks(bumps map[string][]*git.Bump, t *testing.T) {

	if len(bumps) != 2 {
		t.Fail()
	}

	if bumps["module1"] == nil {
		t.Fail()
	}

	if bumps["path/to/module2"] == nil {
		t.Fail()
	}

	if !bumps["path/to/module2"][0].Version.Equals(git.MakeVersion("0.1.0")) {
		t.Fail()
	}

	if !bumps["path/to/module2"][1].Version.Equals(git.MakeVersion("0.2.0")) {
		t.Fail()
	}

	if !bumps["path/to/module2"][2].Version.Equals(git.MakeVersion("0.3.0")) {
		t.Fail()
	}

	if !bumps["path/to/module2"][3].Version.Equals(git.MakeVersion("0.3.1")) {
		t.Fail()
	}

	if !bumps["module1"][0].Version.Equals(git.MakeVersion("0.0.1")) {
		t.Fail()
	}

	if !bumps["module1"][1].Version.Equals(git.MakeVersion("0.1.0")) {
		t.Fail()
	}

	if !bumps["module1"][2].Version.Equals(git.MakeVersion("1.0.0")) {
		t.Fail()
	}
}
