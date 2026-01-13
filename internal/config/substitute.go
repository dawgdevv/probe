package config

import (
	"fmt"
	"regexp"
)

var re = regexp.MustCompile(`\{\{(\w+)\}\}`)

func SubstituteString(input string, env map[string]string) (string, error) {
	result := re.ReplaceAllStringFunc(input, func(match string) string {
		key := re.FindStringSubmatch(match)[1]

		if val, ok := env[key]; ok {
			return val
		}

		return match
	})

	if re.MatchString(result) {
		return "", fmt.Errorf("unresolved variable in string: %s", input)
	}

	return result, nil

}
