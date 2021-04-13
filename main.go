package rfc5424

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

type Message struct {
	Priority  int
	Version   int
	Timestamp time.Time
	Hostname  string
	AppName   string
	ProcID    string
	Message   string
}

// readLength reads the octet-length of the next message from a reader
// and return the octet-length, bytes read and a possible error
func readLength(source io.Reader) (length int, err error) {
	const readSize = 1
	var buf []byte
	p := make([]byte, readSize)

	n, err := source.Read(p)
	if err != nil {
		return 0, err
	}
	if n != readSize {
		return 0, fmt.Errorf("expected to read %v but read %v instead", readSize, n)
	}
	for !bytes.Equal(p, []byte(" ")) {
		buf = append(buf, p[0])
		n, err = source.Read(p)
		if err != nil {
			return 0, err
		}
		if n != readSize {
			return 0, fmt.Errorf("expected to read %v but read %v instead", readSize, n)
		}
	}

	length, err = strconv.Atoi(string(buf))
	if err != nil {
		return 0, err
	}

	return length, nil
}

// Parse parses a bytes buffer for a RFC5424 syslog messages and returns it
func Parse(source io.Reader) (message Message, err error) {
	var thisMessage Message

	length, err := readLength(source)
	if err != nil {
		return thisMessage, err
	}

	messageReader := io.LimitReader(source, int64(length))
	messageBuffer, err := io.ReadAll(messageReader)
	if err != nil {
		return thisMessage, err
	}

	if len(messageBuffer) != length {
		return thisMessage, fmt.Errorf("octet-length [%v] mismatch with payload [%v]", length, len(messageBuffer))
	}

	thisMessage.Message = strings.TrimRight(string(messageBuffer), "\r\n")

	return thisMessage, nil
}

// ParseMultiple parses a bytes buffer for one or more RFC5424 syslog messages
// and returns an array of them
func ParseMultiple(source io.Reader) (messages []Message, err error) {
	message, err := Parse(source)
	messages = append(messages, message)
	if err != nil {
		fmt.Println("error", err)
	}

	return messages, nil
}
