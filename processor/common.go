// SPDX-License-Identifier: MIT OR Unlicense

package processor

import (
	"regexp"
	"strings"
)

var alphaNumericRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
var multipleSpacesRegex = regexp.MustCompile("\\s+")

// Very specific cleaner which is designed to clean to the format lc uses to match things
// so be very careful if you choose to use this
func LcCleanText(content string) string {
	content = strings.ToLower(content)

	content = alphaNumericRegex.ReplaceAllString(content, " ")
	content = multipleSpacesRegex.ReplaceAllString(content, " ")
	content = strings.TrimSpace(content)

	return content
}
