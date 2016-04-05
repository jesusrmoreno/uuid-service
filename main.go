package main

import (
	"errors"
	"flag"
	"fmt"
	"net"
	"net/http"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/jesusrmoreno/sad-squid"
	"github.com/satori/go.uuid"
	lediscfg "github.com/siddontang/ledisdb/config"
	"github.com/siddontang/ledisdb/ledis"
)

var db *ledis.DB

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

	isMember, err := db.SIsMember([]byte(ns), []byte(id))
	if err != nil {
		fmt.Println("[Error] Database error:", err)
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

func allInNameSpaceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	ns := vars["namespace"]
	ids, err := db.SMembers([]byte(ns))
	if err != nil {
		fmt.Println("[Error] Database error:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	response := []byte{}
	for idIndex := range ids {
		id := append(ids[idIndex], []byte("\n")...)
		response = append(response, id...)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(response)
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
	if _, err := db.SAdd([]byte(ns), []byte(idStr)); err != nil {
		fmt.Println("[Error] Could not store id for namespace:", ns, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	fmt.Fprintf(w, idStr)
	fmt.Printf("[GeneratedID] %s, %s\n", ip, idStr)
}

func main() {
	port := flag.String("port", "3000", "The port to run on")
	dbpath := flag.String("dbPath", "var", "The path to the database")
	flag.Parse()

	// Sets the database path
	cfg := lediscfg.NewConfigDefault()
	if dbpath != nil && *dbpath != "var" {
		fmt.Println(*dbpath)
		cfg.DBPath = *dbpath
	}
	l, _ := ledis.Open(cfg)
	db, _ = l.Select(0)

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
