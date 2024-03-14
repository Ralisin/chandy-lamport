package main

import (
	"chandy-lamport/snapshotService"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type PeerFunctionServer struct {
	snapshotService.UnimplementedPeerFunctionServer
}

// initPeerServiceServer serves to initialize service server, returned server must be served
func initPeerServiceServer() (net.Listener, net.Addr, *grpc.Server) {
	lis, err := net.Listen("tcp", ":")
	if err != nil {
		log.Fatalf("Failed to start the Peer service: %s", err)
	}

	// Get port of peer server service
	serviceAddr := lis.Addr()

	// Create a gRPC server with no service registered
	peerServer := grpc.NewServer()

	// Register ServiceRegistry as a service
	peerService := PeerFunctionServer{}
	snapshotService.RegisterPeerFunctionServer(peerServer, peerService)

	return lis, serviceAddr, peerServer
}

// NewPeerAdded update peerList and append the new peer to the end
func (s PeerFunctionServer) NewPeerAdded(_ context.Context, peer *snapshotService.Peer) (*snapshotService.Empty, error) {
	peerList.PeerList = append(peerList.PeerList, peer)

	// TODO: understand if peerList must be sorted

	return nil, nil
}

// SendMessage RPC call used to send message from one peer to another one
func (s PeerFunctionServer) SendMessage(_ context.Context, message *snapshotService.Message) (*snapshotService.Empty, error) {
	// TODO: implement
	log.Printf("%s", message)

	return nil, nil
}
