package main

import (
	"fmt"
	"log"
	"net"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.GreetServiceServer
}

func main() {
	fmt.Println("Hello World")

	mux, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	fmt.Println(mux)

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	err = s.Serve(mux)
	if err != nil {
		log.Fatal("failed to listen", err)
	}
}
