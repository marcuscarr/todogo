package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	// PostgreSQL driver
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "hello world"}`))
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "mysecretpassword"
	dbname   = "postgres"
)

type Server struct {
	router *mux.Router
	db     *sqlx.DB
}

func (s *Server) Start() error {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	s.db = db

	return http.ListenAndServe(":8080", s.router)
}

// New returns a server
func New() *Server {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health-check", healthCheckHandler)

	s := Server{
		router: router,
	}
	return &s
}
