package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jesusrmoreno/sad-squid"
	"github.com/jesusrmoreno/uuid-service/implementation"
	"github.com/jesusrmoreno/uuid-service/interfaces"
	"github.com/satori/go.uuid"
)

var db intf.IDStore

const (
	uuidV4      = "uuid_v4"
	simplesquid = "simplesquid"
	squidID     = "squid"

	errInvalidIDType = "Invalid ID type"
)

func getID(kind string) (string, error) {
	switch kind {
	case uuidV4:
		return uuid.NewV4().String(), nil
	case simplesquid:
		return squid.GenerateSimpleID(), nil
	case squidID:
		return squid.GenerateID(), nil
	default:
		return "", errors.New(errInvalidIDType)
	}
}

func idRequestHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	idKind := mux.Vars(r)["idType"]
	id, err := getID(idKind)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, id)
	fmt.Printf("[GeneratedID] %s, %s\n", ip, id)
}

func idInNameSpaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	ns := vars["namespace"]
	id := vars["uuid"]

	isMember, err := db.Contains(ns, id)
	if err != nil {
		fmt.Println("[Error] Database error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !isMember {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("false"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("true"))
	return
}

func allInNameSpaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ns := vars["namespace"]
	ids, err := db.All(ns)
	if err != nil {
		fmt.Println("[Error] Database error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := ""
	for _, id := range ids {
		response += id + "\n"
	}
	fmt.Fprintf(w, response)
}

func namespacedIDHandler(w http.ResponseWriter, r *http.Request) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Println("[Error] Could not parse IP Address:", r.RemoteAddr)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	ns := mux.Vars(r)["namespace"]
	idKind := mux.Vars(r)["idType"]
	idStr, err := getID(idKind)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := db.Store(ns, idStr); err != nil {
		fmt.Println("[Error] Could not store id for namespace:", ns, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, idStr)
	fmt.Printf("[GeneratedID] %s, %s\n", ip, idStr)
}

func main() {
	port := flag.String("port", "3000", "The port to run on")
	path := flag.String("dbPath", "namespaces.db", "The database name")
	var err error
	db, err = impl.NewBoltStore(*path)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	// Start our router
	r := mux.NewRouter()
	r.StrictSlash(true)
	r.HandleFunc("/new/{idType}", idRequestHandler).
		Methods("GET")
	r.HandleFunc("/{namespace}/", allInNameSpaceHandler).
		Methods("GET")
	r.HandleFunc("/{namespace}/new/{idType}", namespacedIDHandler).
		Methods("GET")
	r.HandleFunc("/{namespace}/contains/{uuid}", idInNameSpaceHandler).
		Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(assetFS()))
	n := negroni.Classic()
	n.UseHandler(r)
	n.Run(":" + *port)
}
