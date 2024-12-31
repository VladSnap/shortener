package main

import (
	"fmt"
	"os"

	"github.com/VladSnap/shortener/internal/app"
	"github.com/VladSnap/shortener/internal/config"
)

func main() {
	fmt.Println("Args:", os.Args)
	opts := config.ParseFlags()
	err := app.RunServer(opts)
	if err != nil {
		panic(err)
	}
}
