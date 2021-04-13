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

func Test_Parse(t *testing.T) {
	msg, err := Parse(bytes.NewReader(testMessage))
	assert.Equal(t, nil, err)
	assert.Equal(t, "<40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up", msg.Message)
}

func Test_Parse2(t *testing.T) {
	testMessageReader := bytes.NewReader(testMessageList)
	msg, err := Parse(testMessageReader)
	assert.Equal(t, nil, err)
	assert.Equal(t, "<40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up", msg.Message)

	msg, err = Parse(testMessageReader)
	assert.Equal(t, nil, err)
	assert.Equal(t, `<40>1 2012-11-30T06:45:26+00:00 host app web.3 - Starting process with command 'bundle exec rackup config.ru -p 24405'`, msg.Message)
}

func Test_ParseMultiple(t *testing.T) {
	msgArray, err := ParseMultiple(bytes.NewReader(testMessageList))
	assert.Equal(t, nil, err)
	assert.Equal(t, "<40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up", msgArray[0].Message)
}
