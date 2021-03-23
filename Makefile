#Compile the protobuf file to generate go files
proto:
	protoc --go-grpc_out=require_unimplemented_servers=false:./chatserver/ --go_out=./chatserver/ chat.proto

build_server:
	go build -o server ./server.go

build_client:
	go build -o client ./client.go
