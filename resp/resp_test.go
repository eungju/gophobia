package resp_test

import (
	. "."
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParsingSimpleString(t *testing.T) {
	dut := NewReader(strings.NewReader("+OK\r\n"))
	actual, _ := dut.Read()
	assert.Equal(t, &RespSimpleString{[]byte("OK")}, actual)
}

func TestParsingError(t *testing.T) {
	dut := NewReader(strings.NewReader("-Error message\r\n"))
	actual, _ := dut.Read()
	assert.Equal(t, &RespError{[]byte("Error message")}, actual)
}

func TestParsingInteger(t *testing.T) {
	dut := NewReader(strings.NewReader(":42\r\n"))
	actual, _ := dut.Read()
	assert.Equal(t, &RespInteger{42}, actual)
}

func TestParsingBulkString(t *testing.T) {
	dut := NewReader(strings.NewReader("$3\r\nget\r\n"))
	actual, _ := dut.Read()
	assert.Equal(t, &RespBulkString{[]byte("get")}, actual)
}

func TestParsingArray(t *testing.T) {
	dut := NewReader(strings.NewReader("*2\r\n$1\r\na\r\n$1\r\nb\r\n"))
	actual, _ := dut.Read()
	assert.Equal(t, &RespArray{[]RespValue{&RespBulkString{[]byte("a")}, &RespBulkString{[]byte("b")}}}, actual)
}

func TestParsingInlineCommand(t *testing.T) {
	dut := NewCommandReader(strings.NewReader("mget a b c\r\n"))
	r, _ := dut.Read()
	assert.Equal(t, []string{"mget", "a", "b", "c"}, r)
}

func TestParsingUnifiedCommand(t *testing.T) {
	dut := NewCommandReader(strings.NewReader("*4\r\n$4\r\nmget\r\n$1\r\na\r\n$1\r\nb\r\n$1\r\nc\r\n"))
	r, _ := dut.Read()
	assert.Equal(t, []string{"mget", "a", "b", "c"}, r)
}
