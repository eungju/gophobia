package resp_test

import (
	. "."
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParsingInlineCommand(t *testing.T) {
	dut := NewRequestReader(strings.NewReader("mget a b c\r\n"))
	r, err := dut.Read()
	assert.NoError(t, err)
	assert.Equal(t, []string{"mget", "a", "b", "c"}, r)
}

func TestParsingUnifiedCommand(t *testing.T) {
	dut := NewRequestReader(strings.NewReader("*4\r\n$4\r\nmget\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n"))
	r, err := dut.Read()
	assert.NoError(t, err)
	assert.Equal(t, []string{"mget", "a", "b", "c"}, r)
}
