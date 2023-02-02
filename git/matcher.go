package git

import (
	"fmt"
	"regexp"
	"strings"
)

const BumpPattern = `\b(((\w+\/?)+)__SEPARATOR__(__MAJOR__|__MINOR__|__PATCH__))\b`
const ModuleSubmatchIndex = 2
const BumpTypeSubmatchIndex = 4

func GetVersionBumpRegexp() *regexp.Regexp {
	pattern := BumpPattern
	pattern = strings.ReplaceAll(pattern, `__SEPARATOR__`, fmt.Sprintf(`\%s`, c.VersionSeparator))
	pattern = strings.ReplaceAll(pattern, `__MAJOR__`, strings.Join(c.BumpStrings.Major, `|`))
	pattern = strings.ReplaceAll(pattern, `__MINOR__`, strings.Join(c.BumpStrings.Minor, `|`))
	pattern = strings.ReplaceAll(pattern, `__PATCH__`, strings.Join(c.BumpStrings.Patch, `|`))

	return regexp.MustCompile(pattern)
}
