package main
import ("fmt"
		"time"
		"bytes"
		"errors"
		"net/http"
		"encoding/json"
		"github.com/gorilla/websocket")

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
}


type client struct {
	conn *websocket.Conn
	name string
	writeTo chan []byte
}

func (c *client) readRoutine() {
	defer c.conn.Close()
	fmt.Println("read routine is running")
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("there was an error ", err)
			appHub.unregister <- c
			break
		}
		msg = bytes.TrimSpace(bytes.Replace(msg, newline, space, -1))
		fmt.Println("read msg")
		fmt.Println(string(msg))
		var deMsg *message = &message{}
		json.Unmarshal(msg, deMsg)
		if deMsg.To == "self" && deMsg.PayloadType == "meta" {
			c.writeTo <- []byte("people")
		} 
		appHub.broadcast <- msg
	}
}

func sendPeopleOnHub() ([]byte, error) {
	var array []string
	for k, _ := range appHub.clients {
		array = append(array, k)
	}

	b, e := json.Marshal(peoplePayload{People: array,
						 PayloadType: "meta"})
	if e != nil {
		return nil, errors.New("failed at encoding")
	}
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

	p, _ := sendPeopleOnHub()
	c.conn.WriteMessage(websocket.TextMessage, p)
	for {
		select {
		case write, ok := <- c.writeTo:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			fmt.Println(string(write))
			if z := bytes.Compare(write, []byte("people")); z == 0 {
				write, _ = sendPeopleOnHub()
			}
			c.conn.WriteMessage(websocket.TextMessage, write)
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}


func upgradeConn(w http.ResponseWriter, r *http.Request) (error) {
	fmt.Println("ws has been hit")
	cookie, e := r.Cookie("name")
	if e != nil {
		fmt.Println("something went wrong in getting the name cookie", e)
		return e
	}
	connection, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("something went wrong in upgrading a conn", err)
		return err
	}
	c := &client{conn:connection, name: cookie.Value, writeTo: make(chan []byte)}
	appHub.register <- c
	go c.readRoutine()
	go c.writeRoutine()
	return nil
}
