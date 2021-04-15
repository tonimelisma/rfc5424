package rfc5424_test

import (
	"bytes"
	"fmt"

	"github.com/tonimelisma/rfc5424"
)

func ExampleParseMultiple() {
	testMessageBuffer := []byte(`83 <40>1 2012-11-30T06:45:29+00:00 host app web.3 - State changed from starting to up
119 <40>1 2012-11-30T06:45:26+00:00 host app web.3 - Starting process with command 'bundle exec rackup config.ru -p 24405'
`)
	testMessageReader := bytes.NewReader(testMessageBuffer)

	messageArray, err := rfc5424.ParseMultiple(testMessageReader)
	if err != nil {
		fmt.Println("error parsing syslog messages:", err.Error())
	}

	for _, message := range messageArray {
		messageFacility, err := message.FacilityString()
		if err != nil {
			fmt.Println("error parsing message facility:", err.Error())
		}

		messageSeverity, err := message.SeverityString()
		if err != nil {
			fmt.Println("error parsing message severity:", err.Error())
		}

		fmt.Printf("%v [%v.%v] %v %v %v: %v\n", message.Timestamp, messageFacility, messageSeverity, message.Hostname, message.AppName, message.ProcID, message.Message)
	}
	// Output:
	// 2012-11-30T06:45:29+00:00 [syslog.emerg] host app web.3: State changed from starting to up
    // 2012-11-30T06:45:26+00:00 [syslog.emerg] host app web.3: Starting process with command 'bundle exec rackup config.ru -p 24405'
}
