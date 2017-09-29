# Log Streamer

Log Streamer is a server that handles logging from remote clients.

Log Streamer uses protocol buffer and gRPC to communicate.

## Usage

**rootPath**: is the chosen path for saving logs

**port**: listening port

**maxFileSize**: the maximum size of log files before the get stored (rename)

**cert, key, ca**: if empty the server starts without TLS, else it starts with TLS

```go
s, err := streamer.New(rootPath, port, maxFileSize, cert, key, ca)
if err != nil {
  return
}

defer s.Shutdown()
go s.Serve()
```
*NOTE*: shutting down the server is user's responsibility.  

## Flags
```
Usage of logStreamer:
  -ca string
    	certificate authority's certificate
  -console
    	dumps log lines to console
  -crt string
    	host's certificate for secured connections
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
