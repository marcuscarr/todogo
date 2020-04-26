package main

import (
	"log"
	"os"

	"github.com/marcuscarr/todogo/server"
)

const (
	port     = 8180
	dbhost   = "localhost"
	dbport   = 5432
	dbuser   = "postgres"
	password = "mysecretpassword"
	dbname   = "marcus_demo"
)

func main() {
	log.Print("Starting server...")

	useDBHost := os.Getenv("DB_HOST")
	if useDBHost == "" {
		useDBHost = dbhost
	}
	s := server.New(server.Config{
		Port:       port,
		DBHost:     useDBHost,
		DBPort:     dbport,
		DBUser:     dbuser,
		DBPassword: password,
		DBName:     dbname,
	})
	log.Fatal(s.Start())
}
