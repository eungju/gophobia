package resp_test

import (
	. "."
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParsingInlineCommand(t *testing.T) {
	dut := NewReader(strings.NewReader("mget a b c\r\n"))
	r, err := dut.Read()
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, []string{"mget", "a", "b", "c"}, r)
}
