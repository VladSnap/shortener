package main

import (
	"github.com/VladSnap/shortener/internal/app"
)

func main() {
	err := app.RunServer()
	if err != nil {
		panic(err)
	}
}
