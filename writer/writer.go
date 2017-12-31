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
	target
	conn *grpc.ClientConn
}

type builderImpl struct {
	target
	address string
	port    int
	cert    string
	key     string
	ca      string
}

type target struct {
	isDir    bool
	path     string
	filename string

	isDB     bool
	dbServer string
	dbName   string
}

// Builder interface holds the option methods
// in order to create an RPCWriter.
type Builder interface {
	WithPath(path string) Builder
	WithFilename(filename string) Builder
	WithDatabase(dbServer, dbName string) Builder
	WithSecurity(cert, key, ca string) Builder
	Build() (*RPCWriter, error)
}

// NewBuilder creates a Builder interface with parameters
// for the creation of RPCWriter when Build method gets called.
func NewBuilder(address string, port int) Builder {
	return builderImpl{
		address: address,
		port:    port,
		target: target{
			isDir: false,
			isDB:  false,
		},
	}
}

// WithFilename adds a non default filename parameter to builder
func (b builderImpl) WithFilename(filename string) Builder {
	b.filename = filename
	b.isDir = true
	return b
}

// WithPath adds a non default path parameter to builder
func (b builderImpl) WithPath(path string) Builder {
	b.path = path
	b.isDir = true
	return b
}

// WithDatabase adds database support
func (b builderImpl) WithDatabase(dbServer, dbName string) Builder {
	b.dbServer = dbServer
	b.dbName = dbName
	b.isDB = true
	b.isDir = false
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
	return newRPCWriter(b)
}

// new creates an RPCWriter
func newRPCWriter(b builderImpl) (*RPCWriter, error) {
	s := scribe{
		address: b.address,
		port:    b.port,
	}

	conn, err := createConnection(b.cert, b.key, b.ca, s)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection to scribe: %v", err)
	}

	if b.isDB {
		b.isDir = false
	}
	return &RPCWriter{
		conn:   conn,
		target: b.target,
	}, nil
}

type scribe struct {
	address string
	port    int
}

func createConnection(cert, key, ca string, sc scribe) (*grpc.ClientConn, error) {
	var conn *grpc.ClientConn
	server := fmt.Sprintf("%s:%d", sc.address, sc.port)
	if cert != "" && key != "" && ca != "" {
		// Load the client certificates from disk
		certificate, err := tls.LoadX509KeyPair(cert, key)
		if err != nil {
			return nil, fmt.Errorf("could not load client key pair: %s", err)
		}

		// Create a certificate pool from the certificate authority
		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(ca)
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
	w, err := NewBuilder(address, port).
		WithFilename(filename).
		WithPath(path).
		WithSecurity(cert, key, ca).
		Build()
	if err != nil {
		return nil, err
	}
	flag := log.Ldate | log.Ltime | log.Lshortfile
	return log.New(w, "", flag), nil
}

// NewLogger creates a *Logger with default prefix and flags and a new RPCWriter
func NewDBLogger(dbServer, dbName, address string, port int, cert, key, ca string) (*log.Logger, error) {
	w, err := NewBuilder(address, port).
		WithDatabase(dbServer, dbName).
		WithSecurity(cert, key, ca).
		Build()
	if err != nil {
		return nil, err
	}
	flag := log.Ldate | log.Ltime | log.Lshortfile
	return log.New(w, "", flag), nil
}
