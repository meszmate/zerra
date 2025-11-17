package crypt

import "regexp"

func ValidateHexColor(s string) bool {
	return regexp.MustCompile(`^#([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$`).MatchString(s)
}
