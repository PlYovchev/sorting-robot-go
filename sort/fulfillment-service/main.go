package main

import (
	"fmt"
	"log"
	"net"

	"github.com/plyovchev/sorting-robot-go/sort/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serverPort = "0.0.0.0:10001"
const sortingRobotPort = "0.0.0.0:10000"

func main() {
	grpcServer, lis, sortingRobotConnection := newFulfillmentServer()

	defer sortingRobotConnection.Close()

	fmt.Printf("gRPC server started. Listening on %s\n", serverPort)
	grpcServer.Serve(lis)
}

func newFulfillmentServer() (*grpc.Server, net.Listener, *grpc.ClientConn) {
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	conn, err := grpc.Dial(sortingRobotPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}

	sortingRobotClient := gen.NewSortingRobotClient(conn)
	grpcServer := grpc.NewServer()
	gen.RegisterFulfillmentServer(grpcServer, newFulfillmentService(sortingRobotClient))
	reflection.Register(grpcServer)

	return grpcServer, lis, conn
}
