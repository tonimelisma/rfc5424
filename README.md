# rfc5424
Fast RFC5424 syslog message parser written in Go

## Usage

```go get github.com/tonimelisma/rfc5424```

## Example

TBD fill in an example here

For a practical example that parses syslog messages from HTTPS POST bodies as a 
log drain for Heroku, see [https://github.com/tonimelisma/golang-heroku-log-drain](https://github.com/tonimelisma/golang-heroku-log-drain)

## Caveats

RFC5424 defines a way to transmit structured data messages in addition to the more typical free-form text log messages.
This structured data is not parsed in any way. I've never actually seen any system use the data.

## Maintenance

Although I will not update the software, I will attempt to answer opened
issues or PRs in a reasonable timeframe.