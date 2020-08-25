package main
import ("fmt"
		"net/http")

var appHub = newHub()

func upgradeRequest(w http.ResponseWriter, r *http.Request) {
	if err := upgradeConn(w, r); err != nil {
		fmt.Println("couldn't upgrade the con: ", err)
	}
}

func ServeFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/main.html")
}

func main() {
	go appHub.run()
	http.HandleFunc("/", ServeFile)
	http.HandleFunc("/ws", upgradeRequest)
	err := http.ListenAndServe("0.0.0.0:8080", nil)
	if err != nil {
		panic(err)
	}
}
