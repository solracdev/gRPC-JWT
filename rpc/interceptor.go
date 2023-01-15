package rpc

import (
	"context"
	"errors"
	"github.com/carlos/grpc/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

var (
	errMetadata = errors.New("can't get metadata from context")
	errNoToken  = errors.New("no token provided")
)

type Interceptor struct {
	authManager jwt.Manager
}

func NewInterceptor(manager jwt.Manager) *Interceptor {
	return &Interceptor{
		authManager: manager,
	}
}

func (i *Interceptor) Unary(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("--> unary interceptor: ", info.FullMethod)
	// TODO this is just for demo purpose
	if !strings.Contains(info.FullMethod, "GenerateToken") {
		err := i.validateToken(ctx)
		if err != nil {
			return nil, err
		}
	}
	return handler(ctx, req)
}

func (i *Interceptor) Stream(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Println("--> stream interceptor: ", info.FullMethod)
	err := i.validateToken(stream.Context())
	if err != nil {
		return err
	}
	return handler(srv, stream)
}

func (i *Interceptor) validateToken(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errMetadata
	}

	var token string
	if exist := md.Get("Authorization"); exist != nil {
		token = exist[0]
	}

	if token == "" {
		return errNoToken
	}

	return i.authManager.Verify(token)
}
