package rfc5424

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var severities = []string{
	"emerg",
	"alert",
	"crit",
	"err",
	"warning",
	"notice",
	"info",
	"debug",
}

var facilities = []string{
	"kernel",
	"user",
	"mail",
	"daemon",
	"auth",
	"syslog",
	"lpr",
	"news",
	"uucp",
	"clock",
	"authpriv",
	"ftp",
	"ntp",
	"audit",
	"logalert",
	"cron",
	"local0",
	"local1",
	"local2",
	"local3",
	"local4",
	"local5",
	"local6",
	"local7",
}

type Message struct {
	Priority  int
	Facility  int
	Severity  int
	Version   int
	Timestamp string
	Hostname  string
	AppName   string
	ProcID    string
	MessageID string
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

// parsePriority parses the priority integer as per RFC5424 section 6.2.1
// and returns the separate facility and severity
func parsePriority(priority int) (facility int, severity int) {
	severity = priority % 8
	facility = (priority - severity) / 8
	return facility, severity
}

// FacilityString returns a string describing the facility of the log message
func (m Message) FacilityString() (facility string, err error) {
	if m.Facility < 0 || m.Facility >= len(facilities) {
		return "", fmt.Errorf("facility out of range [%v]", m.Facility)
	}
	return facilities[m.Facility], nil
}

// SeverityString returns a string describing the severity of the log message
func (m Message) SeverityString() (severity string, err error) {
	if m.Severity < 0 || m.Severity >= len(severities) {
		return "", fmt.Errorf("severity out of range [%v]", m.Severity)
	}
	return severities[m.Severity], nil
}

// Parse parses a bytes buffer for a RFC5424 syslog messages and returns it
func Parse(source io.Reader) (message Message, err error) {
	var thisMessage Message

	length, err := readLength(source)
	if err != nil {
		return thisMessage, err
	}

	payloadReader := io.LimitReader(source, int64(length))
	payloadBuffer, err := io.ReadAll(payloadReader)
	if err != nil {
		return thisMessage, err
	}

	if len(payloadBuffer) != length {
		return thisMessage, fmt.Errorf("octet-length [%v] mismatch with payload [%v]", length, len(payloadBuffer))
	}

	payloadString := strings.TrimRight(string(payloadBuffer), "\r\n")
	payloadSlice := strings.SplitN(payloadString, " ", 7)

	if len(payloadSlice) != 7 {
		return thisMessage, fmt.Errorf("missing fields in payload [%v]", payloadString)
	}

	priorityTempString := strings.TrimLeft(payloadSlice[0], "<")
	priorityTempString = strings.TrimRight(priorityTempString, "0123456789")
	priorityTempString = strings.TrimRight(priorityTempString, ">")
	thisMessage.Priority, err = strconv.Atoi(priorityTempString)
	if err != nil {
		return thisMessage, fmt.Errorf("couldn't parse priority field [%v] in payload [%v]", payloadSlice[0], payloadString)
	}

	thisMessage.Facility, thisMessage.Severity = parsePriority(thisMessage.Priority)

	versionTempString := strings.TrimLeft(payloadSlice[0], "<")
	versionTempString = strings.TrimLeft(versionTempString, "0123456789")
	versionTempString = strings.TrimLeft(versionTempString, ">")
	thisMessage.Version, err = strconv.Atoi(versionTempString)
	if err != nil {
		return thisMessage, fmt.Errorf("couldn't parse version field [%v] in payload [%v]", payloadSlice[0], payloadString)
	}

	thisMessage.Timestamp = payloadSlice[1]
	thisMessage.Hostname = payloadSlice[2]
	thisMessage.AppName = payloadSlice[3]
	thisMessage.ProcID = payloadSlice[4]
	thisMessage.MessageID = payloadSlice[5]
	thisMessage.Message = payloadSlice[6]

	return thisMessage, nil
}

// ParseMultiple parses a bytes buffer for one or more RFC5424 syslog messages
// and returns an array of them
func ParseMultiple(source io.Reader) (messages []Message, err error) {
	message, err := Parse(source)

	if err != nil {
		if err == io.EOF {
			return messages, errors.New("encountered EOF before first message")
		}
		return messages, err
	}

	for err == nil {
		messages = append(messages, message)
		message, err = Parse(source)
	}

	if err == io.EOF {
		return messages, nil
	}

	return messages, err
}
