package main

import (
	"strings"
	"time"
)

type EntryInterpretter interface {
	Time(map[string]interface{}) (*time.Time, error)
	Level(map[string]interface{}) (string, error)
	Message(map[string]interface{}) (string, error)
	Error(map[string]interface{}) (string, error)
}

type defaultInterpreter struct {
	locale *time.Location
}

func (i *defaultInterpreter) Time(m map[string]interface{}) (*time.Time, error) {
	var result *time.Time
	if mvalue, ok := m["time"]; ok {
		switch value := mvalue.(type) {
		case float64:
			nanosecs := int64(value*1e9) % int64(1e9)
			tvalue := time.Unix(int64(value), nanosecs).In(i.locale)
			result = &tvalue
			delete(m, "time")
		case string:
			tvalue, err := time.ParseInLocation(time.RFC3339Nano, value, i.locale)
			if err != nil {
				return nil, err
			}
			result = &tvalue
			delete(m, "time")
		}
	}
	return result, nil
}

func (i *defaultInterpreter) Level(m map[string]interface{}) (string, error) {
	var result string
	if mvalue, ok := m["level"]; ok {
		if value, ok := mvalue.(string); ok {
			result = strings.ToUpper(value)
			if len(result) > 5 {
				if strings.HasPrefix(result, "WARN") {
					result = "WARN"
				} else {
					result = result[:5]
				}
			}
			delete(m, "level")
		}
	}
	return result, nil
}
func (i *defaultInterpreter) Message(m map[string]interface{}) (string, error) {
	var result string
	if mvalue, ok := m["message"]; ok {
		if value, ok := mvalue.(string); ok {
			result = value
			delete(m, "message")
		}
	}
	return result, nil
}
func (i *defaultInterpreter) Error(m map[string]interface{}) (string, error) {
	var result string
	return result, nil
}
