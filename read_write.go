package main
import (
	"fmt"
	"time"
	"errors"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

//the update logic is retarded...fix it

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: originCheck,
}

type client struct {
	conn *websocket.Conn
	name string
	writeTo chan *message
}

func (c *client) readRoutine() {
	defer func() {
		fmt.Println("read routine stopped")
		c.conn.Close()
	}()
	fmt.Println("read routine is running")
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("there was an error ", err)
			appHub.unregister <- c
			return
		}
		var msgObj *message = unmarshalMessage(msg)
		if msgObj.To == msgObj.From {
			c.writeTo <- msgObj
			continue
		}
		appHub.broadcast <- msgObj
	}
}

func sendPeopleOnHub() ([]byte, error) {
	var people peoplePayload = peoplePayload{People: []string{}, PayloadType: "people"}
	appHub.lock.RLock()
	for key, _ := range appHub.clients {
		people.People = append(people.People, key)
	}
	appHub.lock.RUnlock()
	b, e := json.Marshal(people)
	if e != nil {
		return nil, errors.New("failed at encoding")
	}
	fmt.Println(string(b))
	return b, nil
}


func (c *client) writeRoutine() {
	fmt.Println("write routine working")
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
		fmt.Println("writeRouting stopped")
	}()

	for {
		select {
		case write, ok := <- c.writeTo:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if write.PayloadType == "people" { 
				//for now the only meta is people so lets keep it this way...
				byts, _ := sendPeopleOnHub()
				c.conn.WriteMessage(websocket.TextMessage, byts)
				fmt.Println("people hit")
				continue
			}
			byts, e := json.Marshal(write)
			if e != nil {
				fmt.Println("error on marshal")
			}
			fmt.Println(string(byts))
			c.conn.WriteMessage(websocket.TextMessage, byts)
			
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


func upgradeConn(w http.ResponseWriter, r *http.Request) error {	
	qs, ok := r.URL.Query()["id"]
	if !ok {
		http.Error(w, "Make sure query string values name and token are set.", 400)
		return nil
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("something went wrong in upgrading a conn", err)
		return err
	}
	c := &client{conn:connection, name: qs[0], writeTo: make(chan *message)}
	appHub.register <- c
	go c.readRoutine()
	go c.writeRoutine()
	return nil
}