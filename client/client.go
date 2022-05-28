package main

import (
	"context"
	"fmt"
	"io"
	"log"

	sgRPC "github.com/LogeshVel/simple_gRPC/simple/testgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const socket string = "127.0.0.1:8081"

func main() {
	// grpc uses HTTP 2 which is by default uses SSL
	conn, err := grpc.Dial(socket, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Could not connect to : ", socket)
	}
	log.Println("Connected to ", socket)
	defer conn.Close()
	client := sgRPC.NewSimpleServiceClient(conn)
	// unary request
	makeUnaryRequest(client)

	// Client streaming
	makeClientStreaming(client)

	//Server Streaming
	makeServerStreaming(client)

	// Bi-Directional
	makeBidirectional(client)
}

func makeUnaryRequest(c sgRPC.SimpleServiceClient) {
	log.Println("Making Unary Request")
	req := &sgRPC.SimpleRequest{RequestNeed: "To test!"}
	log.Printf("Request - %v\n", req)
	res, err := c.RPCRequest(context.Background(), req)
	handleAndFatalError(err)
	log.Printf("Response - %v\n", res)
}

func makeClientStreaming(c sgRPC.SimpleServiceClient) {
	log.Println("Client Streaming")
	stream, err := c.ClientStreaming(context.Background())
	handleAndFatalError(err)

	for i := 1; i < 10; i++ {
		req := fmt.Sprintf("Request number : %d", i)
		log.Printf("Request - %v\n", req)
		stream.Send(&sgRPC.SimpleRequest{RequestNeed: req})
	}
	response, err := stream.CloseAndRecv()
	handleAndFatalError(err)

	log.Printf("Response - %v\n", response)
}

func makeServerStreaming(c sgRPC.SimpleServiceClient) {
	log.Println("Server Streaming")
	req := &sgRPC.SimpleRequest{RequestNeed: "Need stream response"}
	log.Printf("Request - %v\n", req)
	serverStream, err := c.ServerStreaming(context.Background(), req)
	handleAndFatalError(err)

	for {
		response, err := serverStream.Recv()
		if err == io.EOF {
			break
		}
		handleAndFatalError(err)

		log.Printf("Response - %v\n", response)
	}
}

func makeBidirectional(c sgRPC.SimpleServiceClient) {
	log.Println("Bi-Directional Streaming")
	biStream, err := c.StreamingBiDirectional(context.Background())
	handleAndFatalError(err)

	// here the communication sequence is completely depends on how the server is implemented.
	// if the server is implemetend to give response to all the response at the end or
	// one after another it all compeletely depends on the implementation
	for i := 1; i < 11; i++ {
		req := fmt.Sprintf("My request %d", i)
		log.Printf("Request - %v\n", req)
		biStream.Send(&sgRPC.SimpleRequest{RequestNeed: req})
		if i == 10 {
			biStream.CloseSend()
		}
		reply, err := biStream.Recv()

		handleAndPrintError(err)
		log.Printf("Response - %v\n", reply)
	}
}

func handleAndPrintError(e error) {
	if e != nil {
		log.Println(e)
	}
}

func handleAndFatalError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}
