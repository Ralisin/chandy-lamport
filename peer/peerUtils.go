package main

import (
	"chandy-lamport/snapshotService"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"
)

// Config Configuration struct to represent the JSON data
type Config struct {
	Localhost LocalhostConfig `json:"localhost"`
	Docker    DockerConfig    `json:"docker"`
}

type LocalhostConfig struct {
	ServiceRegistryAddr string `json:"serviceRegistryAddr"`
	ServiceRegistryPort string `json:"serviceRegistryPort"`
	PeerAddr            string `json:"peerAddr"`
}

type DockerConfig struct {
	ServiceRegistryAddr string `json:"serviceRegistryAddr"`
	ServiceRegistryPort string `json:"serviceRegistryPort"`
	PeerAddr            string `json:"peerAddr"`
}

var serviceRegistryAddr string
var serviceRegistryPort string
var peerAddr string

var (
	messageList      []snapshotService.Message
	messageListMutex sync.Mutex

	snapshot      bool
	snapshotMutex sync.Mutex // Needed to sync snapshot through threads
)

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
		log.Fatalf("Error when calling RegisterPeer: %s", err)
	}

	// Register peerStruct info into myPeer
	myPeer.Id, myPeer.Addr = registerPeerResponse.Peer.Id, registerPeerResponse.Peer.Addr

	// Register peerStruct list into peerList global variable
	peerList.PeerList = registerPeerResponse.PeerList
}

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

// sendMessageToPeer send a string message to the peer "peerToCall" using RPC
func sendMessageToPeer(peerToCall *snapshotService.Peer, message *snapshotService.Message) error {
	conn, err := grpc.Dial(peerToCall.Addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return errors.New(fmt.Sprintf("Did not connect: %s", err))
	}
	defer func(conn *grpc.ClientConn) {
		_ = conn.Close()
	}(conn)

	peerService := snapshotService.NewPeerFunctionClient(conn)

	// Send message to peers
	_, err = peerService.SendMessage(context.Background(), message)
	if err != nil {
		return errors.New(fmt.Sprintf("Error when calling NewPeerAdded on peer: %s, err: %s", peerToCall.Addr, err))
	}

	return nil
}

// peerSendMessagesJob randomly send message to other peers, changing message type
func peerSendMessagesJob() {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))

	count := 0

	for {
		// Wait for a new peer to enter into the system
		if len(peerList.PeerList) == 0 {
			continue
		}

		// Generate message to send
		message := &snapshotService.Message{}

		count++
		message.Message = fmt.Sprintf("Agabubu, id sender: %d, count: %d", myPeer.Id, count)
		message.Peer = &myPeer

		if random.Intn(10) < 9 || snapshot {
			message.MType = snapshotService.MessageType_MESSAGE
		} else {
			message.MType = snapshotService.MessageType_MARKER

			snapshotMutex.Lock()
			// TODO make peer snapshot via function
			snapshot = !snapshot
			snapshotMutex.Unlock()
		}

		/* Send messages to all linked peers */
		// emptyPeerList is used to manage contacted peer failures
		emptyPeerList := snapshotService.PeerList{PeerList: nil}
		// Example: send a message to every peer into peerList
		for _, peerToCall := range peerList.PeerList {
			// Send message to peerToCall
			err := sendMessageToPeer(peerToCall, message)
			if err != nil {
				continue
			}

			// Append to emptyPeerList only working peers
			emptyPeerList.PeerList = append(emptyPeerList.PeerList, peerToCall)
		}
		// Update peerList.PeerList with only working peers
		peerList.PeerList = emptyPeerList.PeerList

		randTime := random.Intn(10) + 1
		log.Print("timeSleep:", randTime)
		time.Sleep(time.Duration(randTime) * time.Second)
	}
}

func peerReceiveMessagesJob() {
	// Map to track which peer send marker message
	peerMap := make(map[int32]bool)

	// Slice of messages for every peer
	messagesPerPeer := make(map[int32][]snapshotService.Message)

	for {
		messageListMutex.Lock()
		// Skip for loop if there is not any message (messageType or markerType)
		if len(messageList) == 0 {
			messageListMutex.Unlock()
			continue
		}
		messageListMutex.Unlock()

		// Get first message from queue
		msg := snapshotService.Message{
			MType:   messageList[0].MType,
			Message: messageList[0].Message,
			Peer:    messageList[0].Peer,
		}

		log.Printf("MType: %s, %d -> %d", msg.MType, msg.Peer.Id, myPeer.Id)

		if msg.MType == snapshotService.MessageType_MESSAGE {
			snapshotMutex.Lock()
			// Check if snapshot is going on and message from peer msg.Peer.Id must be recorded
			if snapshot && !peerMap[msg.Peer.Id] {
				// Append in messagesPerPeer the message received
				messagesPerPeer[msg.Peer.Id] = append(
					messagesPerPeer[msg.Peer.Id],
					snapshotService.Message{
						MType:   msg.MType,
						Message: msg.Message,
						Peer:    msg.Peer,
					},
				)
			}
			snapshotMutex.Unlock()

			// TODO compute message
			log.Printf("Message: %s", msg.Message)
		} else if msg.MType == snapshotService.MessageType_MARKER {
			// Set peer's marker message as received
			peerMap[msg.Peer.Id] = true

			snapshotMutex.Lock()
			if !snapshot {
				// TODO make snapshot via function

				// Set snapshot value as done
				snapshot = !snapshot

				// Send marker message to all peers
				for _, peerToCall := range peerList.PeerList {
					message := &snapshotService.Message{
						MType: snapshotService.MessageType_MARKER,
						Peer:  &myPeer,
					}

					time.Sleep(time.Second)

					err := sendMessageToPeer(peerToCall, message)
					if err != nil {
						continue
					}
				}
			}
			snapshotMutex.Unlock()

			if len(peerList.PeerList) == len(peerMap) {
				// Clear peerMap
				peerMap = make(map[int32]bool)

				// TODO check messagesPerPeer content
				log.Printf("All markers received\nmessagesPerPeer: %s", messagesPerPeer)

				// TODO save snapshot in a file
				// Clear messagesPerPeer
				messagesPerPeer = make(map[int32][]snapshotService.Message)

				snapshotMutex.Lock()
				snapshot = !snapshot
				snapshotMutex.Unlock()
			}
		}

		messageListMutex.Lock()
		// Pop first message
		messageList = messageList[1:]
		messageListMutex.Unlock()
	}
}

func ReadConfig(filename string) (Config, error) {
	var config Config

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		return config, fmt.Errorf("error reading JSON file: %v", err)
	}

	// Unmarshal JSON data into Config struct
	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return config, nil
}
