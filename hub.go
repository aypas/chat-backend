package main
import ("fmt"
		"encoding/json")

type hub struct {
	clients map[string]client
	broadcast chan []byte 
	register chan *client 
	unregister chan *client
}

func newHub() hub {
	return hub{
		clients: make(map[string]client),
		broadcast: make(chan []byte),
		register: make(chan *client),
		unregister: make(chan *client),
	}
}

func (h hub) run() {
	for {
		select {
		case connStruct := <- h.register:
			h.clients[connStruct.name] = *connStruct
		case connStruct := <- h.unregister:
			_, ok := h.clients[connStruct.name]
			if ok {
				delete(h.clients, connStruct.name)
			}
		case msg := <- h.broadcast:
			var data *message = &message{} 
			fmt.Println(string(msg), " got to the hub.")
			json.Unmarshal(msg, data)
			fmt.Println("msg meant for ", data.To)
			ch, ok := h.clients[data.To]
			if ok {
				ch.writeTo <- msg
				fmt.Println("sent to", data.To)
			}
		}
	}
}