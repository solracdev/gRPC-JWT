package rpc

import (
	"context"
	"fmt"
	"github.com/carlos/grpc/jwt"
	"github.com/carlos/grpc/proto"
)

type AuthServer struct {
	proto.AuthServiceServer
	jwtManager jwt.Manager
}

func NewAuthServer(manager jwt.Manager) *AuthServer {
	return &AuthServer{
		jwtManager: manager,
	}
}

func (a *AuthServer) GenerateToken(ctx context.Context, req *proto.RequestToken) (*proto.TokenResponse, error) {
	fmt.Println("got unary request from client", req.String())
	token, err := a.jwtManager.Generate()
	if err != nil {
		return nil, err
	}

	fmt.Println("token generated")
	return &proto.TokenResponse{AccessToken: token}, nil
}
