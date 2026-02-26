package main

import (
	"app/internal/server"
)

func main() {
	e := server.New()
	e.Logger.Fatal(e.Start(":8000"))
}
