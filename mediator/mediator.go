package mediator

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"golang.org/x/net/context"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	logServ "github.com/RomanosTrechlis/logStreamer/service/log"
	"github.com/RomanosTrechlis/logStreamer/service/register"
	"github.com/RomanosTrechlis/logStreamer/util/gserver"
	"google.golang.org/grpc"

	p "github.com/RomanosTrechlis/logStreamer/util/format/print"
)

// Mediator grpc server and other relative info
type Mediator struct {
	// mux protectes streamers and streamersCon to be
	// accessed during pinging.
	mux          sync.Mutex
	streamers    map[string]string
	streamersCon map[string]*grpc.ClientConn

	// input stream of protobuf requests
	stream chan pb.LogRequest

	// this can be a composite with the corresponding fields of LogStreamer
	grpcServer *grpc.Server
	grpcPort   int
	stopGrpc   chan struct{}

	// these can also be composed
	startTime time.Time
	counter   int64
	stopAll   chan struct{}
}

// New creates a new mediator
func New(port int, crt, key, ca string) (*Mediator, error) {
	srv, err := gserver.New(crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}

	m := &Mediator{
		stream:       make(chan pb.LogRequest),
		grpcServer:   srv,
		grpcPort:     port,
		streamersCon: make(map[string]*grpc.ClientConn),
		streamers:    make(map[string]string),
	}
	return m, nil
}

// ServiceHandler implements the protobuf service
func (m *Mediator) ServiceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-m.stream:
			var conn *grpc.ClientConn
			for _, v := range m.streamersCon {
				conn = v
				break
			}

			client := pb.NewLogStreamerClient(conn)
			client.Log(context.Background(), &req)
			m.counter++
		case <-stop:
			return
		}
	}
}

// Serve starts mediator server
func (m *Mediator) Serve() {
	p.Print("Log Mediator is starting...")
	m.stopAll = make(chan struct{})
	m.stopGrpc = make(chan struct{})
	m.startTime = time.Now()

	// for log service
	go m.ServiceHandler(m.stopGrpc)
	go gserver.Serve(m.register(), fmt.Sprintf(":%d", m.grpcPort), m.grpcServer)

	go m.pingSubscribers()

	<-m.stopAll
}

func (m *Mediator) pingSubscribers() {
	for _ = range time.Tick(5 * time.Second) {
		m.mux.Lock()
		for key, val := range m.streamers {
			if _, ok := m.streamersCon[key]; !ok {
				conn, err := createConnection(val)
				if err != nil {
					delete(m.streamers, key)
					continue
				}
				m.streamersCon[key] = conn
			}
			if !isSubscriberLive(m.streamersCon[key]) {
				delete(m.streamers, key)
				delete(m.streamersCon, key)
				p.Print(fmt.Sprintf("Deregistering streamer %s at %s", key, val))
			}
		}
		m.mux.Unlock()
	}
}

func isSubscriberLive(conn *grpc.ClientConn) bool {
	c := pb.NewPingerClient(conn)
	req := &pb.PingRequest{
		A: rand.Int31(),
		B: rand.Int31(),
	}
	r, err := c.Ping(context.Background(), req)
	if err != nil {
		return false
	}
	if r.GetRes() != req.GetA()*req.GetB() {
		return false
	}
	return true
}

func createConnection(addr string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to streamer: %v", err)
	}
	return conn, nil
}

func (m *Mediator) register() func() {
	return func() {
		l := &logServ.Logger{
			Stream: m.stream,
		}
		pb.RegisterLogStreamerServer(m.grpcServer, l)

		med := &register.Register{
			Subscribers: m.streamers,
		}
		pb.RegisterRegisterServer(m.grpcServer, med)
	}
}

// Shutdown gracefully stops mediator from serving
func (m *Mediator) Shutdown() {
	m.stopAll <- struct{}{}
	p.Print("Initializing shut down, please wait.")
	m.stopGrpc <- struct{}{}
	time.Sleep(1 * time.Second)
	p.Print(fmt.Sprintf("Log streamer handled %d requests during %v",
		m.counter, time.Since(m.startTime)))
	p.Print("Log Mediator shut down")
}
