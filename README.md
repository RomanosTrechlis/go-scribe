# Log Scribe

[![Build Status](https://travis-ci.org/RomanosTrechlis/go-scribe.svg?branch=master)](https://travis-ci.org/RomanosTrechlis/go-scribe)
[![Go Report Card](https://goreportcard.com/badge/github.com/RomanosTrechlis/go-scribe)](https://goreportcard.com/report/github.com/RomanosTrechlis/go-scribe)
[![codecov](https://codecov.io/gh/RomanosTrechlis/go-scribe/branch/master/graph/badge.svg)](https://codecov.io/gh/RomanosTrechlis/go-scribe)

Log Scribe consists of two project that are both based on protocol buffers and gRPC. The first is the Scribe and the second is the Mediator.

## 1. Scribe

The Scribe is the worker that writes log lines to files.

#### Flags

```
Usage of logScribe:
  -ca string
    	certificate authority's certificate
  -console
    	dumps log lines to console
  -crt string
    	host's certificate for secured connections
  -mediator string
    	mediators address if exists, i.e 127.0.0.1:8080
  -path string
    	path for logs to be persisted (default "../logs")
  -pk string
    	host's private key
  -port int
    	port for server to listen to requests (default 8080)
  -pport int
    	port for pprof server (default 1111)
  -pprof
    	additional server for pprof functionality
  -size string
    	max size for individual files, -1B for infinite size (default "1MB")
```

When the mediator flag has value of type host:port then the Scribe calls the Mediator and gets registered.

## 2. Mediator

The Mediator is used as a master node that balances requests for logging to the registered Scribes (workers).

#### Flags
```
Usage of logMediator:
  -ca string
    	certificate authority's certificate
  -crt string
    	host's certificate for secured connections
  -pk string
    	host's private key
  -port int
    	port for mediator server to listen to requests (default 8000)
  -pport int
    	port for pprof server (default 1111)
  -pprof
    	additional server for pprof functionality
```

When Scribes begin to register, the Mediator starts keeping track of which of them are alive doing health checks every five (5) seconds.
The Mediator also keeps track of which Scribe writes what file, in order to prevent two Scribes writing on the same file at the same time, resulting in a panic from one or both.

## TODO

- [ ] add a one-way SSL authentication for the Scribe (or Mediator).
- [ ] create a more robust algorithm for load balancing among the Scribes.
- [ ] investigate the use of sync.Map instead of sync.Mutex.
- [ ] add more test coverage