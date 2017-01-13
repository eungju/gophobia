package resp

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"fmt"
)

const (
	simpleStringType = '+'
	errorType  = '-'
	intType    = ':'
	bulkStringType = '$'
	arrayType  = '*'
)

var lineSeparator = []byte("\r\n")
var lineSeparatorBytes = len(lineSeparator)

type RequestReader struct {
	source *bufio.Reader
}

func NewRequestReader(source io.Reader) *RequestReader {
	return &RequestReader{source: bufio.NewReader(source)}
}

func (reader *RequestReader) ReadLine() (string, error) {
	line, isPrefix, err := reader.source.ReadLine()
	if err != nil {
		return "", err
	}
	if isPrefix {
		return "", bufio.ErrBufferFull
	}
	return string(line), nil
}

func (reader *RequestReader) Read() ([]string, error) {
	line, err := reader.ReadLine()
	if err != nil {
		return nil, err
	}
	if line[0] == arrayType {
		n, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		array := make([]string, n)
		for i := 0; i < n; i++ {
			element, err := reader.ReadBulkString()
			if err != nil {
				return nil, err
			}
			array[i] = element
		}
		return array, nil
	}
	inlineCommandSeparator := regexp.MustCompile("\\s+")
	return inlineCommandSeparator.Split(string(line), -1), nil
}

func (reader *RequestReader) ReadBulkString() (string, error) {
	line, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	respType := line[0]
	if line[0] != bulkStringType {
		return "", fmt.Errorf("Expected bulk string type but %s", respType)
	}
	bytes, err := strconv.Atoi(line[1:])
	if err != nil {
		return "", err
	}
	buf := make([]byte, bytes)
	readBytes, err := reader.source.Read(buf)
	if bytes != readBytes {
		return "", fmt.Errorf("Expected %d bytes but %d", bytes, readBytes)
	}
	discardBuf := make([]byte, 2)
	discardBytes, err := reader.source.Read(discardBuf)
	if lineSeparatorBytes != discardBytes {
		return "", fmt.Errorf("Expected CRLF but %s", discardBuf)
	}
	return string(buf), nil
}
