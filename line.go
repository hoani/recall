package main

import (
	"encoding/json"
	"fmt"
	"time"
)

type f struct{}

type logData struct {
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
			tvalue, err := time.Parse(time.StampMilli, value)
			if err != nil {
				return err
			}
			d.Time = &tvalue
			delete(m, "time")
		}
	}

	if mvalue, ok := m["level"]; ok {
		if value, ok := mvalue.(string); ok {
			d.Level = value
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

func NewLineFormatter() *f {
	return &f{}
}

func (f *f) Format(input []byte) (string, error) {

	data := logData{}
	if err := json.Unmarshal(input, &data); err != nil {
		fmt.Printf("failed to unmarshal, err: %v\n", err)
		return "", err
	}

	res := fmt.Sprintf("[%v] %v: %v, %v", data.Time.Format(time.StampMilli), data.Level, data.Message, data.Meta)
	return res, nil
}
