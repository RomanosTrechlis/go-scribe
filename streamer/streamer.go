package streamer

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// LogStreamer holds the servers and other relative information
type LogStreamer struct {
	// input stream of protobuf requests
	stream chan pb.LogRequest

	// grpcServer
	grpcServer *grpc.Server
	// GRPC server port
	grpcPort int
	// stopGrpc waits for an empty struct to stop the rpc server.
	stopGrpc chan struct{}

	// maximum log file size
	fileSize int64

	// ticker is responsible for printing
	// status updates every few seconds
	ticker     *time.Ticker
	stopTicker chan struct{}
	startTime  time.Time

	rootPath string
	counter  int64
	stopAll  chan struct{}
}

// New creates a streamer struct
func New(root string, port int, fileSize int64) *LogStreamer {
	return &LogStreamer{
		stream:     make(chan pb.LogRequest),
		grpcServer: grpc.NewServer(),
		grpcPort:   port,
		stopGrpc:   make(chan struct{}),
		fileSize:   fileSize,
		ticker:     time.NewTicker(20 * time.Second),
		stopTicker: make(chan struct{}),
		rootPath:   root,
	}
}

// ServiceHandler implements the protobuf service
func (s *LogStreamer) ServiceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-s.stream:
			s.counter++
			err := handleIncomingRequest(s.rootPath, req.GetPath(), req.GetFilename(), req.GetLine(), s.fileSize)
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
func (s *LogStreamer) Serve() {
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.ServiceHandler(s.stopGrpc)

	// rpc server
	go server(s.stream, fmt.Sprintf(":%d", s.grpcPort), s.grpcServer)

	// ticker
	go s.tickerServ()
	<-s.stopAll
}

// Shutdown gracefully stops log streamer's servers
func (s *LogStreamer) Shutdown() {
	s.stopAll <- struct{}{}
	fmt.Printf("\n%s [INFO] initializing shut down, please wait.\n", PrintTime())
	s.stopGrpc <- struct{}{}
	s.stopTicker <- struct{}{}
	s.ticker.Stop()
	time.Sleep(1 * time.Second)
	fmt.Printf("%s [INFO] Log streamer handled %d requests during %v\n",
		PrintTime(), s.counter, time.Since(s.startTime))
	fmt.Printf("%s [INFO] Log streamer shut down\n", PrintTime())
}

func (s *LogStreamer) tickerServ() {
	for _ = range s.ticker.C {
		select {
		case <-s.stopTicker:
			fmt.Printf("\n%s [INFO] Ticker is stopping...\n", PrintTime())
			return
		default:
			fmt.Printf("%s [INFO] Log Streamer handled %d requests, so far.\n",
				PrintTime(), s.counter)
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

func server(stream chan pb.LogRequest, port string, s *grpc.Server) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log := logger{stream}
	pb.RegisterLogStreamerServer(s, log)

	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
	fmt.Printf("\n%s [INFO] rpc server stopped\n", PrintTime())
}

const (
	//layout string = "2006-01-02T15.04.05Z07.00"
	layout string = "02012006150405"
)

// PrintTime exists for consistency
func PrintTime() string {
	t := time.Now()
	return t.Local().Format(layout)
}
