# rfc5424 [![GoDoc](https://godoc.org/github.com/tonimelisma/rfc5424?status.svg)](https://pkg.go.dev/mod/github.com/tonimelisma/rfc5424) [![Go Report Card](http://goreportcard.com/badge/tonimelisma/rfc5424)](http://goreportcard.com/report/tonimelisma/rfc5424) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tonimelisma/rfc5424) ![License](https://img.shields.io/badge/license-MIT-blue.svg)
Fast RFC5424 syslog message parser written in Go

## Usage

```go get github.com/tonimelisma/rfc5424```

## Example

```
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
```

For a practical example that parses syslog messages from HTTPS POST bodies as a 
log drain for Heroku, see [https://github.com/tonimelisma/golang-heroku-log-drain](https://github.com/tonimelisma/golang-heroku-log-drain)

## Caveats

RFC5424 defines a way to transmit structured data messages in addition to the more typical free-form text log messages.
This structured data is not parsed in any way. I've never actually seen any system use the data.

## Maintenance

Although I will not update the software, I will attempt to answer opened
issues or PRs in a reasonable timeframe.