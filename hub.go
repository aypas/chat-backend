package main
import (
	"fmt"
	"sync"
)

type hub struct {
	clients 	map[string]*client
	lock 		sync.RWMutex
	broadcast 	chan *message  
	register 	chan *client 
	unregister 	chan *client
}

func newHub() *hub {
	return &hub{
		clients: make(map[string]*client),
		lock: sync.RWMutex{},
		broadcast: make(chan *message), 
		register: make(chan *client),
		unregister: make(chan *client),
	}
}

func (h *hub) run() {
	for {
		select {
		case connStruct := <- h.register:
			h.lock.Lock()
			h.clients[connStruct.name] = connStruct
			go sendLogs(connStruct.name, connStruct)
			h.lock.Unlock()
			
		case connStruct := <- h.unregister:
			h.lock.Lock()
			_, ok := h.clients[connStruct.name]
			if ok {
				delete(h.clients, connStruct.name)
			}
			h.lock.Unlock()

		case msg := <- h.broadcast:
			h.lock.RLock()
			ch, ok := h.clients[msg.To]
			if ok {
				ch.writeTo <- msg
				fmt.Println("sent to", msg.To)
			} else {
				go storeMsg(msg)
			}
			h.lock.RUnlock()
		}
	}
}