package resp

import (
	"bufio"
	"io"
	"regexp"
)

const (
	ErrorReply  = '-'
	StatusReply = '+'
	IntReply    = ':'
	StringReply = '$'
	ArrayReply  = '*'
)

type Reader struct {
	source *bufio.Reader
}

func NewReader(source io.Reader) *Reader {
	return &Reader{source: bufio.NewReader(source)}
}

func (self *Reader) Read() ([]string, error) {
	line, isPrefix, err := self.source.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, bufio.ErrBufferFull
	}
	inlineCommandSeparator := regexp.MustCompile("\\s+")
	return inlineCommandSeparator.Split(string(line), -1), nil
}
