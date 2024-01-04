package main

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-stream-example/pb"
	"log"
	"net"
	"time"
)

func main() {
	srv := grpc.NewServer()
	s := new(GreeterServer)
	pb.RegisterGreeterServer(srv, s)

	reflection.Register(srv)

	l, err := net.Listen("tcp", ":8000")
	if err != nil {
		log.Fatalf("could not listen to %s: %v", "8000", err)
	}

	fmt.Println("server is running on port 8000")

	err = srv.Serve(l)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}
}

type GreeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterServer) BidirectionalStream(in pb.Greeter_BidirectionalStreamServer) error {
	for {
		req, err := in.Recv()
		if err != nil {
			return err
		}

		fmt.Println("got msg BidirectionalStream:", req.Name)

		err = in.Send(&pb.Response{Message: req.Name})
		if err != nil {
			return err
		}

		if req.GetAction() == "BREAK" {
			return nil
		}
	}
}

func (s *GreeterServer) ClientStream(in pb.Greeter_ClientStreamServer) error {
	req := new(pb.Request)
	var err error

	for {
		req, err = in.Recv()
		if err != nil {
			return err
		}

		fmt.Println("got msg ClientStream:", req.Name, req.Action)

		if req.GetAction() == "BREAK" {
			fmt.Println("BREAK")
			break
		}
	}

	err = in.SendAndClose(&pb.Response{Message: req.Action})
	if err != nil {
		return err
	}

	return nil
}

func (s *GreeterServer) ServerStream(in *pb.Request, out pb.Greeter_ServerStreamServer) error {

	err := out.Send(&pb.Response{Message: fmt.Sprintf("Hello %s", in.Name)})
	if err != nil {
		return err
	}

	ticker := time.NewTicker(1000 * time.Millisecond)

	for c := 0; c < 10; {
		select {
		case <-ticker.C:
			err = out.Send(&pb.Response{Message: fmt.Sprintf("meledak %d", c)})
			if err != nil {
				return err
			}

			c++
		}
	}

	for i := 0; i < 10; i++ {
		err = out.Send(&pb.Response{Message: fmt.Sprintf("%s %d", in.Name, i)})
		if err != nil {
			return err
		}
	}

	time.Sleep(2 * time.Second)

	err = out.Send(&pb.Response{Message: fmt.Sprintf("Bye %s", in.Name)})
	if err != nil {
		return err
	}

	return nil
}
