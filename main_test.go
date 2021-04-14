package rfc5424

import (
	"bytes"
	"testing"

	"github.com/go-playground/assert/v2"
)

var testMessage = []byte(`83 <40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up
`)
var testMessageList = []byte(`83 <40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up
119 <40>1 2012-11-30T06:45:26+00:00 host app web.3 - Starting process with command 'bundle exec rackup config.ru -p 24405'
`)

func Test_readLength(t *testing.T) {
	length, err := readLength(bytes.NewReader(testMessage))
	if err != nil {
		t.Error(err.Error())
	}
	assert.Equal(t, 83, length)
	assert.Equal(t, nil, err)
}

func Test_parsePriority(t *testing.T) {
	facility, severity := parsePriority(0)
	assert.Equal(t, 0, facility)
	assert.Equal(t, 0, severity)

	facility, severity = parsePriority(165)
	assert.Equal(t, 20, facility)
	assert.Equal(t, 5, severity)
}

func Test_FacilityString(t *testing.T) {
	var m Message
	m.Facility = 0
	facilityString, err := m.FacilityString()
	assert.Equal(t, nil, err)
	assert.Equal(t, "kernel", facilityString)

	m.Facility = 23
	facilityString, err = m.FacilityString()
	assert.Equal(t, nil, err)
	assert.Equal(t, "local7", facilityString)

	m.Facility = -1
	facilityString, err = m.FacilityString()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", facilityString)

	m.Facility = 24
	facilityString, err = m.FacilityString()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", facilityString)
}

func Test_SeverityString(t *testing.T) {
	var m Message
	m.Severity = 0
	severityString, err := m.SeverityString()
	assert.Equal(t, nil, err)
	assert.Equal(t, "emerg", severityString)

	m.Severity = 7
	severityString, err = m.SeverityString()
	assert.Equal(t, nil, err)
	assert.Equal(t, "debug", severityString)

	m.Severity = -1
	severityString, err = m.SeverityString()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", severityString)

	m.Severity = 8
	severityString, err = m.SeverityString()
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", severityString)
}

func Test_Parse(t *testing.T) {
	msg, err := Parse(bytes.NewReader(testMessage))
	assert.Equal(t, nil, err)
	assert.Equal(t, 40, msg.Priority)
	assert.Equal(t, 5, msg.Facility)
	assert.Equal(t, 0, msg.Severity)
	assert.Equal(t, 1, msg.Version)
	assert.Equal(t, "2012-11-30T06:45:29+00:00", msg.Timestamp)
	assert.Equal(t, "host", msg.Hostname)
	assert.Equal(t, "app", msg.AppName)
	assert.Equal(t, "web.3", msg.ProcID)
	assert.Equal(t, "-", msg.MessageID)
	assert.Equal(t, "State changed from starting to up", msg.Message)
}

func Test_Parse2(t *testing.T) {
	testMessageReader := bytes.NewReader(testMessageList)
	msg, err := Parse(testMessageReader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "State changed from starting to up", msg.Message)

	msg, err = Parse(testMessageReader)
	assert.Equal(t, nil, err)
	assert.Equal(t, `Starting process with command 'bundle exec rackup config.ru -p 24405'`, msg.Message)
}

func Test_ParseMultiple(t *testing.T) {
	msgArray, err := ParseMultiple(bytes.NewReader(testMessageList))
	assert.Equal(t, nil, err)
	assert.Equal(t, "State changed from starting to up", msgArray[0].Message)
	assert.Equal(t, `Starting process with command 'bundle exec rackup config.ru -p 24405'`, msgArray[1].Message)
}
