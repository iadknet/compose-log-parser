# Compose Log Parser

A simple go app that can take Docker Compose log output and parse json logs for better readability.

## Examples

Simple example will print out non-json log messages and look for `message` elements in json logs

```bash
docker-compose logs -f | logparse
```

Do not print non-json logs

```bash
docker-compose logs -f | logparse -JSONOnly
```

Print a specific JSON field

```bash
# Sample JSON
{ "fields" :{ "verboseMessage": "A very verbose log message" }}

# Extract using dot path notation
docker-compose logs | logparse -path=fields.verboseMessage
```

Filter based on field value
```bash

# Sample JSON
{ "level": "ERROR", "message": "Something bad happened" }

docker-compose logs | ./logparse -filterPath=level -filterValue=ERROR
```

Help
```bash
$ ./logparse -h

Usage of ./logparse:
  -JSONOnly

        Flag to ignore text logs and only show json.

        this is a text log line and will be ignored
        {"message": "this is a json log and and will be printed"}

        docker-compose logs | ./logparse -JSONOnly

  -filterPath string

        Path to the json field to use for message filtering.
        Must be combined with -filterValue
        { "level": "ERROR", "message": "Something bad happened" }

        docker-compose logs | ./logparse -filterPath=level -filterValue=ERROR

  -filterValue string

        Value to filter for.
        If -filterPath option is not set, this will search the entire log string.

  -path string

        Path to the json field to print out as the log message.
        Nested paths can be provided using dot notation:
        { "fields": { "verboseMessage": "A very verbose log message" }}

        Can be accessed using the following notation:
        docker-compose logs | ./logparse -path=fields.verboseMessage
                         (default "message")
```
