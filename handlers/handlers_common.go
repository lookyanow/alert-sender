package handlers

import (
	"strings"
)

func getPhoneNumbers(m map[string]string) []string {
	if v, ok := m["phones"]; ok {
		return strings.Split(v, ",")
	}
	return nil
}

func getGroupName(m map[string]string) string {
	if v, ok := m["group"]; ok {
		return v
	}
	return ""
}
