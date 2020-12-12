package main

import (
	"github.com/mark-by/tp-db-bykhovets/infrastructure/persistance"
	"github.com/mark-by/tp-db-bykhovets/interfaces/handlers"
)

func main() {
	repositories := persistance.New()

	handlers.ServeAPI(repositories)
}