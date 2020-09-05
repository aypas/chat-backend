package main
import ("fmt"
		"sync")

type safeClientMap struct {
	clients map[string]*client
	lock sync.RWMutex
}

type hub struct {
	clients *safeClientMap
	broadcast chan *message  
	register chan *client 
	unregister chan *client
}

func newHub() hub {
	return hub{
		clients: &safeClientMap{ clients: make(map[string]*client), lock: sync.RWMutex{} },
		broadcast: make(chan *message), 
		register: make(chan *client),
		unregister: make(chan *client),
	}
}

func (h hub) run() {
	for {
		select {
		case connStruct := <- h.register:
			h.clients.lock.Lock()
			h.clients.clients[connStruct.name] = connStruct
			h.clients.lock.Unlock()
			
		case connStruct := <- h.unregister:
			h.clients.lock.Lock()
			_, ok := h.clients.clients[connStruct.name]
			if ok {
				delete(h.clients.clients, connStruct.name)
			}
			h.clients.lock.Unlock()

		case msg := <- h.broadcast:
			h.clients.lock.RLock()
			fmt.Println(*msg, " got to the hub.")
			fmt.Println("msg meant for ", msg.To)
			ch, ok := h.clients.clients[msg.To]
			if ok {
				ch.writeTo <- msg
				fmt.Println("sent to", msg.To)
			}
			h.clients.lock.RUnlock()
		}
	}
}