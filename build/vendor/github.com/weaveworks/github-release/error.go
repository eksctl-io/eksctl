package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

/* usually when something goes wrong, github sends something like this back */
type Message struct {
	Message string        `json:"message"`
	Errors  []GithubError `json:"errors"`
}

type GithubError struct {
	Resource string `json:"resource"`
	Code     string `json:"code"`
	Field    string `json:"field"`
}

/* transforms a stream into a Message, if it's valid json */
func ToMessage(r io.Reader) (*Message, error) {
	var msg Message
	if err := json.NewDecoder(r).Decode(&msg); err != nil {
		return nil, err
	}

	return &msg, nil
}

func (m *Message) String() string {
	str := fmt.Sprintf("msg: %v, errors: ", m.Message)

	errstr := make([]string, len(m.Errors))
	for idx, err := range m.Errors {
		errstr[idx] = fmt.Sprintf("[field: %v, code: %v]",
			err.Field, err.Code)
	}

	return str + strings.Join(errstr, ", ")
}
