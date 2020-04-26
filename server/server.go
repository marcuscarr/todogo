package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/jmoiron/sqlx"
	// PostgreSQL driver
	_ "github.com/lib/pq"
)

// Config are the server parameters
type Config struct {
	Port       int
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
}

// Server holds the http router and the database handle
type Server struct {
	cfg    Config
	router *mux.Router
	db     *sqlx.DB
}

// Start connects to the database and serves the endpoint
func (s *Server) Start() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.cfg.DBHost, s.cfg.DBPort, s.cfg.DBUser, s.cfg.DBPassword, s.cfg.DBName)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	s.db = db

	s.router.HandleFunc("/_status", s.statusCheck)
	s.router.HandleFunc("/todo", s.createTodo).Methods("POST")
	s.router.HandleFunc("/todos/{id}", s.getOneTodo).Methods("GET")
	s.router.HandleFunc("/todos/{id}", s.upsertTodo).Methods("PUT")
	s.router.HandleFunc("/todos/{id}", s.updateTodo).Methods("PATCH")
	s.router.HandleFunc("/todos/{id}", s.deleteTodo).Methods("DELETE")
	s.router.HandleFunc("/todos", s.getAllTodos).Methods("GET")

	return http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), s.router)
}

// New returns a new server
func New(cfg Config) *Server {
	router := mux.NewRouter().StrictSlash(true)
	s := Server{
		cfg:    cfg,
		router: router,
	}

	return &s
}
