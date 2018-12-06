package main

import (
	"fmt"
	"net"
	"net/http"

	pb "github.com/lsytj0413/ena/demo/gogrpc/proto"
	"golang.org/x/net/context"
	"golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/metadata"
)

const (
	Address = "127.0.0.1:50052"
)

type helloService struct {
}

var HelloService = &helloService{}

func (h *helloService) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc.Errorf(codes.Unauthenticated, "less token")
	}

	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return nil, grpc.Errorf(codes.Unauthenticated, "token invalid")
	}

	resp := new(pb.HelloReply)
	// resp.Message = "Hello " + in.Name + "."
	resp.Message = fmt.Sprintf("Hello %s.\n Token info: appid=%s,appkey=%s", in.Name, appid, appkey)

	return resp, nil
}

func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return grpc.Errorf(codes.Unauthenticated, "less token")
	}

	var (
		appid  string
		appkey string
	)
	if val, ok := md["appid"]; ok {
		appid = val[0]
	}
	if val, ok := md["appkey"]; ok {
		appkey = val[0]
	}

	if appid != "101010" || appkey != "i am key" {
		return grpc.Errorf(codes.Unauthenticated, "token invalid")
	}

	return nil
}

func main() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption

	// For tls
	creds, err := credentials.NewServerTLSFromFile("../../keys/server.pem", "../../keys/server.key")
	if err != nil {
		grpclog.Fatalf("Failed to generate credentials %v", err)
	}
	opts = append(opts, grpc.Creds(creds))

	var interceptor grpc.UnaryServerInterceptor
	interceptor = func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		err = auth(ctx)
		if err != nil {
			return
		}

		return handler(ctx, req)
	}
	opts = append(opts, grpc.UnaryInterceptor(interceptor))
	s := grpc.NewServer(opts...)

	pb.RegisterHelloServer(s, HelloService)

	go func() {
		trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
			return true, true
		}
		go http.ListenAndServe(":50051", nil)
		fmt.Println("Trace on 50051")
	}()

	fmt.Println("Listen on: " + Address)
	s.Serve(listen)
}
