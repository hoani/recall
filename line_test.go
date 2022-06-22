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
				"time":    time.Unix(0, 0).Add(time.Second + time.Millisecond + time.Microsecond),
				"level":   "error",
				"message": "something went wrong",
			},
			expected: "[Jan  1 12:00:01.001] ERROR something went wrong",
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
