package main

import (
	"fmt"
	"os"

	"github.com/VladSnap/shortener/internal/app"
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/config/validation"
)

func main() {
	fmt.Println("Args:", os.Args)
	opts := config.LoadConfig()
	err := app.RunServer(opts)
	fmt.Println("Run shorneter server. Args:", os.Args)
	confValidator := &validation.OptionsValidator{}
	opts, err := config.LoadConfig(confValidator)
	if err != nil {
		panic(err)
	}
}
