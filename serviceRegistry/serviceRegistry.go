package main

import (
	"chandy-lamport/snapshotService"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

const port = "3030"

var peerList snapshotService.PeerList

// Start Service Registry and handle peers registration
func main() {
	/* Initialize Service Registry service */
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Failed to start the Service Registry: %s", err)
	}

	// Create a gRPC server with no service registered
	serverRegister := grpc.NewServer()

	// Register ServiceRegistry as a service
	serviceRegistry := ServiceRegistryServer{}
	snapshotService.RegisterServiceRegistryServer(serverRegister, serviceRegistry)

	log.Printf("Service Registry started")

	// Listen for Remote Procedure Call
	if err := serverRegister.Serve(lis); err != nil {
		log.Fatalf("Failed to serve process over port []: %s", err)
	}
}
