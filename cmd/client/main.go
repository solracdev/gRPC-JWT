package main

import (
	"context"
	"github.com/carlos/grpc/proto"
	"github.com/carlos/grpc/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"log"
)

const port = ":8050"

var names = proto.NameList{
	Name: []string{"john", "Doe", "Max", "Payne"},
}

func main() {
	conn, err := grpc.Dial("localhost"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can't connect %v", err)
	}
	defer func() { _ = conn.Close() }()

	a := rpc.NewAuthClient(conn)
	token, err := a.CallGenerateToken()
	if err != nil {
		log.Fatal(err)
	}

	ctx := metadata.AppendToOutgoingContext(context.Background(), "Authorization", token)
	c := rpc.NewGreetClient(conn)
	c.CallSayHello(ctx)
	//c.CallSayHelloServerStream(ctx, &names)
	//c.CallSayHelloClientStream(ctx, &names)
	//c.CallSayHelloBidirectionalStream(ctx, &names)
}
