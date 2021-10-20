package valid

import (
	"net/url"
	"regexp"
)

// TODO validate URL and type (regex)

func ValidateURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	return err == nil
}

var slugRegex = regexp.MustCompile(`^\w+(-\w+)*$`)

func ValidateSlug(s string) bool {
	return slugRegex.MatchString(s)
}
