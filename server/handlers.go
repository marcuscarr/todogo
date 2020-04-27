package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/marcuscarr/todogo/todo"
)

type request struct {
	ID          *int    `json:"ID"`
	Title       *string `json:"Title"`
	Description *string `json:"Description"`
	Status      *bool   `json:"Status"`
}

func (s *Server) getOneTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		returnError(w, r, err)
		return
	}
	t, err := todo.Retrieve(s.db, id)
	if err != nil {
		log.Print("error find todo with id:", id, " err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if t == nil {
		log.Print("did not find todo with id:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	log.Print("found todo with id:", id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)
	json.NewEncoder(w).Encode(t)
}

func (s *Server) createTodo(w http.ResponseWriter, r *http.Request) {
	var req request
	err := decodeJSONBody(w, r, &req)
	if err != nil {
		returnError(w, r, err)
		return
	}
	id, err := todo.Create(s.db, nil, req.Title, req.Description, req.Status)
	if err != nil {
		log.Print("error creating todo:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := map[string]int{"id": id}
	log.Print("created todo; sending result:", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) updateTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		returnError(w, r, err)
		return
	}
	var req request
	err = decodeJSONBody(w, r, &req)
	if err != nil {
		returnError(w, r, err)
		return
	}
	count, err := todo.Update(s.db, id, req.Title, req.Description, req.Status)
	if err != nil {
		log.Print("error updating todo with id:", id, " err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if count == 0 {
		log.Print("did not find todo with id:", id)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	res := map[string]int{"updated": count}
	log.Print("updated todo; sending result:", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) upsertTodo(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		returnError(w, r, err)
		return
	}
	var req request
	err = decodeJSONBody(w, r, &req)
	if err != nil {
		returnError(w, r, err)
		return
	}
	id, err = todo.Upsert(s.db, id, req.Title, req.Description, req.Status)
	if err != nil {
		log.Print("error upserting todo with id:", id, " err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := map[string]int{"id": id}
	log.Print("upserted todo; sending result:", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) deleteTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		returnError(w, r, err)
		return
	}
	count, err := todo.Delete(s.db, id)
	if err != nil {
		log.Print("error deleting todo with id:", id, "err:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	res := map[string]int{"deleted": count}
	log.Print("deleted todo; sending result:", res)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (s *Server) getAllTodos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	todos, err := todo.GetAll(s.db)
	if err != nil {
		log.Print("error getting all todos:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	log.Print("sending all todos; count:", len(todos))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(todos)
}

func (s *Server) statusCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	err := s.db.Ping()
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status": "no database connection"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ready"}`))
}
