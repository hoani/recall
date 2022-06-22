package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type f struct {
	index      int64
	locale     *time.Location
	timeLength int
}

type logData struct {
	locale  *time.Location
	Time    *time.Time
	Level   string
	Message string
	Meta    map[string]interface{}
}

func (d *logData) UnmarshalJSON(data []byte) error {
	m := make(map[string]interface{})
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}

	if mvalue, ok := m["time"]; ok {
		switch value := mvalue.(type) {
		case float64:
			tvalue := time.Unix(int64(value), 0)
			d.Time = &tvalue
			delete(m, "time")
		case string:
			tvalue, err := time.ParseInLocation(time.RFC3339Nano, value, d.locale)
			if err != nil {
				return err
			}
			d.Time = &tvalue
			delete(m, "time")
		}
	}

	if mvalue, ok := m["level"]; ok {
		if value, ok := mvalue.(string); ok {
			d.Level = strings.ToUpper(value)
			if len(d.Level) > 5 {
				if strings.HasPrefix(d.Level, "WARN") {
					d.Level = "WARN"
				} else {
					d.Level = d.Level[:5]
				}
			}
			delete(m, "level")
		}
	}

	if mvalue, ok := m["message"]; ok {
		if value, ok := mvalue.(string); ok {
			d.Message = value
			delete(m, "message")
		}
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
	output += fmt.Sprintf(" %-5v %v", data.Level, data.Message)
	if len(data.Meta) != 0 {
		output += fmt.Sprintf(", %v", data.Meta)
	}
	return output, nil
}
