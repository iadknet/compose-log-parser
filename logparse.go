package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Jeffail/gabs"
)

// Config struct for storing app config data
type Config struct {
	MessagePath string
	FilterPath  string
	FilterValue string
	JSONOnly    bool
}

var config = &Config{}

func init() {

	messagePath := flag.String(
		"path",
		"message",
		`
Path to the json field to print out as the log message.
Nested paths can be provided using dot notation: 
{ "fields": { "verboseMessage": "A very verbose log message" }}

Can be accessed using the following notation:
docker-compose logs | ./logparse -path=fields.verboseMessage
		`,
	)

	filterPath := flag.String(
		"filterPath",
		"",
		`
Path to the json field to use for message filtering.
Must be combined with -filterValue
{ "level": "ERROR", "message": "Something bad happened" }

docker-compose logs | ./logparse -filterPath=level -filterValue=ERROR
		`,
	)

	filterValue := flag.String(
		"filterValue",
		"",
		`
Value to filter for. 
If -filterPath option is not set, this will search the entire log string.
		`,
	)

	JSONOnly := flag.Bool(
		"JSONOnly",
		false,
		`
Flag to ignore text logs and only show json.

this is a text log line and will be ignored
{"message": "this is a json log and and will be printed"}

docker-compose logs | ./logparse -JSONOnly
		`,
	)

	flag.Parse()

	config.MessagePath = *messagePath
	config.FilterPath = *filterPath
	config.FilterValue = *filterValue
	config.JSONOnly = *JSONOnly
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {

		logLine := scanner.Text()
		composeLogLine := parseComposeLog(logLine)

		if composeLogLine.Display != false {
			fmt.Printf("%-10s | \033[0m %s\n",
				composeLogLine.Service,
				composeLogLine.Message,
			)
		}

	}

	if scanner.Err() != nil {
		panic(scanner.Err)
	}
}

// ComposeLogLine Parse docker-compose log output
type ComposeLogLine struct {
	Service string
	Message string
	Display bool
}

func parseComposeLog(logLine string) ComposeLogLine {
	// fmt.Println(logLine)
	// split on pipe
	serviceMessageSlice := strings.SplitN(logLine, "|", 2)

	var service string
	var message string
	if len(serviceMessageSlice) == 2 {
		service = serviceMessageSlice[0]
		message = serviceMessageSlice[1]
	}

	// extract json from log message
	jsonRegex, err := regexp.Compile(`^[^{]*(.*)[^}]*$`)
	if err != nil {
		fmt.Println("Error building regex")
		panic(err)
	}
	jsonResults := jsonRegex.FindStringSubmatch(message)[1]

	display := true

	if jsonResults != "" {
		parsedMessage, err := gabs.ParseJSON([]byte(jsonResults))

		if err == nil {
			messageData := parsedMessage.Path(config.MessagePath).Data()
			filterData := parsedMessage.Path(config.FilterPath).Data()

			if messageData != nil {
				message = messageData.(string)
				message = strings.Replace(message, "\\n", "\n    ", -1)
			}

			if config.FilterPath != "" && filterData != config.FilterValue {
				display = false
			}

		}
	} else {
		if config.JSONOnly {
			display = false
		} else if config.FilterValue != "" && strings.Contains(message, config.FilterValue) {
			display = false
		}

	}

	return ComposeLogLine{
		Service: service,
		Message: message,
		Display: display,
	}

}
