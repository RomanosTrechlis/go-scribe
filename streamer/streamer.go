package streamer

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

const (
	logLayout string = "2006-01-02T15.04.05Z07.00"
	layout    string = "02012006150405"
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
func New(root string, port int, fileSize int64, crt, key, ca string) (*LogStreamer, error) {
	var srv *grpc.Server
	if crt != "" && key != "" && ca != "" {
		fmt.Printf("%s [INFO] Log streamer will start with TLS\n", printTime(logLayout))
		// Load the certificates from disk
		certificate, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, fmt.Errorf("could not load server key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(ca)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %s", err)
		}

		// Append the client certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, fmt.Errorf("failed to append client certs")
		}

		// Create the TLS credentials
		creds := credentials.NewTLS(&tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			ClientCAs:    certPool,
		})

		// Create the gRPC server with the credentials
		srv = grpc.NewServer(grpc.Creds(creds))
	} else {
		srv = grpc.NewServer()
	}

	return &LogStreamer{
		stream:     make(chan pb.LogRequest),
		grpcServer: srv,
		grpcPort:   port,
		stopGrpc:   make(chan struct{}),
		fileSize:   fileSize,
		ticker:     time.NewTicker(20 * time.Second),
		stopTicker: make(chan struct{}),
		rootPath:   root,
	}, nil
}

// ServiceHandler implements the protobuf service
func (s *LogStreamer) ServiceHandler(stop chan struct{}) {
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
func (s *LogStreamer) Serve() {
	fmt.Printf("%s [INFO] Log streamer is starting...\n", printTime(logLayout))
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.ServiceHandler(s.stopGrpc)

	// rpc server
	go server(s.stream, fmt.Sprintf("127.0.0.1:%d", s.grpcPort), s.grpcServer)

	// ticker
	go s.tickerServ()
	<-s.stopAll
}

// Shutdown gracefully stops log streamer from serving
func (s *LogStreamer) Shutdown() {
	s.stopAll <- struct{}{}
	fmt.Printf("\n%s [INFO] initializing shut down, please wait.\n",
		printTime(logLayout))
	s.stopGrpc <- struct{}{}
	s.stopTicker <- struct{}{}
	s.ticker.Stop()
	time.Sleep(1 * time.Second)
	fmt.Printf("%s [INFO] Log streamer handled %d requests during %v\n",
		printTime(logLayout), s.counter, time.Since(s.startTime))
	fmt.Printf("%s [INFO] Log streamer shut down\n", printTime(logLayout))
}

func (s *LogStreamer) tickerServ() {
	for _ = range s.ticker.C {
		select {
		case <-s.stopTicker:
			fmt.Printf("\n%s [INFO] Ticker is stopping...\n", printTime(logLayout))
			return
		default:
			fmt.Printf("%s [INFO] Log Streamer handled %d requests, so far.\n",
				printTime(logLayout), s.counter)
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

func server(stream chan pb.LogRequest, addr string, s *grpc.Server) {
	lis, err := net.Listen("tcp", addr)
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
	fmt.Printf("\n%s [INFO] rpc server stopped\n", printTime(logLayout))
}
