package mediator

import (
	"fmt"
	"math/rand"
	"strings"
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
	mux sync.Mutex
	// streamers has as key the streamer id
	// and value its address.
	streamers map[string]string
	// streamersCon has as key the streamer id
	// and value a valid connection
	streamersCon map[string]*grpc.ClientConn
	// streamersCon has as key the streamer id
	// and value a counter of requestts handled by that id.
	streamersCounter map[string]int64
	// streamResponsibility has as key a character
	// and value a streamer id
	streamResponsibility map[string]string

	// input stream of protobuf requests
	stream chan pb.LogRequest

	gRPC gserver.GRPC

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
		gRPC: gserver.GRPC{
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
	m.gRPC.Stop = make(chan struct{})
	m.startTime = time.Now()

	// for log service
	go m.serviceHandler(m.gRPC.Stop)
	go gserver.Serve(m.register(), fmt.Sprintf(":%d", m.gRPC.Port), m.gRPC.Server)

	go m.startPingingSubcribers()

	<-m.stopAll
}

// Shutdown gracefully stops mediator from serving
func (m *Mediator) Shutdown() {
	m.stopAll <- struct{}{}
	p.Print("Initializing shut down, please wait.")
	close(m.gRPC.Stop)
	time.Sleep(1 * time.Second)
	p.Print(fmt.Sprintf("Mediator handled %d requests during %v",
		m.counter, time.Since(m.startTime)))
	p.Print("Log Mediator shut down")
}

func (m *Mediator) getConnection(s string) (*grpc.ClientConn, error) {
	m.mux.Lock()
	defer m.mux.Unlock()
	var conn *grpc.ClientConn
	for k, v := range m.streamResponsibility {
		toCheck := strings.ToLower(string(s[0]))
		if toCheck >= k {
			return m.streamersCon[v], nil
		}
		conn = m.streamersCon[v]
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

// load balancing
func (m *Mediator) reCalculateStreamerResponsibility() {
	r := "abcdefghijklmnopqrstuvwxyz0123456789"
	streamerNum := len(m.streamers)
	if streamerNum == 0 {
		return
	}
	rNum := len(r)

	m.streamResponsibility = make(map[string]string)
	mid := (rNum - 1) / streamerNum

	val := mid
	for s := range m.streamers {
		p.Print(s)
		m.streamResponsibility[string(r[val])] = s
		val += mid + 1
	}
	// for k, v := range m.streamResponsibility {
	// 	p.Print(k + " " + v)
	// }
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
		pb.RegisterLogStreamerServer(m.gRPC.Server, l)

		med := &register.Register{
			Subscribers: m.streamers,
		}
		pb.RegisterRegisterServer(m.gRPC.Server, med)
	}
}
