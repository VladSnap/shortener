// Package config отвечает за конфигурирование приложения shortener.
package config

import (
	"flag"
	"fmt"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options - Структура конфига приложения.
type Options struct {
	// ListenAddress - Server listen address.
	ListenAddress string `env:"SERVER_ADDRESS"`
	// BaseURL - Base url for short url
	BaseURL string `env:"BASE_URL"`
	// FileStoragePath - File path to storage all shorten url
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	// DataBaseConnString - Database connection string
	DataBaseConnString string `env:"DATABASE_DSN"`
	// AuthCookieKey - Key for signing auth cookies
	AuthCookieKey string `env:"AUTH_COOKIE_KEY"`
	// Performance - Enable pprof for performance testing
	Performance bool
}

// MarshalLogObject - Сериализует структуру конфига для эффективного логирования.
func (opts *Options) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ListenAddress", opts.ListenAddress)
	enc.AddString("BaseURL", opts.BaseURL)
	enc.AddString("FileStoragePath", opts.FileStoragePath)
	enc.AddString("DataBaseConnString", opts.DataBaseConnString)
	enc.AddBool("Performance", opts.Performance)
	enc.AddString("AuthCookieKey", opts.AuthCookieKey)
	return nil
}

// ConfigValidater - Интерфейс валидатора конфига.
type ConfigValidater interface {
	Validate(opts *Options) error
}

// LoadConfig - Загружает конфигурацию приложения, парсит флаги и переменные окружения.
func LoadConfig(validater ConfigValidater) (*Options, error) {
	opts := ParseFlags(validater)

	err := ParseEnvConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("config env parsing failed: %w", err)
	}

	err = validater.Validate(opts)

	if err != nil {
		return nil, fmt.Errorf("config validating failed: %w", err)
	}

	log.Zap.Info("Config loaded", zap.Object("config", opts))
	return opts, nil
}

// ParseFlags - Парсит консольные флаги приложения.
func ParseFlags(validater ConfigValidater) *Options {
	opts := new(Options)

	flag.StringVar(&opts.ListenAddress, "a", ":8080", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "http://localhost:8080", "base url for short url")
	flag.StringVar(&opts.FileStoragePath, "f", "", "file path to storage all shorten url")
	flag.StringVar(&opts.DataBaseConnString, "d", "", "database connection string")
	flag.BoolVar(&opts.Performance, "p", false, "pprof")
	flag.StringVar(&opts.AuthCookieKey, "k", "", "key for signing auth cookies")

	flag.Parse()

	return opts
}

// ParseEnvConfig - Парсит переменные окружения.
func ParseEnvConfig(opts *Options) error {
	err := env.Parse(opts)

	if err != nil {
		return fmt.Errorf("failed env parsing: %w", err)
	}

	return nil
}
