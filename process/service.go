package main

import (
	"chandy-lamport/remoteProcedures"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
)

type PeerFunctionServer struct {
	remoteProcedures.UnimplementedPeerFunctionServer
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
	remoteProcedures.RegisterPeerFunctionServer(peerServer, peerService)

	return lis, serviceAddr, peerServer
}

func (s PeerFunctionServer) NewPeerAdded(_ context.Context, peer *remoteProcedures.Peer) (*remoteProcedures.Empty, error) {
	peerList.PeerList = append(peerList.PeerList, peer)

	return nil, nil
}
