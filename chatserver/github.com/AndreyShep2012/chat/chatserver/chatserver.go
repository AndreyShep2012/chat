package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

func (mh *messageHandle) last() messageUnit {
	r := messageUnit{}
	if len(mh.MQue) > 0 {
		mh.mu.Lock()
		r = mh.MQue[len(mh.MQue)-1]
		mh.mu.Unlock()
	}
	return r
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

//ChatService -
func (is *ChatServer) ChatService(csi Services_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	// receive messages
	go receiveFromStream(csi, clientUniqueCode, errch)

	// send messages
	go sendToStream(csi, clientUniqueCode, errch)

	return <-errch
}

func receiveFromStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	for {
		msg, err := csi_.Recv()
		if err != nil {
			log.Println("Err in receive", err)
			errch_ <- err
		} else {
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        msg.Name,
				MessageBody:       msg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode_,
			})

			messageHandleObject.mu.Unlock()

			log.Printf("%v", messageHandleObject.last())

			log.Println(messageHandleObject.last())
		}
	}
}

func sendToStream(csi_ Services_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	for {
		for {
			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			msg := messageHandleObject.MQue[0]
			senderUniqueCode := msg.ClientUniqueCode
			senderName4Client := msg.ClientName
			message4Client := msg.MessageBody

			messageHandleObject.mu.Unlock()

			//send message to designed client (do not send message to the same client)
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&FromServer{
					Name: senderName4Client,
					Body: message4Client,
				})

				if err != nil {
					errch_ <- err
				}

				messageHandleObject.mu.Lock()

				if len(messageHandleObject.MQue) > 1 {
					messageHandleObject.MQue = messageHandleObject.MQue[1:]
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}

				messageHandleObject.mu.Unlock()
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}
