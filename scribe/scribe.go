package scribe

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	logServ "github.com/RomanosTrechlis/go-scribe/service/log"
	"github.com/RomanosTrechlis/go-scribe/service/ping"
	p "github.com/RomanosTrechlis/go-scribe/util/format/print"
	"github.com/RomanosTrechlis/go-scribe/util/gserver"
)

const (
	layout string = "02012006150405"
)

// logScribe holds the servers and other relative information
type logScribe struct {
	gRPC gserver.GRPC

	// input stream of protobuf requests
	stream chan pb.LogRequest

	// maximum log file size
	fileSize int64

	// mediator is the address of the mediator middleware
	mediator string

	// rootPath keeps the initial logging path passed by the user
	rootPath string
	// counter counts the requests handled by logScribe
	counter int64
	startTime time.Time
	stopAll chan struct{}
}

// New creates a Scribe struct
func New(root string, port int, fileSize int64, mediator, crt, key, ca string) (*logScribe, error) {
	srv, err := gserver.New(crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}
	return &logScribe{
		gRPC: gserver.GRPC{
			Server: srv,
			Port:   port,
			Stop:   make(chan struct{}),
		},
		stream:   make(chan pb.LogRequest),
		fileSize: fileSize,
		rootPath: root,
		mediator: mediator,
	}, nil
}

// serviceHandler implements the protobuf service
func (s *logScribe) serviceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-s.stream:
			s.counter++
			err := handleIncomingRequest(s.rootPath, req.GetPath(),
				req.GetFilename(), req.GetLine(), s.fileSize)
			if err != nil {
				fmt.Printf("hanldeIncomingRequest returned with error: %v", err)
				return
			}
		case <-stop:
			p.Print("serviceHandler stopped")
			return
		}
	}
}

// Serve initializes log Scribe's servers
func (s *logScribe) Serve() {
	p.Print("Log Scribe is starting...")
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.serviceHandler(s.gRPC.Stop)

	// rpc server
	go gserver.Serve(s.register(), fmt.Sprintf(":%d", s.gRPC.Port), s.gRPC.Server)

	<-s.stopAll
	p.Print("gRPC server stopped.")
}

// Shutdown gracefully stops log Scribe from serving
func (s *logScribe) Shutdown() {
	close(s.stopAll)
	p.Print("Initializing shut down, please wait.")
	close(s.gRPC.Stop)
	p.Print(fmt.Sprintf("Log Scribe handled %d requests during %v", s.counter, time.Since(s.startTime)))
	p.Print("Log Scribe shut down")
}

// Tick prints a count of requests handled.
func (s *logScribe) Tick(interval time.Duration) {
	for _ = range time.Tick(interval * time.Second) {
		select {
		case <-s.stopAll:
			p.Print("Tick is stopping")
			return
		default:
			p.Print(fmt.Sprintf("Log Scribe handled %d requests, so far.", s.counter))
		}
	}
}

func (s *logScribe) register() func() {
	return func() {
		log := logServ.Logger{Stream: s.stream}
		pb.RegisterLogScribeServer(s.gRPC.Server, log)

		if s.mediator != "" {
			p := &ping.Pinger{}
			pb.RegisterPingerServer(s.gRPC.Server, p)
		}
	}
}

func handleIncomingRequest(rootPath, path, filename, line string, size int64) error {
	var mu sync.RWMutex
	mu.Lock()
	defer mu.Unlock()
	if err := writeLine(rootPath, path, filename, line, size); err != nil {
		return fmt.Errorf("failed to write line: %v", err)
	}
	return nil
}
