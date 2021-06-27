package handlers

import (
	"reflect"
	"testing"
)

func TestGetPhoneNumbers(t *testing.T) {
	m := make(map[string]string)
	m["phones"] = ""

	result := getPhoneNumbers(m)
	if len(result) != 1 {
		t.Errorf("Result list lenght is %d, not %d", len(result), 1)
	}

	m = make(map[string]string)
	m["test"] = ""

	result = getPhoneNumbers(m)
	if len(result) != 0 {
		t.Errorf("Result list lenght is not null: %d", len(result))
	}

	m = make(map[string]string)
	m["phones"] = "79191111111,79192222222"
	tl := []string{"79191111111", "79192222222"}
	result = getPhoneNumbers(m)
	if len(result) != 2 && reflect.DeepEqual(result, tl) {
		t.Errorf("Result slice is not equal, result slice %v, test slice %v", result, tl)
	}
}
