# gRPC and JWT Authentication

This repository demonstrates how to implement JSON Web Token (JWT) authentication in a gRPC service.

## Prerequisites
- Go
- protobuf
- protoc
- grpc
- jwt-go
- make

## Getting Started

1. Clone the repository
```bash
git clone https://github.com/xSolrac87/gRPC-JWT
```

2. Install the dependencies
```bash
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
go get -u github.com/dgrijalva/jwt-go
```
3. Generate the protobuf files
```bash
protoc --go_out=. --go-grpc_out=. proto/greet.proto
protoc --go_out=. --go-grpc_out=. proto/auth.proto
```

4. Run the server
```bash
make build-server
make run-server
```

5. Run the client
```bash
make build-client
make run-client
```

## How it works
The gRPC service uses JWT for authentication. When a client makes a request to the server, it must include a JWT in the Authorization header. The server will then validate the JWT and only process the request if it is valid.

The server uses a simple hardcoded JWT for demonstration purposes. In a production environment, you would want to use a more secure method for generating and storing JWTs, such as using a JWT library and storing them in a secure database.

The client will add the JWT to the headers before making the request to the server, this JWT is hardcoded on the client file for demonstration purposes, in a production environment the JWT should be obtained from a secure authentication system.

## Conclusion
This repository provides a simple example of how to implement JWT authentication in a gRPC service. You can use the code and concepts demonstrated here as a starting point for building your own secure gRPC service.

Please note that this example is for demonstration purposes only and should not be used in production without proper security measures.