# rfc5424 [![GoDoc](https://godoc.org/github.com/tonimelisma/rfc5424?status.svg)](https://pkg.go.dev/mod/github.com/tonimelisma/rfc5424) [![Go Report Card](http://goreportcard.com/badge/tonimelisma/rfc5424)](http://goreportcard.com/report/tonimelisma/rfc5424) ![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/tonimelisma/rfc5424) ![License](https://img.shields.io/badge/license-MIT-blue.svg) [![Build Status](https://github.com/tonimelisma/rfc5424/actions/workflows/go.yml/badge.svg)](https://github.com/tonimelisma/rfc5424/actions) [![Coverage Status](https://img.shields.io/coveralls/github/tonimelisma/rfc5424)](https://coveralls.io/github/tonimelisma/rfc5424?branch=master)
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
		fmt.Printf("%v [%v.%v] %v %v %v: %v\n", message.Timestamp, message.Facility, message.Severity, message.Hostname, message.AppName, message.ProcID, message.Message)
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
This library does not parse this structured data. It is provided as-is in the ``Message`` field of the struct, just like
regular unstructured log messages. I've never actually seen any system use the structured data, and most implementations
break RFC5424 by transmitting unstructured data instead of structured in the seventh field, and thus parsing it as
structured would break compatibility with most systems.

## Maintenance
I consider this library feature-complete for my use cases and don't foresee much activity in the repository.
However, *this software is still actively maintained*. Any issues or PRs will be dealt with in a reasonable
amount of time.