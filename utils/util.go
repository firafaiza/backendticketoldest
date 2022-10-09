package utils

import "strings"

func CheckEmptyStringRequest(list map[string]string) string {
	for x := range list {
		if list[x] == "" || len(strings.TrimSpace(list[x])) == 0 {
			return x
		}
	}
	return ""
}
