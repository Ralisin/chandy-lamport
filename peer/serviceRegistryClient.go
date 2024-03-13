package main

import (
	"chandy-lamport/snapshotService"
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
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection: %s", err)
		}
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
