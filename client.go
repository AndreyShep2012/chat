package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AndreyShep2012/chat/chatserver/github.com/AndreyShep2012/chat/chatserver"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Enter server addr")

	reader := bufio.NewReader(os.Stdin)
	serverID, err := reader.ReadString('\n')
	if err != nil {
		fmt.Errorf("can't read server addr %v", err)
	}

	serverID = strings.Trim(serverID, "\r\n")
	log.Println("connecting to", serverID)

	conn, err := grpc.Dial(serverID, grpc.WithInsecure())
	if err != nil {
		fmt.Errorf("can't connect to server %v", err)
	}

	defer conn.Close()

	client := chatserver.NewServicesClient(conn)

	stream, err := client.ChatService(context.Background())
	if err != nil {
		fmt.Errorf("can't create chat service %v", err)
	}

	ch := clientHandle{
		stream: stream,
	}
	ch.clientConfig()
	go ch.sendMessage()
	go ch.receiveMessage()

	blocker := make(chan bool)
	<-blocker
}

type clientHandle struct {
	stream     chatserver.Services_ChatServiceClient
	clientName string
}

func (ch *clientHandle) clientConfig() {
	fmt.Println("Enter name")

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Errorf("can't read client name %v", err)
	}

	ch.clientName = strings.Trim(name, "\r\n")
}

func (ch *clientHandle) sendMessage() {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, err := reader.ReadString('\n')
		if err != nil {
			fmt.Errorf("can't read client message %v", err)
		}

		message = strings.Trim(message, "\r\n")
		err = ch.stream.Send(&chatserver.FromClient{
			Name: ch.clientName,
			Body: message,
		})

		if err != nil {
			fmt.Errorf("can't send client message %v", err)
		}
	}
}

func (ch *clientHandle) receiveMessage() {
	for {
		msg, err := ch.stream.Recv()
		if err != nil {
			fmt.Errorf("can't recv message %v", err)
		}

		fmt.Printf("%s : %s\n", msg.Name, msg.Body)
	}
}
