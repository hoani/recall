package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type f struct {
	index  int64
	locale *time.Location
}

type logData struct {
	locale  *time.Location
	Time    *time.Time
	Level   string
	Message string
	Meta    map[string]interface{}
}

func (d *logData) UnmarshalJSON(data []byte) error {
	var err error
	i := &defaultInterpreter{d.locale}
	m := make(map[string]interface{})
	if err = json.Unmarshal(data, &m); err != nil {
		return err
	}

	if d.Time, err = i.Time(m); err != nil {
		return err
	}
	if d.Level, err = i.Level(m); err != nil {
		return err
	}
	if d.Message, err = i.Message(m); err != nil {
		return err
	}

	d.Meta = m
	return nil
}

func NewLineFormatter(locale *time.Location) *f {
	return &f{
		index:  1,
		locale: locale,
	}
}

func (f *f) Format(input []byte) (string, error) {
	defer func() { f.index++ }()

	data := logData{
		locale: f.locale,
	}
	if err := json.Unmarshal(input, &data); err != nil {
		fmt.Printf("failed to unmarshal, err: %v\n", err)
		return "", err
	}

	var output string
	if data.Time == nil {
		indexStr := strconv.FormatInt(f.index, 10)
		output += "[" + strings.Repeat(" ", len(time.StampMilli)-len(indexStr)) + indexStr + "]"
	} else {
		output += "[" + data.Time.Format(time.StampMilli) + "]"
	}

	output += fmt.Sprintf(" %-5v", data.Level)
	if len(data.Message) != 0 {
		output += " " + data.Message
	}
	if len(data.Meta) != 0 {
		output += fmt.Sprintf(", %v", data.Meta)
	}
	return output, nil
}
