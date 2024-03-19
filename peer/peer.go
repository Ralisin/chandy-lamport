package main

import (
	"chandy-lamport/snapshotService"
	"log"
)

// Variables to manage peer structures
var (
	myPeer   snapshotService.Peer
	peerList snapshotService.PeerList
)

func main() {
	// Initialize Peer Service server
	lis, serviceAddr, peerServer := initPeerServiceServer()
	log.Printf("Peer service server initialized: %s", lis.Addr())

	// Register peer service on Service Registry and get current peerList in Service Registry
	registerPeerServiceOnServiceRegistry(serviceAddr.String())
	log.Printf("Peer registered on Service Registry")

	/* Launch peerSendMessagesJob on different thread so main-one could stay on Serve state and wait for messages
	 * This thread sends messages to other linked peer in the net (currently network structure is a fully connected graph
	 */
	go peerSendMessagesJob()

	/* Launch peerSendMessagesJob on different thread
	 * This thread manage received messages
	 */
	go peerReceiveMessagesJob()

	// Listen for Remote Procedure Call
	if err := peerServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve process over port []: %s", err)
	}
}
