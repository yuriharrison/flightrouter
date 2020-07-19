# Flight Router

![Build](https://github.com/yuriharrison/flightrouter/workflows/CI/badge.svg?style=flat-square)

**Example data:** ./fixtures/flights.csv

### Makefile

```
$make
build                Build application
format               Run GOFMT to format code
help                 This help
lint                 Run GOVET and GOLINT to check code quality
test                 Run tests
```

## Usage

### CLI

```shell
$ ./flightrouter ../fixtures/flights.csv
Importing file...
Search for flight (e.g. GRU-APQ): GRU-APQ
Cheapest route available: GRU-PLU-APQ -> $25.00
```

### API

```
$ ./flightrouter api --help

Starts the Web API

Usage:
   flightrouter {flags}

Flags:
   -d, --data                    csv file to pre fetch data (default: NULL)
   -h, --help                    displays usage information of the application or a command (default: false)
   -p, --port                    port (default: 8080)
```

- GET `/flights/search/:route`
  - `curl localhost:8080/flights/search/GRU-APQ`
- PUT `/flight`
  - `curl -X PUT -d "origin=GRU" -d "destination=APQ" -d "price=10" localhost:8080/flight`
- DELETE `/flight/:route`
  - `curl -X DELETE -w "%{http_code}" localhost:8080/flight/GRU-APQ`
- PUT `/flights/import/csv`
  - `curl -X PUT -F "document=@./fixtures/flights.csv;type=text/csv" -w "%{http_code}" localhost:8080/flights/import/csv`
