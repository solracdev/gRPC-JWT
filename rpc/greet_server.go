package rpc

import (
	"context"
	"fmt"
	"github.com/carlos/grpc/proto"
	"io"
	"log"
	"time"
)

type HelloServer struct {
	proto.GreetServiceServer
}

// SayHello unary gRPC method
func (h *HelloServer) SayHello(ctx context.Context, req *proto.NoParam) (*proto.HelloResponse, error) {
	fmt.Println("got unary request from client", req.String())
	return &proto.HelloResponse{Message: "Hello"}, nil
}

// SayHelloServerStreaming streaming server-side
func (h *HelloServer) SayHelloServerStreaming(req *proto.NameList, stream proto.GreetService_SayHelloServerStreamingServer) error {
	log.Printf("got request with names: %v", req.Name)
	ticker := time.NewTicker(2 * time.Second)
	for _, name := range req.Name {
		res := &proto.HelloResponse{
			Message: "Hello " + name,
		}
		if err := stream.Send(res); err != nil {
			return err
		}
		<-ticker.C
	}
	return nil
}

// SayHelloClientStreaming client sending streaming
func (h *HelloServer) SayHelloClientStreaming(stream proto.GreetService_SayHelloClientStreamingServer) error {
	// TODO dont like this part, refactor message initialization with make, and change SendAndClose to do it out of for
	var message []string
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.MessageList{Message: message})
		}

		if err != nil {
			return err
		}

		log.Printf("got request with name: %v", req.Name)
		message = append(message, "Hello "+req.Name)
	}
}

func (h *HelloServer) SayHelloBidirectionalStreaming(stream proto.GreetService_SayHelloBidirectionalStreamingServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		log.Println("got request with name: ", req.Name)
		res := &proto.HelloResponse{
			Message: "Hello " + req.Name,
		}
		if err = stream.Send(res); err != nil {
			return err
		}
	}
}
