package main

import (
	"chandy-lamport/snapshotService"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

const serviceRegistryAddr = "localhost"
const serviceRegistryPort = "3030"

var serviceRegistry snapshotService.ServiceRegistryClient

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
		log.Fatalf("Error when calling RegisterProcess: %s", err)
	}

	// Register peerStruct info into currPeer
	currPeer.Id, currPeer.Addr = registerPeerResponse.Peer.Id, registerPeerResponse.Peer.Addr

	// Register peerStruct list into peerList global variable
	peerList.PeerList = registerPeerResponse.PeerList
}

// sendMessageToPeer send a string message to the peer "peerToCall" using RPC
func sendMessageToPeer(peerToCall *snapshotService.Peer, messageStr string) error {
	conn, err := grpc.Dial(peerToCall.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New(fmt.Sprintf("Did not connect: %s", err))
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	message := &snapshotService.Message{
		MType:   snapshotService.MessageType_MESSAGE,
		Message: messageStr,
	}

	peerService := snapshotService.NewPeerFunctionClient(conn)

	// Send message to peers
	_, err = peerService.SendMessage(context.Background(), message)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when calling NewPeerAdded on peer: %s, err: %s", peerToCall.Addr, err))
	}

	return nil
}
