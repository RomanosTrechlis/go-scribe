package mediator

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	"github.com/RomanosTrechlis/go-scribe/service"
	"github.com/RomanosTrechlis/go-scribe/util/gserver"
	"google.golang.org/grpc"

	p "github.com/RomanosTrechlis/go-scribe/util/format/print"
)

// Mediator grpc server and other relative info
type Mediator struct {
	// mux protects scribes and scribesCon
	// while pinging subscribers.
	mux sync.Mutex
	// scribes has as key the scribe id
	// and value its address.
	scribes map[string]string
	// scribesCon has as key the scribe id
	// and value a valid connection
	scribesCon map[string]*grpc.ClientConn
	// scribesCon has as key the scribe id
	// and value a counter of requestts handled by that id.
	scribesCounter map[string]int64
	// streamResponsibility has as key a character
	// and value a scribe id
	scribeResponsibility map[string]string

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
		scribesCon:           make(map[string]*grpc.ClientConn),
		scribes:              make(map[string]string),
		scribeResponsibility: make(map[string]string),
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

			client := pb.NewLogScribeClient(conn)
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
	for k, v := range m.scribeResponsibility {
		toCheck := strings.ToLower(string(s[0]))
		if toCheck >= k {
			return m.scribesCon[v], nil
		}
		conn = m.scribesCon[v]
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
	if len(m.scribes) == 0 {
		return
	}
	for scribe, addr := range m.scribes {
		ok := m.checkSubscriberConnection(scribe, addr)
		if ok {
			continue
		}
		m.checkSubscriberAlive(scribe, addr)

	}
	m.reCalculateScribeResponsibility()
}

// load balancing
func (m *Mediator) reCalculateScribeResponsibility() {
	r := "abcdefghijklmnopqrstuvwxyz0123456789"
	scribeNum := len(m.scribes)
	if scribeNum == 0 {
		return
	}
	rNum := len(r)

	m.scribeResponsibility = make(map[string]string)
	mid := (rNum - 1) / scribeNum

	val := mid
	for s := range m.scribes {
		p.Print(s)
		m.scribeResponsibility[string(r[val])] = s
		val += mid + 1
	}
}

func (m *Mediator) checkSubscriberConnection(key, val string) bool {
	if _, ok := m.scribesCon[key]; !ok {
		conn, err := createConnection(val)
		if err != nil {
			delete(m.scribes, key)
			delete(m.scribesCon, key)
			p.Print(fmt.Sprintf("Deregistering scribe %s at %s", key, val))
			return true
		}
		m.scribesCon[key] = conn
		return true
	}
	return false
}

func (m *Mediator) checkSubscriberAlive(key, val string) {
	if !m.isSubscriberAlive(m.scribesCon[key]) {
		delete(m.scribes, key)
		delete(m.scribesCon, key)
		p.Print(fmt.Sprintf("Deregistering scribe %s at %s", key, val))
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
		return nil, fmt.Errorf("failed to connect to scribe: %v", err)
	}
	return conn, nil
}

func (m *Mediator) register() func() {
	return func() {
		l := &service.Logger{
			Stream: m.stream,
		}
		pb.RegisterLogScribeServer(m.gRPC.Server, l)

		med := &service.Register{
			Subscribers: m.scribes,
		}
		pb.RegisterRegisterServer(m.gRPC.Server, med)
	}
}
