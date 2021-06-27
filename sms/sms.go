package sms

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type Sender struct {
	SmsURL      string
	SmsUser     string
	SmsPassword string
}

type sms struct {
	Message    string   `json:"message"`
	Recipients []string `json:"recipients"`
}

func NewSender(smsUrl, smsUser, smsPassword string) (*Sender, error) {
	return &Sender{
		SmsURL:      smsUrl,
		SmsUser:     smsUser,
		SmsPassword: smsPassword,
	}, nil
}

// SendSMS send sms message to api
func (s Sender) SendSMS(message string, status string, recipients []string) error {
	if len(recipients) == 0 {
		return errors.New("Empty recipients list")
	}

	str := fmt.Sprintf("%s: %s", strings.ToUpper(status), message) // "{STATUS}: message text" - message format
	m := sms{
		Message:    str,
		Recipients: recipients,
	}
	requestBody, err := json.Marshal(m)
	if err != nil {
		return err
	}

	var netClient = &http.Client{
		Timeout: time.Second * 2, // todo: may be should get timeout value from config params
	}
	req, _ := http.NewRequest("POST", s.SmsURL, bytes.NewBuffer(requestBody))
	req.SetBasicAuth(s.SmsUser, s.SmsPassword)
	_, err1 := netClient.Do(req)
	if err1 != nil {
		return err1
	}
	return nil
}
