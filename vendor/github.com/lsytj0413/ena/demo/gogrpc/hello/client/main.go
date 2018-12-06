package main

import (
	"fmt"

	pb "github.com/lsytj0413/ena/demo/gogrpc/proto"
	"golang.org/x/net/context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

const (
	// Address ...
	Address = "127.0.0.1:50052"

	// OpenTLS ...
	OpenTLS = true
)

type customCredential struct{}

func (c customCredential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"appid":  "101010",
		"appkey": "i am key",
	}, nil
}

func (c customCredential) RequireTransportSecurity() bool {
	return OpenTLS
}

func main() {
	// creds, err := credentials.NewClientTLSFromFile("../../keys/server.pem", "server name")
	// if err != nil {
	// 	grpclog.Fatalf("Failed to create TLS credentials %v", err)
	// }

	// // conn, err := grpc.Dial(Address, grpc.WithInsecure())
	// conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(creds))

	// if err != nil {
	// 	grpclog.Fatalln(err)
	// }

	// defer conn.Close()

	var err error
	var opts []grpc.DialOption

	if OpenTLS {
		creds, err := credentials.NewClientTLSFromFile("../../keys/server.pem", "server name")
		if err != nil {
			grpclog.Fatalf("Failed to create TLS credentials %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	opts = append(opts, grpc.WithPerRPCCredentials(new(customCredential)))
	conn, err := grpc.Dial(Address, opts...)
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer conn.Close()

	c := pb.NewHelloClient(conn)
	reqBody := new(pb.HelloRequest)
	reqBody.Name = "gRPC"
	r, err := c.SayHello(context.Background(), reqBody)
	if err != nil {
		grpclog.Fatalln(err)
	}

	fmt.Println(r.Message)
}
