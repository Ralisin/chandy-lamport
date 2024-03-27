package main

import (
	"chandy-lamport/snapshotService"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net"
)

type PeerFunctionServer struct {
	snapshotService.UnimplementedPeerFunctionServer
}

var serviceRegistry snapshotService.ServiceRegistryClient

// initPeerServiceServer serves to initialize peer server
func initPeerServiceServer() (net.Listener, net.Addr, *grpc.Server) {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:", peerAddr))
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

// registerPeerServiceOnServiceRegistry register peer service on Service Registry
func registerPeerServiceOnServiceRegistry(serviceAddr string) {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", serviceRegistryAddr, serviceRegistryPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	serviceRegistry = snapshotService.NewServiceRegistryClient(conn)

	peerStruct := snapshotService.Peer{Addr: serviceAddr}

	// Register process to Service Registry
	registerPeerResponse, err := serviceRegistry.RegisterPeer(context.Background(), &peerStruct)
	if err != nil {
		log.Fatalf("Error when calling RegisterPeer: %s", err)
	}

	// Register peerStruct info into myPeer
	myPeer.Id, myPeer.Addr = registerPeerResponse.Peer.Id, registerPeerResponse.Peer.Addr

	// Register peerStruct list into peerList global variable
	peerList.PeerList = registerPeerResponse.PeerList
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
