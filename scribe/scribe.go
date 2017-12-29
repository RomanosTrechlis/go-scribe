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
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
)

const (
	layout string = "02012006150405"
)

// LogScribe holds the servers and other relative information
type LogScribe struct {
	// target struct keeps info on where to write logs
	target

	// GRPC server
	gserver.GRPC

	// input stream of protobuf requests
	stream chan pb.LogRequest

	// mediator is the address of the mediator middleware
	mediator string

	// counter counts the requests handled by LogScribe
	counter   int64
	startTime time.Time
	stopAll   chan struct{}
}

// New creates a Scribe struct
func New(root string, port int, fileSize int64, mediator, crt, key, ca string) (*LogScribe, error) {
	srv, err := gserver.New(crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}

	t, err := NewTarget(false, true, "", "", root, fileSize)
	if err != nil {
		return nil, err
	}

	return &LogScribe{
		target: *t,
		GRPC: gserver.GRPC{
			Server: srv,
			Port:   port,
			Stop:   make(chan struct{}),
		},
		stream:   make(chan pb.LogRequest),
		mediator: mediator,
	}, nil
}

// NewScribe returns a log scribe that can be parametrized to either write on files or a db.
func NewScribe(port int, gServer *grpc.Server, isDB bool, dbServer, dbStore, rootPath string, fileSize int64) (*LogScribe, error) {
	t, err := NewTarget(isDB, !isDB, dbServer, dbStore, rootPath, fileSize)
	if err != nil {
		return nil, err
	}

	return &LogScribe{
		target: *t,
		GRPC: gserver.GRPC{
			Server: gServer,
			Port:   port,
			Stop:   make(chan struct{}),
		},
		stream: make(chan pb.LogRequest),
	}, nil
}

// Serve initializes log Scribe's servers
func (s *LogScribe) Serve() {
	p.Print("Log Scribe is starting...")
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.serviceHandler(s.Stop)

	// rpc server
	go gserver.Serve(s.register(), fmt.Sprintf(":%d", s.Port), s.Server)

	<-s.stopAll
	p.Print("gRPC server stopped.")
}

// Shutdown gracefully stops log Scribe from serving
func (s *LogScribe) Shutdown() {
	close(s.stopAll)
	p.Print("Initializing shut down, please wait.")
	close(s.Stop)
	p.Print(fmt.Sprintf("Log Scribe handled %d requests during %v", s.counter, time.Since(s.startTime)))
	p.Print("Log Scribe shut down")
}

// Tick prints a count of requests handled.
func (s *LogScribe) Tick(interval time.Duration) {
	for range time.Tick(interval * time.Second) {
		select {
		case <-s.stopAll:
			p.Print("Tick is stopping")
			return
		default:
			p.Print(fmt.Sprintf("Log Scribe handled %d requests, so far.", s.counter))
		}
	}
}

// serviceHandler implements the protobuf service
func (s *LogScribe) serviceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-s.stream:
			s.counter++
			err := s.handleIncomingRequest(req)
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

func (s *LogScribe) register() func() {
	return func() {
		log := logServ.Logger{Stream: s.stream}
		pb.RegisterLogScribeServer(s.Server, log)

		if s.mediator != "" {
			pinger := &ping.Pinger{}
			pb.RegisterPingerServer(s.Server, pinger)
		}
	}
}

func (s *LogScribe) handleIncomingRequest(r pb.LogRequest) error {
	if s.isDB {
		return handleDBRequest(s.database, r.Filename, r.Line)
	}

	var mu sync.RWMutex
	return handleFileRequest(mu, s.rootPath, r.Path, r.Filename, r.Line, s.fileSize)
}

func handleFileRequest(mu sync.RWMutex, rootPath, path, filename, line string, size int64) error {
	mu.Lock()
	defer mu.Unlock()
	if err := writeLine(rootPath, path, filename, line, size); err != nil {
		return fmt.Errorf("failed to write line: %v", err)
	}
	return nil
}

func handleDBRequest(db *mgo.Database, filename, line string) error {
	return db.C(filename).Insert(line)
}
