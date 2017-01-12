package resp_test

import (
	. "."
	"reflect"
	"strings"
	"testing"
)

func TestParsingInlineCommand(t *testing.T) {
	dut := NewReader(strings.NewReader("mget a b c\r\n"))
	r, err := dut.Read()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(r, []string{"mget", "a", "b", "c"}) {
		t.Error("Unexpected")
	}
}
