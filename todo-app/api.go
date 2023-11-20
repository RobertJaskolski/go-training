package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddress string
}

func NewAPIServer(listenAddress string) *APIServer {
	return &APIServer{
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

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandlerFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}

// Handlers

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
