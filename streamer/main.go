package main

import (
	"golang.org/x/net/context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type logger struct{}

// Log is the ptotobuf service implementation
func (l logger) Log(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	countRequests++
	stream <- in
	return &pb.LogResponse{Res: "handling"}, nil
}

func server(port string, s *grpc.Server) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	log := logger{}
	pb.RegisterLogStreamerServer(s, log)

	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
	fmt.Printf("\n%s [INFO] rpc server stopped\n", printTime())
}

func channelHandler(stop chan struct{}, s *grpc.Server) {
	for {
		select {
		case req := <-stream:
			err := handleIncomingRequest(req)
			if err != nil {
				fmt.Printf("hanldeIncomingRequest returned with error: %v", err)
				shutdown(s)
				return
			}
		case <-stop:
			shutdown(s)
			return
		}
	}
}

func handleIncomingRequest(req *pb.LogRequest) error {
	if console {
		fmt.Printf("%s/%s.log: %s\n", req.GetPath(), req.GetFilename(), req.GetLine())
	}

	var mu sync.RWMutex
	mu.Lock()
	defer mu.Unlock()
	if err := writeLine(req, path, maxSize); err != nil {
		return fmt.Errorf("failed to write line: %v", err)
	}
	return nil
}

func shutdown(s *grpc.Server) {
	s.GracefulStop()
	close(stream)
	fmt.Printf("%s [INFO] rpc channel handler is closing\n", printTime())
}

func printTime() string {
	t := time.Now()
	return t.Local().Format(layout)
}

var (
	port          int
	stream        chan *pb.LogRequest
	path          string
	maxSize       int
	countRequests int64
	pprofInfo     bool
	pport         int
	console       bool
)

const (
	layout string = time.RFC3339
)

func init() {
	fmt.Printf("%s [INFO] Log streamer is starting...\n", printTime())
	// rpc server listening port
	flag.IntVar(&port, "port", 8080, "port for server to listen to requests")
	// enable/disable pprof functionality
	flag.BoolVar(&pprofInfo, "pprof", false,
		"additional server for pprof functionality")
	// enable/disable console dumps
	flag.BoolVar(&console, "console", false, "dumps log lines to console")
	// pprof port for http server
	flag.IntVar(&pport, "pport", 1111, "port for pprof server")
	// path must already exist
	flag.StringVar(&path, "path", "../logs", "path for logs to be persisted")
	// the size of log files before they get renamed for storing purposes.
	size := flag.String("size", "1MB",
		"max size for individual files, -1B for infinite size")
	flag.Parse()

	i, err := lexicalToNumber(*size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse size input to bytes: %v", err)
		os.Exit(2)
	}
	maxSize = i

	// prints some logo and info
	printLogo()
	infoBlock(port, pport, maxSize, path, pprofInfo)
}

func main() {
	countRequests = 0
	start := time.Now()

	// validate path passed
	if err := checkPath(path); err != nil {
		fmt.Printf("path passed is not valid: %v\n", err)
		return
	}

	// stream is the channel that takes LogRequests.
	stream = make(chan *pb.LogRequest)
	// stopRPC waits for an empty struct to stop the rpc server.
	stopRPC := make(chan struct{})
	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	// grpc server is created in the main
	// to pass it around as parameter and
	// stop it gracefully
	s := grpc.NewServer()

	// go func listens to stream and stop channels
	go channelHandler(stopRPC, s)

	// rpc server
	go server(fmt.Sprintf(":%d", port), s)

	// starts pprof server for debuging purposes.
	if pprofInfo {
		go pprofServer(pport)
	}

	// ticker for printing some information on the requests handled
	stopTicker := make(chan struct{})
	go func() {
		for _ = range time.Tick(20 * time.Second) {
			select {
			case <-stopTicker:
				fmt.Printf("%s [INFO] Ticker is stopping...\n", printTime())
				return
			default:
				fmt.Printf("%s [INFO] Log Streamer handled %d requests, so far.\n",
					printTime(), countRequests)
			}
		}
	}()

	<-stopAll
	fmt.Printf("\n%s [INFO] initializing shut down, please wait.\n", printTime())
	stopRPC <- struct{}{}
	stopTicker <- struct{}{}
	time.Sleep(1 * time.Second)
	fmt.Printf("%s [INFO] Log streamer handled %d requests during %v\n",
		printTime(), countRequests, time.Since(start))
	fmt.Printf("%s [INFO] Log streamer shut down\n", printTime())
}
