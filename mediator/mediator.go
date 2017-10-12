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
	// mux protectes streamers and streamersCon
	// while pinging subscribers.
	mux                  sync.Mutex
	streamers            map[string]string
	streamersCon         map[string]*grpc.ClientConn
	streamResponsibility map[string]string

	// input stream of protobuf requests
	stream chan pb.LogRequest

	gserver.GRPC

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
		stream: make(chan pb.LogRequest),
		GRPC: gserver.GRPC{
			Server: srv,
			Port:   port,
		},
		streamersCon:         make(map[string]*grpc.ClientConn),
		streamers:            make(map[string]string),
		streamResponsibility: make(map[string]string),
	}
	return m, nil
}

// serviceHandler implements the protobuf service
func (m *Mediator) serviceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-m.stream:
			conn, err := m.getConnection(req.GetFilename())
			if err != nil {
				p.Print(err.Error())
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
	m.GRPC.Stop = make(chan struct{})
	m.startTime = time.Now()

	// for log service
	go m.serviceHandler(m.GRPC.Stop)
	go gserver.Serve(m.register(), fmt.Sprintf(":%d", m.GRPC.Port), m.GRPC.Server)

	go m.startPingingSubcribers()

	<-m.stopAll
}

// Shutdown gracefully stops mediator from serving
func (m *Mediator) Shutdown() {
	m.stopAll <- struct{}{}
	p.Print("Initializing shut down, please wait.")
	close(m.GRPC.Stop)
	time.Sleep(1 * time.Second)
	p.Print(fmt.Sprintf("Mediator handled %d requests during %v",
		m.counter, time.Since(m.startTime)))
	p.Print("Log Mediator shut down")
}

func (m *Mediator) getConnection(s string) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	for _, v := range m.streamersCon {
		conn = v
		break
	}
	return conn, nil
}

func (m *Mediator) startPingingSubcribers() {
	for _ = range time.Tick(5 * time.Second) {
		m.mux.Lock()
		m.pingSubscribers()
		m.mux.Unlock()
	}
}

func (m *Mediator) pingSubscribers() {
	if len(m.streamers) == 0 {
		return
	}
	for streamer, addr := range m.streamers {
		ok := m.checkSubscriberConnection(streamer, addr)
		if ok {
			continue
		}
		m.checkSubscriberAlive(streamer, addr)

	}
	m.reCalculateStreamerResponsibility()
}

// testing load balancing
func (m *Mediator) reCalculateStreamerResponsibility() {
	r := "abcdefghijklmnopqrstuvwxyz0123456789"
	streamerNum := len(m.streamers)
	if streamerNum == 0 {
		return
	}
	rNum := 36

	m.streamResponsibility = make(map[string]string)
	mid := (rNum - 1) / streamerNum

	val := mid
	p.Print(fmt.Sprintf("slicing at %d with %s value and index of %d", val, r, rNum))
	for s := range m.streamers {
		p.Print(s)
		m.streamResponsibility[string(r[val])] = s
		val += mid
	}

	for k, v := range m.streamResponsibility {
		p.Print(k + " " + v)
	}
}

func (m *Mediator) checkSubscriberConnection(key, val string) bool {
	if _, ok := m.streamersCon[key]; !ok {
		conn, err := createConnection(val)
		if err != nil {
			delete(m.streamers, key)
			delete(m.streamersCon, key)
			p.Print(fmt.Sprintf("Deregistering streamer %s at %s", key, val))
			return true
		}
		m.streamersCon[key] = conn
		return true
	}
	return false
}

func (m *Mediator) checkSubscriberAlive(key, val string) {
	if !m.isSubscriberAlive(m.streamersCon[key]) {
		delete(m.streamers, key)
		delete(m.streamersCon, key)
		p.Print(fmt.Sprintf("Deregistering streamer %s at %s", key, val))
	}
}

func (m *Mediator) isSubscriberAlive(conn *grpc.ClientConn) bool {
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
		pb.RegisterLogStreamerServer(m.GRPC.Server, l)

		med := &register.Register{
			Subscribers: m.streamers,
		}
		pb.RegisterRegisterServer(m.GRPC.Server, med)
	}
}
