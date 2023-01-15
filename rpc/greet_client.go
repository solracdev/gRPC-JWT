package rpc

import (
	"context"
	"fmt"
	"github.com/carlos/grpc/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

type GreetClient struct {
	gc proto.GreetServiceClient
}

func NewGreetClient(conn grpc.ClientConnInterface) *GreetClient {
	return &GreetClient{gc: proto.NewGreetServiceClient(conn)}
}

// CallSayHello unary function call
func (c *GreetClient) CallSayHello(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	resp, err := c.gc.SayHello(ctx, &proto.NoParam{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("%s", resp.Message)
}

// CallSayHelloServerStream server-side streaming function
func (c *GreetClient) CallSayHelloServerStream(ctx context.Context, names *proto.NameList) {
	log.Printf("streaming has strarted")
	stream, err := c.gc.SayHelloServerStreaming(ctx, names)
	if err != nil {
		log.Fatalf("could not send names %v", err)
	}

	for {
		message, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while streaming %v", err)
		}
		log.Println(message)
	}

	fmt.Println("stream finished")
}

// CallSayHelloClientStream client sending streaming to server
func (c *GreetClient) CallSayHelloClientStream(ctx context.Context, names *proto.NameList) {
	log.Println("client streaming started")
	// generate stream
	stream, err := c.gc.SayHelloClientStreaming(ctx)
	if err != nil {
		log.Fatalf("could not create stream %v", err)
	}

	ticker := time.NewTicker(2 * time.Second)
	for _, name := range names.Name {
		req := &proto.HelloRequest{
			Name: name,
		}

		if err = stream.Send(req); err != nil {
			log.Fatalf("error while sending %v", err)
		}

		log.Println("sent request with name: ", name)
		<-ticker.C
	}
	log.Println("client complete sending stream")

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while recieving request from server %v", err)
	}

	log.Println(resp)
}

// CallSayHelloBidirectionalStream bidirectional streaming client <-> server
func (c *GreetClient) CallSayHelloBidirectionalStream(ctx context.Context, names *proto.NameList) {
	log.Printf("bidirectional streaming started")

	// generate stream
	stream, err := c.gc.SayHelloBidirectionalStreaming(ctx)
	if err != nil {
		log.Fatalf("could not create stream %v", err)
	}

	// chanel to block til all stream from server is consumed
	ch := make(chan struct{})

	// listener from server
	go func() {
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("error while recieving from server %v", err)
			}
			log.Println(message)
		}
		close(ch)
	}()

	// sending stream data to server
	ticker := time.NewTicker(2 * time.Second)
	for _, name := range names.Name {
		req := &proto.HelloRequest{
			Name: name,
		}

		if err = stream.Send(req); err != nil {
			log.Fatalf("error while sending request %v", err)
		}
		<-ticker.C
	}

	if err = stream.CloseSend(); err != nil {
		log.Fatalf("error while closing stream %v", err)
	}

	<-ch
	log.Println("bidirectional streaming done")
}
