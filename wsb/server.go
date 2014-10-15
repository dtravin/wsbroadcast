package wsb

import (
	"code.google.com/p/go.net/websocket"
	"log"
	"net/http"
)

// Chat server.
type Server struct {
	path              string
	clients           []*Client
	addClient         chan *Client
	removeClient      chan *Client
	sendAll           chan *[]byte
	handShakeMessage  []byte
	inputControl      chan bool
	inputWS           *websocket.Conn
	inputStreamActive bool
	inputStreamURI    string
}

// Create new chat server.
func NewServer(path string, inputStreamURI string) *Server {
	clients := make([]*Client, 0)
	addClient := make(chan *Client)
	removeClient := make(chan *Client)
	sendAll := make(chan *[]byte, 1000)
	inputControl := make(chan bool, 0)
	handShakeMessage := make([]byte, 0)
	return &Server{path, clients, addClient, removeClient, sendAll, handShakeMessage, inputControl, nil, false, inputStreamURI}
}

func (self *Server) InputControl() chan<- bool {
	return (chan<- bool)(self.inputControl)
}

func (self *Server) AddClient() chan<- *Client {
	return (chan<- *Client)(self.addClient)
}

func (self *Server) RemoveClient() chan<- *Client {
	return (chan<- *Client)(self.removeClient)
}

func (self *Server) SendAll() chan<- *[]byte {
	return (chan<- *[]byte)(self.sendAll)
}

func (self *Server) ReadInput() {
	origin := "ws://localhost"

	for {
		select {
		case activateInputStream := <-self.inputControl:
			log.Printf("INPUT CONTROL received = ", activateInputStream)
			if activateInputStream && !self.inputStreamActive {
				log.Printf("Opening input stream")
				ws, err := websocket.Dial(self.inputStreamURI, "", origin)
				if err != nil {
					log.Fatal(err)
				}
				self.inputWS = ws
				var msg = make([]byte, 8)
				var n int
				if n, err = self.inputWS.Read(msg); err != nil {
					log.Fatal(err)
				}
				self.handShakeMessage = msg
				log.Printf("Received JSMPEG handshake message from stream: %s %s", msg, n)

				for _, c := range self.clients {
					c.chHandShake <- self.handShakeMessage
				}
				self.inputStreamActive = true

			} else if self.inputStreamActive {
				log.Printf("Closing input stream")
				if self.inputWS != nil {
					self.inputWS.Close()
				}
				self.inputStreamActive = false
			}
		default:
			var n int
			var err error
			var msg = make([]byte, 1024)
			if self.inputStreamActive {
				if n, err = self.inputWS.Read(msg); err != nil {
					log.Fatal(err, n)
				}
				self.SendAll() <- &msg
			} /* else {
				log.Println("Waiting for connections")
			}*/
		}
	}
}

func (self *Server) Listen() {

	onConnected := func(ws *websocket.Conn) {
		client := NewClient(ws, self)
		self.addClient <- client
		client.Listen()
		defer ws.Close()
	}
	http.Handle(self.path, websocket.Handler(onConnected))

	for {
		select {

		case c := <-self.addClient:
			log.Println("Added new client ", c.id)
			self.clients = append(self.clients, c)
			if !self.inputStreamActive {
				self.inputControl <- true
			} else {
				c.Write() <- &self.handShakeMessage
			}

		case c := <-self.removeClient:
			log.Println("Remove client ", c.id)
			for i := range self.clients {
				if self.clients[i] == c {
					self.clients = append(self.clients[:i], self.clients[i+1:]...)
					break
				}
			}
			if len(self.clients) == 0 {
				log.Println("Closing source stream")
				self.inputControl <- false
			}
			log.Println("Active clients count is ", len(self.clients))

		case msg := <-self.sendAll:
			for _, c := range self.clients {
				c.Write() <- msg
			}
		}
	}
}
