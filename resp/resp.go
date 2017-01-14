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
	errorType = '-'
	integerType = ':'
	bulkStringType = '$'
	arrayType = '*'
)

var lineSeparator = []byte("\r\n")
var lineSeparatorBytes = len(lineSeparator)
var inlineCommandSeparator = regexp.MustCompile("\\s+")

type RespValue interface {
	AsSimpleString() (*RespSimpleString, error)
	AsError() (*RespError, error)
	AsInteger() (*RespInteger, error)
	AsBulkString() (*RespBulkString, error)
	AsArray() (*RespArray, error)
}

type RespSimpleString struct {
	Value []byte
}

func (resp *RespSimpleString) AsSimpleString() (*RespSimpleString, error) {
	return resp, nil
}

func (resp *RespSimpleString) AsError() (*RespError, error) {
	return nil, fmt.Errorf("Expected RespError but RespSimpleString")
}

func (resp *RespSimpleString) AsInteger() (*RespInteger, error) {
	return nil, fmt.Errorf("Expected RespInteger but RespSimpleString")
}

func (resp *RespSimpleString) AsBulkString() (*RespBulkString, error) {
	return nil, fmt.Errorf("Expected RespBulkString but RespSimpleString")
}

func (resp *RespSimpleString) AsArray() (*RespArray, error) {
	return nil, fmt.Errorf("Expected RespArray but RespSimpleString")
}

type RespError struct {
	Value []byte
}

func (resp *RespError) AsSimpleString() (*RespSimpleString, error) {
	return nil, fmt.Errorf("Expected RespSimpleString but RespError")
}

func (resp *RespError) AsError() (*RespError, error) {
	return resp, nil
}

func (resp *RespError) AsInteger() (*RespInteger, error) {
	return nil, fmt.Errorf("Expected RespInteger but RespError")
}

func (resp *RespError) AsBulkString() (*RespBulkString, error) {
	return nil, fmt.Errorf("Expected RespBulkString but RespError")
}

func (resp *RespError) AsArray() (*RespArray, error) {
	return nil, fmt.Errorf("Expected RespArray but RespError")
}

type RespInteger struct {
	Value int64
}

func (resp *RespInteger) AsSimpleString() (*RespSimpleString, error) {
	return nil, fmt.Errorf("Expected RespSimpleString but RespInteger")
}

func (resp *RespInteger) AsError() (*RespError, error) {
	return nil, fmt.Errorf("Expected RespInteger but RespInteger")
}

func (resp *RespInteger) AsInteger() (*RespInteger, error) {
	return resp, nil
}

func (resp *RespInteger) AsBulkString() (*RespBulkString, error) {
	return nil, fmt.Errorf("Expected RespBulkString but RespInteger")
}

func (resp *RespInteger) AsArray() (*RespArray, error) {
	return nil, fmt.Errorf("Expected RespArray but RespInteger")
}

type RespBulkString struct {
	Value []byte
}

func (resp *RespBulkString) AsSimpleString() (*RespSimpleString, error) {
	return nil, fmt.Errorf("Expected RespSimpleString but RespBulkString")
}

func (resp *RespBulkString) AsError() (*RespError, error) {
	return nil, fmt.Errorf("Expected RespError but RespBulkString")
}

func (resp *RespBulkString) AsInteger() (*RespInteger, error) {
	return nil, fmt.Errorf("Expected RespInteger but RespBulkString")
}

func (resp *RespBulkString) AsBulkString() (*RespBulkString, error) {
	return resp, nil
}

func (resp *RespBulkString) AsArray() (*RespArray, error) {
	return nil, fmt.Errorf("Expected RespArray but RespBulkString")
}

type RespArray struct {
	Value []RespValue
}

func (resp *RespArray) AsSimpleString() (*RespSimpleString, error) {
	return nil, fmt.Errorf("Expected RespSimpleString but RespArray")
}

func (resp *RespArray) AsError() (*RespError, error) {
	return nil, fmt.Errorf("Expected RespError but RespArray")
}

func (resp *RespArray) AsInteger() (*RespInteger, error) {
	return nil, fmt.Errorf("Expected RespInteger but RespArray")
}

func (resp *RespArray) AsBulkString() (*RespBulkString, error) {
	return nil, fmt.Errorf("Expected RespBulkString but RespArray")
}

func (resp *RespArray) AsArray() (*RespArray, error) {
	return resp, nil
}

type Reader struct {
	source *bufio.Reader
}

func NewReader(source io.Reader) *Reader {
	return &Reader{source: bufio.NewReader(source)}
}

func (reader *Reader) Read() (RespValue, error) {
	line, err := reader.readLine()
	if err != nil {
		return nil, err
	}
	respType := line[0]
	if simpleStringType == respType {
		return &RespSimpleString{line[1:]}, nil
	} else if errorType == respType {
		return &RespError{line[1:]}, nil
	} else if integerType == respType {
		value, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return &RespInteger{int64(value)}, nil
	} else if bulkStringType == respType {
		bytes, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return reader.readBulkString(bytes)
	} else if arrayType == respType {
		length, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		return reader.readArray(length)
	}
	return nil, fmt.Errorf("Unknown type %s", respType)
}

func (reader *Reader) readBulkString(bytes int) (*RespBulkString, error) {
	value := make([]byte, bytes)
	readBytes, err := reader.source.Read(value)
	if err != nil {
		return nil, err
	}
	if bytes != readBytes {
		return nil, fmt.Errorf("Expected %d bytes but %d", bytes, readBytes)
	}
	discardBuf := make([]byte, 2)
	discardBytes, err := reader.source.Read(discardBuf)
	if lineSeparatorBytes != discardBytes {
		return nil, fmt.Errorf("Expected CRLF but %s", discardBuf)
	}
	return &RespBulkString{value}, nil
}

func (reader *Reader) readArray(length int) (*RespArray, error) {
	value := make([]RespValue, length)
	for i := 0; i < length; i++ {
		element, err := reader.Read()
		if err != nil {
			return nil, err
		}
		value[i] = element
	}
	return &RespArray{value}, nil
}

func (reader *Reader) readLine() ([]byte, error) {
	line, isPrefix, err := reader.source.ReadLine()
	if err != nil {
		return nil, err
	}
	if isPrefix {
		return nil, bufio.ErrBufferFull
	}
	return line, nil
}

type CommandReader struct {
	Reader
}

func NewCommandReader(source io.Reader) *CommandReader {
	return &CommandReader{Reader{bufio.NewReader(source)}}
}

func (reader *CommandReader) Read() ([]string, error) {
	line, err := reader.readLine()
	if err != nil {
		return nil, err
	}
	if line[0] == arrayType {
		length, err := strconv.Atoi(string(line[1:]))
		if err != nil {
			return nil, err
		}
		array, err := reader.readArray(length)
		command := make([]string, length)
		for i := 0; i < length; i++ {
			element, err := array.Value[i].AsBulkString()
			if err != nil {
				return nil, err
			}
			command[i] = string(element.Value)
		}
		return command, nil
	}
	return inlineCommandSeparator.Split(string(line), -1), nil
}
