package main

import (
	"fmt"
	"log"

	"github.com/marcuscarr/todogo/server"
)

func main() {
	fmt.Println("Starting server...")

	s := server.New()
	log.Fatal(s.Start())
}
