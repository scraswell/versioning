package main

import (
	"flag"
	"fmt"

	git "github.com/scraswell/versioning/git"
)

var (
	dir        string
	branch     string
	dryrun     bool
	findgomods bool
)

func init() {
	flag.StringVar(&dir, "dir", "/Users/sean/code/sre/cdp", "The path to the folder that contains a .git repository.")
	flag.StringVar(&branch, "branch", "main", "The branch or reference for which to parse the logs for version bumps.")
	flag.BoolVar(&findgomods, "findgomods", false, "A value indicating whether we should find golang modules in the specified directory.")
	flag.BoolVar(&dryrun, "dryrun", true, "A value indicating whether to actually add the tags to the repo or just print which tags would be added.")
}

func main() {
	flag.Parse()

	if findgomods {
		modules := git.FindModules(dir)

		for _, v := range modules {
			fmt.Println(v)
		}
	} else if dryrun {
		bumps := git.FindBumps(dir, branch)

		for mod := range bumps {
			for _, bump := range bumps[mod] {
				fmt.Printf("%s => %s\n", git.GetTag(mod, bump), bump.Commit)
			}
		}
	} else {
		panic("not yet implemented.")
	}
}
