package main

import (
	"log"
	"net"
	"os"

	"github.com/AndreyShep2012/chat/chatserver/github.com/AndreyShep2012/chat/chatserver"
	"google.golang.org/grpc"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Could not listen port @ %v :: %v", port, err)
	}

	log.Println("Listening @ : " + port)
	gprcserver := grpc.NewServer()

	cs := chatserver.ChatServer{}
	chatserver.RegisterServicesServer(gprcserver, &cs)

	if err := gprcserver.Serve(listen); err != nil {
		log.Fatalf("Could not start grpc server %v", err)
	}
}
