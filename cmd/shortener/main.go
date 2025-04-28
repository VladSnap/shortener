package main

import (
	"fmt"
	"os"

	"github.com/VladSnap/shortener/internal/app"
	"github.com/VladSnap/shortener/internal/config"
	"github.com/VladSnap/shortener/internal/config/validation"
	"github.com/VladSnap/shortener/internal/log"
	"github.com/VladSnap/shortener/internal/services"
	"go.uber.org/zap"
)

var resourceManager *services.ResourceManager

func main() {
	logWorkDir(false)
	resourceManager = services.NewResourceManager()
	defer func() {
		err := resourceManager.Cleanup()

		if err != nil {
			panic(fmt.Errorf("failed resourceManager clean: %w", err))
		}
	}()
	// Регаем функцию Sync Zap логов
	resourceManager.Register(log.Close)

	log.Zap.Info("run shorneter server", zap.Strings("Args", os.Args))

	confValidator := &validation.OptionsValidator{}
	opts, err := config.LoadConfig(confValidator)
	if err != nil {
		panic(err)
	}

	server, err := app.CreateServer(opts, resourceManager)
	if err != nil {
		panic(err)
	}

	err = server.RunServer()

	if err != nil {
		log.Zap.Error("failed stop server", zap.Error(err))
	}
	log.Zap.Info("main.go end")
}

func logWorkDir(isPrint bool) {
	if !isPrint {
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		log.Zap.Error("failed get workdir", zap.Error(err))
		return
	}

	log.Zap.Info("workdir", zap.String("path", dir))
	log.Zap.Info("print all subdirs:")

	entries, err := os.ReadDir(dir)
	if err != nil {
		log.Zap.Error("failed read workdir", zap.Error(err))
	}

	for _, e := range entries {
		log.Zap.Info(e.Name())
	}
}
