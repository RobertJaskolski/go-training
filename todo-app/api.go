package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
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
	router.HandleFunc("/task/{id}", makeHTTPHandlerFunc(s.handleTaskWithID))

	log.Println("JSON API Server running on port: ", s.listenAddress)
	err := http.ListenAndServe(s.listenAddress, router)
	if err != nil {
		return
	}
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
	intID, err := strconv.Atoi(id)

	if err != nil {
		return fmt.Errorf("invalid id given %s", id)
	}

	if r.Method == "GET" {
		return s.handleGetTask(intID, w, r)
	}

	if r.Method == "PATCH" {
		return s.handlePatchTasks(intID, w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteTasks(intID, w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handlePostTasks(w http.ResponseWriter, r *http.Request) error {
	payload := new(CreateTaskDTO)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	if payload.Title == "" || payload.Description == "" {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "Bad request body"})
	}

	task, err := s.store.CreateTask(payload)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, task)
}

func (s *APIServer) handleGetTask(id int, w http.ResponseWriter, r *http.Request) error {
	account, err := s.store.GetTaskByID(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleGetTasks(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetTasks()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handlePatchTasks(id int, w http.ResponseWriter, r *http.Request) error {
	payload := new(UpdateTaskDTO)
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return err
	}

	if payload.Title == "" && payload.Description == "" {
		return WriteJSON(w, http.StatusBadRequest, APIError{Error: "Bad request body"})
	}

	task, err := s.store.UpdateTask(id, payload)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusAccepted, task)
}

func (s *APIServer) handleDeleteTasks(id int, w http.ResponseWriter, r *http.Request) error {
	err := s.store.DeleteTask(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusAccepted, APIESuccess{Success: true})
}
