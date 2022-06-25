package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/i582/cfmt/cmd/cfmt"
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
	output += f.formatTime(data.Time)
	output += " " + f.formatLevel(data.Level)

	if len(data.Message) != 0 {
		output += " " + data.Message
	}
	if len(data.Meta) != 0 {
		output += fmt.Sprintf(", %v", data.Meta)
	}
	return output, nil
}

func (f *f) formatTime(t *time.Time) string {
	if t == nil {
		index := strconv.FormatInt(f.index, 10)
		padding := strings.Repeat(" ", len(time.StampMilli)-len(index))
		return cfmt.Sprintf("[%s%s]", padding, index)
	} else {
		return cfmt.Sprintf("{{[%s]}}::blue", t.Format(time.StampMilli))
	}
}

func (f *f) formatLevel(level string) string {
	color := "#888888"
	switch strings.ToUpper(level) {
	case "INFO":
		color = "#33ddff"
	case "WARN":
		color = "#ffaa00"
	case "DEBUG":
		color = "#bb99bb"
	case "ERROR":
		color = "#ee4232"
	case "PANIC":
		color = "#ff4288"
	}
	return cfmt.Sprintf("{{%-5v}}::"+color+"|bold", level)
}
