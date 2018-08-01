package gserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net"

	p "github.com/RomanosTrechlis/go-scribe/internal/util/format/print"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// New creates a grpc server with or without SSL
func New(crt, key, ca string) (*grpc.Server, error) {
	if crt == "" || key == "" || ca == "" {
		return grpc.NewServer(), nil
	}

	// one way ssl
	if crt != "" && key != "" && ca == "" {
		// todo(romanos): yet to be implemented
	}

	// two way ssl
	p.Print("Log streamer will start with TLS")
	// Load the certifmediatoricates from disk
	certificate, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, fmt.Errorf("could not load server key pair: %s", err)
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	caBytes, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, fmt.Errorf("could not read ca certifiapicate: %s", err)
	}

	// Append the client certificates from the CA
	if ok := certPool.AppendCertsFromPEM(caBytes); !ok {
		return nil, fmt.Errorf("failed to append client certs")
	}

	// Create the TLS credentials
	creds := credentials.NewTLS(&tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    certPool,
	})

	// Create the gRPC server with the credentials
	return grpc.NewServer(grpc.Creds(creds)), nil
}

// Serve listen for LogRequests
func Serve(register func(), addr string, s *grpc.Server) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	register()

	reflection.Register(s)
	err = s.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
	}
	p.Print("RPC server stopped")
}

// GRPC holds info about the gRPC server.
type GRPC struct {
	// grpcServer
	Server *grpc.Server
	// GRPC server port
	Port int
	// stopGrpc waits for an empty struct to stop the rpc server.
	Stop chan struct{}
}
