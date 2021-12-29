package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	"github.com/liridonrama/grpc-go-course/calculator/calculatorpb"
	"google.golang.org/grpc"
)

type server struct {
	calculatorpb.UnimplementedCalculatorServiceServer
}

func (s *server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	vx := req.Sum.GetValues()

	result := 0.0

	for _, v := range vx {
		result += v
	}

	sR := &calculatorpb.SumResponse{
		Result: result,
	}

	return sR, nil
}

func (s *server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	number := req.GetNumber()
	var k int32 = 2

	for number > 1 {
		if number%k == 0 {
			fmt.Println("this is a factor:", number)
			err := stream.Send(&calculatorpb.PrimeNumberDecompositionResponse{
				PrimeNumber: k,
			})
			if err != nil {
				fmt.Println("error while trying to extract prime numbers:", err)
				return err
			}

			number /= k
		} else {
			k++
		}
	}

	return nil
}

func (s *server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	var num float64
	msgCount := 0

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error while trying to read from stream:", err)
		}

		num += msg.GetNumber()
		msgCount++
	}

	res := &calculatorpb.ComputeAverageResponse{
		Result: num / float64(msgCount),
	}

	stream.SendAndClose(res)

	return nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	var max int32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalln("Error while trying to read from stream:", err)
		}

		fmt.Println("FindMaximum — Message received:", req.GetNumber())

		num := req.GetNumber()

		if num > max {
			max = num
		}

		stream.Send(&calculatorpb.FindMaximumResponse{
			Result: max,
		})

		time.Sleep(time.Second)
	}

	return nil
}

func main() {
	mux, err := net.Listen("tcp", ":6543")
	if err != nil {
		log.Fatal("failed to listen", err)
	}

	s := grpc.NewServer()

	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	err = s.Serve(mux)
	if err != nil {
		log.Fatal("failed to listen", err)
	}
}
