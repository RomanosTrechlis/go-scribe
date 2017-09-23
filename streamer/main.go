package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/logStream/api"
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
	if err := writeLine(req); err != nil {
		return fmt.Errorf("failed to write line: %v", err)
	}
	return nil
}

func writeLine(req *pb.LogRequest) error {
	logPath := fmt.Sprintf("%s/%s/%s.log", path, req.GetPath(), req.GetFilename())
	info, err := os.Stat(logPath)
	if os.IsNotExist(err) {
		// path doesn't exist and we need to create it.
		err = os.MkdirAll(filepath.Join(path, req.GetPath()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(path, req.GetPath()), err)
		}
	}

	createNewFile, err := fileExceedsMaxSize(info, req)
	if err != nil {
		return fmt.Errorf("failed to rename file: %v", err)
	}
	if createNewFile {
		// re create file if the old has exceeded max size
		err = os.MkdirAll(filepath.Join(path, req.GetPath()), os.ModePerm)
		if err != nil {
			return fmt.Errorf("couldn't create path '%s': %v",
				filepath.Join(path, req.GetPath()), err)
		}
	}

	f, err := os.OpenFile(logPath,
		syscall.O_CREAT|syscall.O_APPEND|syscall.O_WRONLY, os.ModePerm)
	if err != nil {
		return fmt.Errorf("couldn't create to path '%s': %v", logPath, err)
	}
	defer f.Close()

	line := req.GetLine()
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	_, err = f.WriteString(line)
	if err != nil {
		return fmt.Errorf("couldn't write line: %v", err)
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
	flag.IntVar(&maxSize, "size", 1000000,
		"max size in bytes for individual files, -1 for infinite size")
	flag.Parse()

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
