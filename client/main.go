package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/RomanosTrechlis/logStream/api"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("test")
	conn, err := grpc.Dial("127.0.0.1:8080",
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	var i int
	for {
		i++
		c := pb.NewLogStreamerClient(conn)
		req := &pb.LogRequest{
			Filename: "test",
			Path:     "path",
			Line:     fmt.Sprintf("%d: This is a test", i),
		}
		r, err := c.Log(context.Background(), req)
		if err != nil {
			log.Fatalf("failled: %v", err)
		}
		fmt.Printf("%s", r.Res)
	}

}
