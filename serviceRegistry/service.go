package main

import (
	"chandy-lamport/snapshotService"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sort"
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

	// emptyPeerList is used to manage contacted peer failures
	emptyPeerList := snapshotService.PeerList{PeerList: nil}
	// Update the list of each peer with newPeer
	for _, peerEl := range peerList.PeerList {
		// Informs a peer in peerList of the entry of a new peer
		err := callServiceNewPeerAdded(peerEl, &newPeer)
		if err != nil {
			continue
		}

		emptyPeerList.PeerList = append(emptyPeerList.PeerList, peerEl)
	}
	peerList.PeerList = emptyPeerList.PeerList

	// Create response with list without newPeer
	registerPeerResponse := snapshotService.RegisterPeerResponse{
		Peer:     &newPeer,
		PeerList: peerList.PeerList,
	}

	// Append newPeer to peerList
	peerList.PeerList = append(peerList.PeerList, &newPeer)

	// Sort peerList
	sort.Slice(peerList.PeerList, func(i, j int) bool { return peerList.PeerList[i].Id < peerList.PeerList[j].Id })

	// Print some log to show ServiceRegistry workload
	log.Print("Curr peerList: ", peerList.PeerList)

	return &registerPeerResponse, nil
}

func callServiceNewPeerAdded(peerToCall *snapshotService.Peer, newPeer *snapshotService.Peer) error {
	conn, err := grpc.Dial(peerToCall.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	peerService := snapshotService.NewPeerFunctionClient(conn)

	// Register process to Service Registry
	_, err = peerService.NewPeerAdded(context.Background(), newPeer)
	if err != nil {
		return err
	}

	return nil
}

func getNewId() int32 {
	// Handle the case when the listPeers is empty
	if len(peerList.PeerList) == 0 {
		return 1
	}

	// Sort peerList
	sort.Slice(peerList.PeerList, func(i, j int) bool { return peerList.PeerList[i].Id < peerList.PeerList[j].Id })

	var minId int32 = 1
	for _, peerEl := range peerList.PeerList {
		if minId < peerEl.Id {
			return minId
		}

		minId++
	}

	return minId
}
