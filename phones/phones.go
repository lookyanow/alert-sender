package phones

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

type ReceiverGroup struct {
	Name   string   `json:"name"`
	Phones []string `json:"phones"`
}

func NewReceiverGroup(filename string) (*[]ReceiverGroup, error) {
	rec, err := getPhonesFromJsonFile(filename)
	if err != nil {
		s := fmt.Sprintf("Phones file: '%s', parsed with error: %s", filename, err)
		return nil, errors.New(s)
	}
	return &rec, nil
}

func getPhonesFromJsonFile(filename string) ([]ReceiverGroup, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var data []ReceiverGroup
	err1 := json.Unmarshal(f, &data)
	if err1 != nil {
		return nil, err1
	}

	return data, nil

}
