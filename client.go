package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"grpc-stream-example/pb"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var client pb.GreeterClient

func main() {

	conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client = pb.NewGreeterClient(conn)

	srv := grpc.NewServer()
	s := new(ClientServer)
	pb.RegisterClientServer(srv, s)

	reflection.Register(srv)

	l, err := net.Listen("tcp", ":8001")
	if err != nil {
		log.Fatalf("could not listen to %s: %v", "8001", err)
	}

	fmt.Println("server is running on port 8001")

	err = srv.Serve(l)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}

}

type ClientServer struct {
	pb.UnimplementedClientServer
}

func (s *ClientServer) BidirectionalStream(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	out, err := client.BidirectionalStream(context.Background())
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		for c := 0; c < 10; {
			select {
			case <-ticker.C:
				r.Name = fmt.Sprintf("%s %d", r.Name, c)
				if c == 5 {
					r.Action = "BREAK"
				}

				err = out.Send(r)
				if err != nil {
					fmt.Println("error send:", err)
					return
				}

				c++
			}
		}
	}()

	go func() {
		defer wg.Done()

		for {
			resp, err := out.Recv()
			if err != nil {
				fmt.Println("error recv:", err)
				return
			}

			fmt.Println("got msg BidirectionalStream:", resp.Message)
			if resp.Message == "BREAK" {
				return
			}
		}
	}()

	wg.Wait()

	return &pb.Response{
		Message: "done",
	}, nil
}

func (s *ClientServer) ServerStream(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	out, err := client.ServerStream(ctx, r)
	if err != nil {
		return nil, err
	}

	for {
		resp, err := out.Recv()
		if err != nil {
			fmt.Println("error recv:", err)
			return nil, err
		}

		fmt.Println("got msg ServerStream:", resp.Message)
		if strings.Contains(resp.Message, "Bye") {
			return resp, nil
		}
	}
}

func (s *ClientServer) ClientStream(ctx context.Context, r *pb.Request) (*pb.Response, error) {
	out, err := client.ClientStream(ctx)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(1000 * time.Millisecond)
	defer ticker.Stop()

	for c := 0; c < 10; {
		select {
		case <-ticker.C:

			r.Name = fmt.Sprintf("%s %d", r.Name, c)
			if c == 5 {
				r.Action = "BREAK"
			}

			err = out.Send(r)
			if err != nil {
				fmt.Println("error send:", err)
			}

			c++
		}
	}

	return out.CloseAndRecv()
}
