package rpc

import (
	"context"
	"github.com/carlos/grpc/proto"
	"google.golang.org/grpc"
	"time"
)

type AuthClient struct {
	ac proto.AuthServiceClient
}

func NewAuthClient(conn grpc.ClientConnInterface) *AuthClient {
	return &AuthClient{ac: proto.NewAuthServiceClient(conn)}
}

func (a *AuthClient) CallGenerateToken() (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Second)
	defer cancel()

	resp, err := a.ac.GenerateToken(ctx, &proto.RequestToken{})
	if err != nil {
		return "", err
	}
	return resp.AccessToken, nil
}
