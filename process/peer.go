package main

import (
	"chandy-lamport/remoteProcedures"
	"log"
)

var currPeer remoteProcedures.Peer
var peerList remoteProcedures.PeerList

func main() {
	/* Initialize Peer Service server */
	lis, serviceAddr, peerServer := initPeerServiceServer()
	log.Printf("Peer service server initialized")

	/* Register peer service on Service Registry */
	registerPeerServiceOnServiceRegistry(serviceAddr.String())
	log.Printf("Peer registered on Service Registry")

	// Listen for Remote Procedure Call
	if err := peerServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve process over port []: %s", err)
	}
}
