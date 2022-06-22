package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertLine_Success(t *testing.T) {
	testCases := []struct {
		name     string
		data     map[string]interface{}
		expected string
	}{
		{
			name: "simple with error",
			data: map[string]interface{}{
				"time":    time.Unix(0, 0).Add(time.Second + time.Millisecond + time.Microsecond).In(time.UTC),
				"level":   "error",
				"message": "something went wrong",
			},
			expected: "[Jan  1 00:00:01.001] ERROR something went wrong",
		},
		{
			name: "info in 2016",
			data: map[string]interface{}{
				"time":    time.Date(2016, 1, 10, 7, 7, 7, 777e6, time.UTC),
				"level":   "info",
				"message": "everybody knows me now",
			},
			expected: "[Jan 10 07:07:07.777] INFO  everybody knows me now",
		},
		{
			name: "no time field",
			data: map[string]interface{}{
				"level":   "warning",
				"message": "you're almost there, but don't panic",
			},
			expected: "[                  1] WARN  you're almost there, but don't panic",
		},
		{
			name: "missing message",
			data: map[string]interface{}{
				"time":  time.Unix(0, 0).In(time.UTC),
				"level": "panic",
			},
			expected: "[Jan  1 00:00:00.000] PANIC",
		},
		{
			name: "unknown level",
			data: map[string]interface{}{
				"time":  time.Unix(0, 0).In(time.UTC),
				"level": "something",
			},
			expected: "[Jan  1 00:00:00.000] SOMET",
		},
		{
			name: "missing level",
			data: map[string]interface{}{
				"time":    time.Unix(0, 0).In(time.UTC),
				"message": "hello",
			},
			expected: "[Jan  1 00:00:00.000]       hello",
		},
		{
			name: "numerical time",
			data: map[string]interface{}{
				"time": float64(time.Unix(0, 0).In(time.UTC).Unix()) + 0.25,
			},
			expected: "[Jan  1 00:00:00.250]      ",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := json.Marshal(tc.data)
			assert.Nil(t, err)

			c := NewLineFormatter(time.UTC)
			s, err := c.Format(b)
			assert.Nil(t, err)
			assert.Equal(t, tc.expected, s)
		})
	}
}
