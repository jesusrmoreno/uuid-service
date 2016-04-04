package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
)

// Generates a UUID V4 and sends it to the server, also logs the ID and IP
func requestIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", 500)
		return
	}
	id := uuid.NewV4()
	idStr := id.String()
	fmt.Printf("[GeneratedID] %s, %s\n", ip, idStr)
	fmt.Fprintf(w, idStr)
}

func main() {
	port := flag.String("port", "3000", "The port to run on")
	flag.Parse()
	r := mux.NewRouter()
	r.HandleFunc("/", requestIDHandler)
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":" + *port)
}
