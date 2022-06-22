package main

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_ConvertLine(t *testing.T) {
	data := map[string]interface{}{
		"time":    time.Unix(0, 0).Add(time.Second + time.Millisecond + time.Microsecond),
		"level":   "error",
		"message": "something went wrong",
	}
	b, err := json.Marshal(data)
	assert.Nil(t, err)

	c := NewLineFormatter(time.UTC)
	s, err := c.Format(b)
	assert.Nil(t, err)
	assert.Equal(t, "[Jan  1 12:00:01.001] error: something went wrong", s)
}
