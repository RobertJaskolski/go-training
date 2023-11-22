package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// TYPES -------------------------------

type APIServer struct {
	listenAddress string
	store         Storage
}

// METHODS -------------------------------

func NewAPIServer(listenAddress string, store Storage) *APIServer {
	return &APIServer{
		store:         store,
		listenAddress: listenAddress,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/task", makeHTTPHandlerFunc(s.handleTask))
	router.HandleFunc("/task{id}", makeHTTPHandlerFunc(s.handleTaskWithID))

	log.Println("JSON API Server running on port: ", s.listenAddress)
	http.ListenAndServe(s.listenAddress, router)
}

// HANDLERS -------------------------------

func (s *APIServer) handleTask(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handlePostTasks(w, r)
	}

	if r.Method == "GET" {
		return s.handleGetTasks(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleTaskWithID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]

	fmt.Println(id)

	if r.Method == "GET" {
		return s.handleGetTask(w, r)
	}

	if r.Method == "PATCH" {
		return s.handlePatchTasks(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteTasks(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handlePostTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetTask(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handlePatchTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteTasks(w http.ResponseWriter, r *http.Request) error {
	return nil
}
