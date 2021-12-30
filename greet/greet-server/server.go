package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/liridonrama/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (s *server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	fmt.Printf("Greet function was invoked with %v", req)

	fN, lN := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()

	result := fmt.Sprintf("Hello, %v %v", fN, lN)

	res := &greetpb.GreetResponse{
		Result: result,
	}

	return res, nil
}

func (s *server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimesFunction was invoked:", req)

	fN := req.Greeting.GetFirstName()

	for i := 0; i < 10; i++ {
		res := &greetpb.GreetManyTimesResponse{
			Result: fmt.Sprintf("Hello %v number %v", fN, strconv.Itoa(i+1)),
		}

		stream.Send(res)
		time.Sleep(time.Second)
	}

	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet was invoked")
	result := "Hello "

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&greetpb.LongGreetResponse{
				Result: result,
			})
		}

		if err != nil {
			log.Fatalln("Error while reading stream:", err)
		}

		fN := req.GetGreeting().GetFirstName()
		result += fN + "! \n"
	}
}

func (*server) GreetEveryone(stream greetpb.GreetService_GreetEveryoneServer) error {
	fmt.Println("GreetEveryone was invoked")

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalln("Error while reading stream:", err)
			return err
		}

		fN := msg.GetGreeting()

		result := fmt.Sprintf("Hello, %v!", fN)
		err = stream.Send(&greetpb.GreetEveryoneResponse{
			Result: result,
		})

		if err != nil {
			log.Fatalln("Error while sending data to client:", err)
		}
	}
}

func (s *server) GreetWithDeadline(ctx context.Context, req *greetpb.GreetWithDeadlineRequest) (*greetpb.GreetWithDeadlineResponse, error) {
	fmt.Println("GreetWithDeadline function was invoked with:", req)

	for i := 0; i < 3; i++ {
		if ctx.Err() == context.Canceled {
			fmt.Println("Client canceled request")

			return nil, status.Error(codes.Canceled, "the client canceled the request")
		}

		time.Sleep(time.Second)
	}

	fN, lN := req.GetGreeting().GetFirstName(), req.GetGreeting().GetLastName()

	result := fmt.Sprintf("Hello, %v %v", fN, lN)

	res := &greetpb.GreetWithDeadlineResponse{
		Result: result,
	}

	return res, nil
}

func main() {
	mux, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatalln("failed to listen", err)
	}

	s := grpc.NewServer()

	greetpb.RegisterGreetServiceServer(s, &server{})

	err = s.Serve(mux)
	if err != nil {
		log.Fatalln("failed to listen", err)
	}
}
