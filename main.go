package main
import ("fmt"
	"log"
	"net/http")


var appHub = newHub()
var badResponse = []byte("This is an endpoint for WebSocket upgrade requests. You must call it in a way that supports the WebSocket protocol.")
const addr = "0.0.0.0:8080"

func upgradeRequest(w http.ResponseWriter, r *http.Request) {
	if err := upgradeConn(w, r); err != nil {
		fmt.Println("couldn't upgrade the con: ", err)
		w.Write(badResponse)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/main.html")
}

func logMiddleware(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s @ %s", r.Method ,r.URL)
		fn(w, r)
	}
}

func authMiddleware(fn func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	fmt.Println("nyet")
}

func main() {
	go appHub.run()
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", logMiddleware(serveFile))
	mux.HandleFunc("/ws", logMiddleware(upgradeRequest))
	fmt.Println("listening on address", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		panic(err)
	}
}
