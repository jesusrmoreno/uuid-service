package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/satori/go.uuid"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

var db *ledis.DB

func idInNameSpaceHandler(w http.ResponseWriter, r *http.Request) {
	ns := r.URL.Query().Get("namespace")
	id := r.URL.Query().Get("uuid")

	isMember, err := db.SIsMember([]byte(ns), []byte(id))
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if isMember == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("false"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
	return
}

func namespacedIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ns := r.URL.Query().Get("namespace")
	id := uuid.NewV4()
	idStr := id.String()
	if _, err := db.SAdd([]byte(ns), []byte(idStr)); err != nil {
		fmt.Println("[Error] Could not store id for namespace:", ns)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, idStr)
	fmt.Printf("[GeneratedID] %s, %s\n", ip, idStr)
}

// Generates a UUID V4 and sends it to the server, also logs the ID and IP
func requestIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	id := uuid.NewV4()
	idStr := id.String()
	fmt.Fprintf(w, idStr)
	fmt.Printf("[GeneratedID] %s, %s\n", ip, idStr)
}

func main() {
	port := flag.String("port", "3000", "The port to run on")
	flag.Parse()
	cfg := lediscfg.NewConfigDefault()
	l, _ := ledis.Open(cfg)
	db, _ = l.Select(0)

	r := mux.NewRouter()
	r.HandleFunc("/", requestIDHandler)
	r.HandleFunc("/ns", namespacedIDHandler)
	r.HandleFunc("/exists", idInNameSpaceHandler)
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":" + *port)
}
