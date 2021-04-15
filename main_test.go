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
	assert.Equal(t, "kernel", facility)
	assert.Equal(t, "emerg", severity)

	facility, severity = parsePriority(165)
	assert.Equal(t, "local4", facility)
	assert.Equal(t, "notice", severity)
}

func Test_FacilityString(t *testing.T) {
	facilityString, err := parseFacility(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, "kernel", facilityString)

	facilityString, err = parseFacility(23)
	assert.Equal(t, nil, err)
	assert.Equal(t, "local7", facilityString)

	facilityString, err = parseFacility(-1)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", facilityString)

	facilityString, err = parseFacility(24)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", facilityString)
}

func Test_SeverityString(t *testing.T) {
	severityString, err := parseSeverity(0)
	assert.Equal(t, nil, err)
	assert.Equal(t, "emerg", severityString)

	severityString, err = parseSeverity(7)
	assert.Equal(t, nil, err)
	assert.Equal(t, "debug", severityString)

	severityString, err = parseSeverity(-1)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", severityString)

	severityString, err = parseSeverity(8)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "", severityString)
}

func Test_Parse(t *testing.T) {
	msg, err := Parse(bytes.NewReader(testMessage))
	assert.Equal(t, nil, err)
	assert.Equal(t, 40, msg.Priority)
	assert.Equal(t, "syslog", msg.Facility)
	assert.Equal(t, "emerg", msg.Severity)
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
