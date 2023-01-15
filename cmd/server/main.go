package main

import (
	"github.com/carlos/grpc/jwt"
	"github.com/carlos/grpc/proto"
	"github.com/carlos/grpc/rpc"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const port = ":8050"

var (
	secretKey     = []byte("my_super_secret_key_that_should_be_anywhere_else")
	tokenDuration = 15 * time.Minute
)

func main() {
	list, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to start the server %v", err)
	}

	// Auth server
	jwtManager := jwt.NewManager(secretKey, tokenDuration)
	authServer := rpc.NewAuthServer(*jwtManager)

	// Interceptor
	interceptor := rpc.NewInterceptor(*jwtManager)

	// Grpc Server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.Unary),
		grpc.StreamInterceptor(interceptor.Stream),
	)

	// Register
	proto.RegisterGreetServiceServer(grpcServer, &rpc.HelloServer{})
	proto.RegisterAuthServiceServer(grpcServer, authServer)
	if err := grpcServer.Serve(list); err != nil {
		log.Fatalf("failed to start the grpc server %v", err)
	}
}
