package main

import (
	"chandy-lamport/snapshotService"
	"errors"
	"fmt"
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

	log.Printf("lis.Addr: %s", lis.Addr().String())

	// Get address of peer server service
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

	return nil, nil
}

// SendMessage RPC call used to send message from one peer to another one
func (s PeerFunctionServer) SendMessage(_ context.Context, message *snapshotService.Message) (*snapshotService.Empty, error) {
	if message.MType != snapshotService.MessageType_MESSAGE && message.MType != snapshotService.MessageType_MARKER {
		return nil, errors.New(
			fmt.Sprintf(
				"message type not valid. Received %s, expected snapshotService.MessageType_MESSAGE"+
					" or snapshotService.MessageType_MARKER",
				message.MType,
			))
	}

	messageListMutex.Lock()
	defer messageListMutex.Unlock()
	messageList = append(
		messageList,
		snapshotService.Message{
			MType:   message.MType,
			Message: message.Message,
			Peer:    message.Peer,
		})

	/*
		if message.MType == snapshotService.MessageType_MARKER {
			markerListMutex.Lock()
			defer markerListMutex.Unlock()

			markerList = append(
				markerList,
				snapshotService.Message{
					MType:   message.MType,
					Message: message.Message,
					Peer:    message.Peer,
				})
		}
	*/

	return nil, nil
}
