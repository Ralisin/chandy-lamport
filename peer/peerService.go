package main

import (
	"chandy-lamport/snapshotService"
	"errors"
	"fmt"
	"golang.org/x/net/context"
)

type PeerFunctionServer struct {
	snapshotService.UnimplementedPeerFunctionServer
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
