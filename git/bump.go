package git

import (
	"fmt"

	"github.com/blang/semver/v4"
)

const BumpTypeMajor = "major"
const BumpTypeMinor = "minor"
const BumpTypePatch = "patch"
const DefaultVersion = "0.0.0"

type Bump struct {
	Type    string
	Commit  string
	Version semver.Version
}

func GetTag(moduleRelativePath string, bump *Bump) string {
	return fmt.Sprintf("%s/v%s", moduleRelativePath, bump.Version.String())
}

func getDefaultVersion() semver.Version {
	version, err := semver.Make(DefaultVersion)
	if err != nil {
		panic(fmt.Errorf("unable to parse semver for default version %w", err))
	}

	return version
}

func MakeVersion(version string) semver.Version {
	semversion, err := semver.Make(version)
	if err != nil {
		panic(fmt.Errorf("unable to parse semver for default version %w", err))
	}

	return semversion
}

func IncrementMajor(version *semver.Version) *semver.Version {
	err := version.IncrementMajor()
	if err != nil {
		panic(fmt.Errorf("unable to increment major component %w", err))
	}

	return version
}

func IncrementMinor(version *semver.Version) *semver.Version {
	err := version.IncrementMinor()
	if err != nil {
		panic(fmt.Errorf("unable to increment minor component %w", err))
	}

	return version
}

func IncrementPatch(version *semver.Version) *semver.Version {
	err := version.IncrementPatch()
	if err != nil {
		panic(fmt.Errorf("unable to increment major component %w", err))
	}

	return version
}

func getBumpType(parsedType string) string {
	if Contains(c.BumpStrings.Major, parsedType) {
		return BumpTypeMajor
	} else if Contains(c.BumpStrings.Minor, parsedType) {
		return BumpTypeMinor
	} else {
		return BumpTypePatch
	}
}
