package wsb

import (
	"code.google.com/p/go.net/websocket"
	"crypto/rand"
	"encoding/base64"
	"log"
)

const STRLEN = 10

func uuid() string {
	b := make([]byte, STRLEN)
	rand.Read(b)
	en := base64.StdEncoding // or URLEncoding
	d := make([]byte, en.EncodedLen(len(b)))
	en.Encode(d, b)
	return base64.URLEncoding.EncodeToString(d)
}

type Client struct {
	ws          *websocket.Conn
	server      *Server
	ch          chan *[]byte
	chHandShake chan []byte
	done        chan bool
	id          string
}

const channelBufSize = 1024

func NewClient(ws *websocket.Conn, server *Server) *Client {

	if ws == nil {
		panic("ws cannot be nil")
	} else if server == nil {
		panic("server cannot be nil")
	}

	ch := make(chan *[]byte, channelBufSize)
	chHandShake := make(chan []byte)
	done := make(chan bool)

	return &Client{ws, server, ch, chHandShake, done, uuid()}
}

func (self *Client) Conn() *websocket.Conn {
	return self.ws
}

func (self *Client) Write() chan<- *[]byte {
	return (chan<- *[]byte)(self.ch)
}

func (self *Client) Done() chan<- bool {
	return (chan<- bool)(self.done)
}

func (self *Client) Listen() {
	go self.listenRead()
	self.listenWrite()
}

func (self *Client) listenWrite() {
	log.Println("Listening write to client ", self.id)

	for {
		select {
		case handShake := <-self.chHandShake:
			log.Println("Sending handshake to client: ", handShake)
			websocket.Message.Send(self.ws, handShake)

		case msg := <-self.ch:
			//log.Println("Sending message to client ", self.id, " size=", len(*msg))
			websocket.Message.Send(self.ws, *msg)

		case <-self.done:
			self.server.RemoveClient() <- self
			self.done <- true // for listenRead method
			return
		}
	}
}

func (self *Client) listenRead() {
	log.Println("Listening read from client ", self.id)
	for {
		select {
		case <-self.done:
			self.server.RemoveClient() <- self
			self.done <- true // for listenWrite method
			return

		default:
			var msg []byte
			err := websocket.Message.Receive(self.ws, &msg)
			if err != nil {
				self.done <- true
			}
		}
	}
}
