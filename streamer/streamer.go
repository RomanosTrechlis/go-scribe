package streamer

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	logServ "github.com/RomanosTrechlis/logStreamer/service/log"
	"github.com/RomanosTrechlis/logStreamer/service/ping"
	p "github.com/RomanosTrechlis/logStreamer/util/format/print"
	"github.com/RomanosTrechlis/logStreamer/util/gserver"
)

const (
	layout string = "02012006150405"
)

// logStreamer holds the servers and other relative information
type logStreamer struct {
	gserver.GRPC
	ticker

	// input stream of protobuf requests
	stream chan pb.LogRequest

	// maximum log file size
	fileSize int64

	// mediator is the address of the mediator middleware
	mediator string

	// rootPath keeps the initial logging path passed by the user
	rootPath string
	// counter counts the requests handled by logStreamer
	counter int64
	stopAll chan struct{}
}

// ticker is responsible for printing
// status updates every few seconds
type ticker struct {
	ticker    *time.Ticker
	stop      chan struct{}
	startTime time.Time
}

// New creates a streamer struct
func New(root string, port int, fileSize int64, mediator, crt, key, ca string) (*logStreamer, error) {
	srv, err := gserver.New(crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}
	return &logStreamer{
		GRPC: gserver.GRPC{
			Server: srv,
			Port:   port,
			Stop:   make(chan struct{}),
		},
		stream:   make(chan pb.LogRequest),
		fileSize: fileSize,
		ticker: ticker{
			ticker: time.NewTicker(20 * time.Second),
			stop:   make(chan struct{}),
		},
		rootPath: root,
		mediator: mediator,
	}, nil
}

// serviceHandler implements the protobuf service
func (s *logStreamer) serviceHandler(stop chan struct{}) {
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
			return
		}
	}
}

// Serve initializes log streamer's servers
func (s *logStreamer) Serve() {
	p.Print("Log streamer is starting...")
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.serviceHandler(s.GRPC.Stop)

	// rpc server
	go gserver.Serve(s.register(), fmt.Sprintf(":%d", s.GRPC.Port), s.GRPC.Server)

	// ticker
	go s.tickerServ()
	<-s.stopAll
}

// Shutdown gracefully stops log streamer from serving
func (s *logStreamer) Shutdown() {
	s.stopAll <- struct{}{}
	p.Print("Initializing shut down, please wait.")
	close(s.GRPC.Stop)
	close(s.ticker.stop)
	s.ticker.ticker.Stop()
	time.Sleep(1 * time.Second)
	p.Print(fmt.Sprintf("Log streamer handled %d requests during %v", s.counter, time.Since(s.startTime)))
	p.Print("Log streamer shut down")
}

func (s *logStreamer) tickerServ() {
	for _ = range s.ticker.ticker.C {
		select {
		case <-s.ticker.stop:
			p.Print("Ticker is stopping...")
			return
		default:
			p.Print(fmt.Sprintf("Log Streamer handled %d requests, so far.", s.counter))
		}
	}
}

func (s *logStreamer) register() func() {
	return func() {
		log := logServ.Logger{Stream: s.stream}
		pb.RegisterLogStreamerServer(s.GRPC.Server, log)

		if s.mediator != "" {
			p := &ping.Pinger{}
			pb.RegisterPingerServer(s.GRPC.Server, p)
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
