package main

import (
	"chandy-lamport/snapshotService"
	"fmt"
	"log"
	"math/rand"
	"time"
)

var currPeer snapshotService.Peer
var peerList snapshotService.PeerList

func main() {
	/* Initialize Peer Service server */
	lis, serviceAddr, peerServer := initPeerServiceServer()
	log.Printf("Peer service server initialized: %s", lis.Addr())

	/* Register peer service on Service Registry and get current peerList in Service Registry */
	registerPeerServiceOnServiceRegistry(serviceAddr.String())
	log.Printf("Peer registered on Service Registry")

	/* Launch peerJob on different thread so main-one could stay on Serve state and wait for messages */
	go peerJob()

	/* Listen for Remote Procedure Call */
	if err := peerServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve process over port []: %s", err)
	}
}

func peerJob() {
	for {
		if len(peerList.PeerList) == 0 {
			continue
		}

		/* Send messages to all linked peers */
		// emptyPeerList is used to manage contacted peer failures
		emptyPeerList := snapshotService.PeerList{PeerList: nil}
		// Example: send a message to every peer into peerList
		for _, peerToCall := range peerList.PeerList {
			messageStr := fmt.Sprintf("caller peer: %d -> called peer: %d", currPeer.Id, peerToCall.Id)

			// Send a message to
			err := sendMessageToPeer(peerToCall, messageStr)
			if err != nil {
				continue
			}

			// Append to emptyPeerList only working peers
			emptyPeerList.PeerList = append(emptyPeerList.PeerList, peerToCall)
		}
		// Update peerList.PeerList with only working peers
		peerList.PeerList = emptyPeerList.PeerList

		random := rand.New(rand.NewSource(time.Now().UnixNano()))
		randTime := random.Intn(10) + 1
		log.Print(randTime)
		time.Sleep(time.Duration(randTime) * time.Second)
	}
}
