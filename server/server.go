package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	sgRPC "github.com/LogeshVel/simple_gRPC/simple/testgrpc"
	"google.golang.org/grpc"
)

const socket string = "localhost:8090"

type Server struct {
	sgRPC.SimpleServiceServer
}

func main() {
	lisn, err := net.Listen("tcp", socket)
	if err != nil {
		log.Fatalln("Errored while Listen to : ", socket, err)
	}
	log.Println("Listening at ", socket)
	s := grpc.NewServer()
	sgRPC.RegisterSimpleServiceServer(s, &Server{}) // registering our grpc server with our grpc service.
	err = s.Serve(lisn)
	if err != nil {
		log.Fatalln("Errored while Serving : ", socket, err)
	}
}

// Implementing all these rpc methods defined in the Proto file.
// These Methods Signatures are found in the compiled output of the Proto file.

// type SimpleServiceClient interface {
// 	RPCRequest(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (*SimpleResponse, error)
// 	ServerStreaming(ctx context.Context, in *SimpleRequest, opts ...grpc.CallOption) (SimpleService_ServerStreamingClient, error)
// 	ClientStreaming(ctx context.Context, opts ...grpc.CallOption) (SimpleService_ClientStreamingClient, error)
// 	StreamingBiDirectional(ctx context.Context, opts ...grpc.CallOption) (SimpleService_StreamingBiDirectionalClient, error)
// }

func (s *Server) RPCRequest(ctx context.Context, req *sgRPC.SimpleRequest) (*sgRPC.SimpleResponse, error) {
	log.Println("Unary request")
	log.Printf("Request - %v\n", req)
	r := fmt.Sprintf("Here is your response for the request msg %v", req.RequestNeed)
	response := &sgRPC.SimpleResponse{Response: r}
	log.Printf("Response - %v\n", response)
	return response, nil
}

func (s *Server) ClientStreaming(stream sgRPC.SimpleService_ClientStreamingServer) error {
	log.Println("ClientStreaming RPC")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			response := &sgRPC.SimpleResponse{Response: "Here is your Client Streaming response"}
			log.Printf("Response - %v\n", response)
			stream.SendAndClose(response)
			break
		}
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("Request - %v\n", req)
	}
	return nil
}

func (s *Server) ServerStreaming(req *sgRPC.SimpleRequest, stream sgRPC.SimpleService_ServerStreamingServer) error {
	log.Println("ServerStreaming RPC")
	log.Printf("Request- %v", req)
	for i := 1; i < 10; i++ {
		res := fmt.Sprintf("Here is the response %d", i)
		log.Printf("Response - %v", res)
		stream.Send(&sgRPC.SimpleResponse{Response: res})
	}

	return nil
}

func (s *Server) StreamingBiDirectional(stream sgRPC.SimpleService_StreamingBiDirectionalServer) error {
	log.Println("StreamingBiDirectional RPC")

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Errored in stream Recv", err)
			break
		}
		log.Printf("Request - %v", msg)
		r := fmt.Sprintf("Response for your request - %v", msg.RequestNeed)
		log.Printf("Response - %v\n", r)
		stream.Send(&sgRPC.SimpleResponse{Response: r})

	}
	return nil
}
