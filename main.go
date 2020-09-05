package main
import ("fmt"
		"log"
		"net/http")


var appHub = newHub()
var badResponse = []byte("This is an endpoint for WebSocket upgrade requests. You must call it in a way that supports the WebSocket protocol.")
var secret = loadSecret()
const addr = "0.0.0.0:8008"

func originCheck(_ *http.Request) bool {
	return true
}

func upgradeRequest(w http.ResponseWriter, r *http.Request) {
	err := upgradeConn(w, r)
	if err != nil {
		fmt.Println("couldn't upgrade the con: ", err)
		w.Write(badResponse)
	}
}

func serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/main.html")
}

func logMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf(" -- %s @ %s", r.Method, r.URL)
		fn(w, r)
	}
}

func authMiddleware(fn http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		token, ok := r.URL.Query()["token"]
		if !ok {
			http.Error(w, "Make sure the query string token and name values are set.", 400)
			return
		}
		if len(token) != 1 {
			http.Error(w, "Invalid token format. There needs to one, and only one value set.", 400)
			return
		}
		if id := validate(token[0]); !id {
			fmt.Println(id)
			http.Error(w, "Token is invalid or expired", 401)
			return
		}
		fn(w, r)
	}
}

func routes() *http.ServeMux {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", logMiddleware(authMiddleware(serveFile)))
	mux.HandleFunc("/ws/", logMiddleware(authMiddleware(upgradeRequest)))
	return mux
}

func testRoutes() *http.ServeMux {
	var mux *http.ServeMux = http.NewServeMux()
	mux.HandleFunc("/", logMiddleware(serveFile))
	mux.HandleFunc("/ws/", logMiddleware(upgradeRequest))
	return mux
}

func main() {
	go appHub.run()
	var mux *http.ServeMux = testRoutes()
	fmt.Println("listening on address", addr)
	err := http.ListenAndServe(addr, mux)
	if err != nil {
		fmt.Println("ppop")
		panic(err)
	}
}
