package main

import (
	"github.com/mark-by/tp-db-bykhovets/application/app"
	"github.com/mark-by/tp-db-bykhovets/infrastructure/persistance"
	"github.com/mark-by/tp-db-bykhovets/interfaces/handlers"
)

func main() {
	handlers.SetApp(app.New(persistance.New()))
	handlers.ServeAPI()
}
