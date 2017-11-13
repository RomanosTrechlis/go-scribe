package writer

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/RomanosTrechlis/go-scribe/api"
)

// RPCWriter implements Writer interface
type RPCWriter struct {
	conn     *grpc.ClientConn
	path     string
	filename string
}

type builderImpl struct {
	path     string
	filename string
	address  string
	port     int
	cert     string
	key      string
	ca       string
}

// Builder interface holds the option methods
// in order to create an RPCWriter.
type Builder interface {
	WithPath(path string) Builder
	WithFilename(filename string) Builder
	WithSecurity(cert, key, ca string) Builder
	Build() (*RPCWriter, error)
}

// NewBuilder creates a Builder interface with parameters
// for the creation of RPCWriter when Build method gets called.
func NewBuilder(address string, port int) Builder {
	return builderImpl{
		address:  address,
		port:     port,
		path:     "temp",
		filename: "temporary",
	}
}

// WithFilename adds a non default filename parameter to builder
func (b builderImpl) WithFilename(filename string) Builder {
	b.filename = filename
	return b
}

// WithPath adds a non default path parameter to builder
func (b builderImpl) WithPath(path string) Builder {
	b.path = path
	return b
}

// WithSecurity adds parameters necessary for ssl authentication
func (b builderImpl) WithSecurity(cert, key, ca string) Builder {
	b.ca = ca
	b.cert = cert
	b.key = key
	return b
}

// Build creates a new RPCWriter given the Builder parameters.
func (b builderImpl) Build() (*RPCWriter, error) {
	return New(b.path, b.filename, b.address, b.port, b.cert, b.key, b.ca)
}

// New creates an RPCWriter
func New(path, filename, address string, port int, cert, key, ca string) (*RPCWriter, error) {
	s := scribe{
		address: address,
		port:    port,
	}
	config := gRPCConnectionConfig{
		cert: cert,
		key:  key,
		ca:   ca,
	}

	conn, err := createConnection(config, s)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to scribe: %v", err)
	}

	return &RPCWriter{
		conn:     conn,
		path:     path,
		filename: filename,
	}, nil
}

type gRPCConnectionConfig struct {
	cert string
	key  string
	ca   string
}

type scribe struct {
	address string
	port    int
}

func createConnection(sec gRPCConnectionConfig, sc scribe) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	server := fmt.Sprintf("%s:%d", sc.address, sc.port)
	if sec.cert != "" && sec.key != "" && sec.ca != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(sec.cert, sec.key)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(sec.ca)
		if err != nil {
			return nil, fmt.Errorf("could not read ca certificate: %s", err)
		}

		// Append the certificates from the CA
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return nil, fmt.Errorf("failed to append ca certs")

		}

		creds := credentials.NewTLS(&tls.Config{
			ServerName:   sc.address, // NOTE: this is required!
			Certificates: []tls.Certificate{certificate},
			RootCAs:      certPool,
		})

		// Create a connection with the TLS credentials
		conn, err = grpc.Dial(server,
			grpc.WithTransportCredentials(creds),
			grpc.WithTimeout(1*time.Second))
		if err != nil {
			return nil, fmt.Errorf("did not connect: %v", err)
		}
		return conn, nil
	}
	var err error
	conn, err = grpc.Dial(server,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, nil
}

// Write implements the Write method of Writer interface.
// Inside this method there is a call to scribe.
func (w *RPCWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	c := pb.NewLogScribeClient(w.conn)
	req := &pb.LogRequest{
		Filename: w.filename,
		Path:     w.path,
		Line:     string(p),
	}
	r, err := c.Log(context.Background(), req)
	if err != nil {
		return 0, fmt.Errorf("failled to write bytes: %v", err)
	}
	if r.GetRes() != "true" {
		return 0, fmt.Errorf("failled to write bytes: %v", err)
	}
	return n, nil
}

// NewLogger creates a *Logger with default prefix and flags and a new RPCWriter
func NewLogger(path, filename, address string, port int, cert, key, ca string) (*log.Logger, error) {
	w, err := New(path, filename, address, port, cert, key, ca)
	if err != nil {
		return nil, err
	}
	flag := log.Ldate | log.Ltime | log.Lshortfile
	return log.New(w, "", flag), nil
}
