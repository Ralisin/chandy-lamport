package main

import (
	"chandy-lamport/snapshotService"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type ServiceRegistryServer struct {
	snapshotService.UnimplementedServiceRegistryServer
}

// RegisterPeer Service to register a new process in the process of available processes
func (s ServiceRegistryServer) RegisterPeer(_ context.Context, peerToRegister *snapshotService.Peer) (*snapshotService.RegisterPeerResponse, error) {
	// Create new peer to add to list with address of peer service
	newPeer := snapshotService.Peer{
		Id:   getNewId(),
		Addr: peerToRegister.Addr,
	}

	// Create response with list without newPeer
	registerPeerResponse := snapshotService.RegisterPeerResponse{
		Peer:     &newPeer,
		PeerList: peerList.PeerList,
	}

	// Update all peers list with newPeer
	for _, peerEl := range peerList.PeerList {
		callServiceNewPeerAdded(peerEl, &newPeer)
	}

	// Append newPeer to peerList
	peerList.PeerList = append(peerList.PeerList, &newPeer)

	// Print some log to show ServiceRegistry workload
	log.Println(peerList.PeerList)

	return &registerPeerResponse, nil
}

func callServiceNewPeerAdded(peerToCall *snapshotService.Peer, newPeer *snapshotService.Peer) {
	conn, err := grpc.Dial(peerToCall.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Error closing connection: %s", err)
		}
	}(conn)

	peerService := snapshotService.NewPeerFunctionClient(conn)

	// Register process to Service Registry
	_, err = peerService.NewPeerAdded(context.Background(), newPeer)
	if err != nil {
		log.Fatalf("Error when calling NewPeerAdded on peer: %s, err: %s", peerToCall.Addr, err)
	}
}

func getNewId() int32 {
	// Handle the case when the listPeers is empty
	if len(peerList.PeerList) == 0 {
		return 1
	}

	return peerList.PeerList[len(peerList.PeerList)-1].Id + 1
}
