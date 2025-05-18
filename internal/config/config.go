// Package config отвечает за конфигурирование приложения shortener.
package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/VladSnap/shortener/internal/log"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type BoolFlag struct {
	Set   bool
	Value bool
}

// Options - Структура конфига приложения.
type Options struct {
	// ListenAddress - Server listen address.
	ListenAddress string `env:"SERVER_ADDRESS" json:"server_address,omitempty"`
	// BaseURL - Base url for short url
	BaseURL string `env:"BASE_URL" json:"base_url,omitempty"`
	// FileStoragePath - File path to storage all shorten url
	FileStoragePath string `env:"FILE_STORAGE_PATH" json:"file_storage_path,omitempty"`
	// DataBaseConnString - Database connection string
	DataBaseConnString string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	// Enable https for run server
	EnableHTTPS *bool `env:"ENABLE_HTTPS" json:"enable_https,omitempty"`

	// AuthCookieKey - Key for signing auth cookies
	AuthCookieKey string `env:"AUTH_COOKIE_KEY" json:"-"`
	// Performance - Enable pprof for performance testing
	Performance *bool `json:"-"`
	// ConfigPath path to config file
	ConfigPath string `env:"CONFIG" json:"-"`
}

// MarshalLogObject - Сериализует структуру конфига для эффективного логирования.
func (opts *Options) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("ListenAddress", opts.ListenAddress)
	enc.AddString("BaseURL", opts.BaseURL)
	enc.AddString("FileStoragePath", opts.FileStoragePath)
	enc.AddString("DataBaseConnString", opts.DataBaseConnString)
	enc.AddBool("EnableHTTPS", *opts.EnableHTTPS)
	enc.AddBool("Performance", *opts.Performance)
	enc.AddString("AuthCookieKey", opts.AuthCookieKey)
	enc.AddString("ConfigPath", opts.ConfigPath)
	return nil
}

// ConfigValidater - Интерфейс валидатора конфига.
type ConfigValidater interface {
	Validate(opts *Options) error
}

// InitConfig - Загружает конфигурацию приложения, парсит флаги и переменные окружения.
func InitConfig(validater ConfigValidater) (*Options, error) {
	opts := new(Options)

	ParseFlags(opts, validater)

	err := ParseEnvConfig(opts)
	if err != nil {
		return nil, fmt.Errorf("config env parsing failed: %w", err)
	}

	var jsonOpts *Options
	if opts.ConfigPath != "" {
		jsonOpts, err = LoadJSONConfig(opts.ConfigPath)
		if err != nil {
			return nil, fmt.Errorf("config validating failed: %w", err)
		}
	}
	fmt.Println(jsonOpts)

	opts = mergeOptions(opts, jsonOpts)

	err = validater.Validate(opts)

	if err != nil {
		return nil, fmt.Errorf("config validating failed: %w", err)
	}

	log.Zap.Info("Config loaded", zap.Object("config", opts))
	return opts, nil
}

// LoadJSONConfig - Загружает конфигурацию из json файла.
// Если файл не найден, возвращает ошибку.
// Если файл не является json, возвращает ошибку.
// Если файл не содержит валидных данных, возвращает ошибку.
// Если файл содержит валидные данные, возвращает указатель на структуру Options.
func LoadJSONConfig(path string) (*Options, error) {
	opts := new(Options)

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed LoadJsonConfig: %w", err)
	}

	err = json.Unmarshal(fileData, opts)
	if err != nil {
		return nil, fmt.Errorf("failed unmarshal json config: %w", err)
	}

	return opts, nil
}

// ParseFlags - Парсит консольные флаги приложения.
func ParseFlags(opts *Options, validater ConfigValidater) {
	flag.StringVar(&opts.ListenAddress, "a", "", "server listen address")
	flag.StringVar(&opts.BaseURL, "b", "", "base url for short url")
	flag.StringVar(&opts.FileStoragePath, "f", "", "file path to storage all shorten url")
	flag.StringVar(&opts.DataBaseConnString, "d", "", "database connection string")
	flag.Func("p", "pprof", setPointerBool(func(v bool) {
		opts.Performance = new(bool)
		*opts.Performance = v
	}))
	flag.StringVar(&opts.AuthCookieKey, "k", "", "key for signing auth cookies")
	flag.Func("s", "enable https for run server", setPointerBool(func(v bool) {
		opts.EnableHTTPS = new(bool)
		*opts.EnableHTTPS = v
	}))
	flag.StringVar(&opts.ConfigPath, "c", "", "path to config file")

	flag.Parse()
}

// ParseEnvConfig - Парсит переменные окружения.
func ParseEnvConfig(opts *Options) error {
	err := env.Parse(opts)

	if err != nil {
		return fmt.Errorf("failed env parsing: %w", err)
	}

	return nil
}

// setPointerBool - Устанавливает значение указателя на bool для работы с flag-ами.
func setPointerBool(setValue func(bool)) func(s string) error {
	return func(s string) error {
		switch strings.ToLower(s) {
		case "true", "1":
			setValue(true)
		case "false", "0":
			setValue(false)
		}
		return nil
	}
}

// mergeOptions - Объединяет значения из флагов/окружения и файла конфигурации.
func mergeOptions(flagEnvOpts, fileOpts *Options) *Options {
	if fileOpts == nil {
		return flagEnvOpts
	}
	merged := *flagEnvOpts // start with flag/env values

	if merged.ListenAddress == "" && fileOpts.ListenAddress != "" {
		merged.ListenAddress = fileOpts.ListenAddress
	}
	if merged.BaseURL == "" && fileOpts.BaseURL != "" {
		merged.BaseURL = fileOpts.BaseURL
	}
	if merged.FileStoragePath == "" && fileOpts.FileStoragePath != "" {
		merged.FileStoragePath = fileOpts.FileStoragePath
	}
	if merged.DataBaseConnString == "" && fileOpts.DataBaseConnString != "" {
		merged.DataBaseConnString = fileOpts.DataBaseConnString
	}
	if merged.EnableHTTPS == nil && fileOpts.EnableHTTPS != nil {
		merged.EnableHTTPS = fileOpts.EnableHTTPS
	}
	if merged.Performance == nil && fileOpts.Performance != nil {
		merged.Performance = fileOpts.Performance
	}
	if merged.AuthCookieKey == "" && fileOpts.AuthCookieKey != "" {
		merged.AuthCookieKey = fileOpts.AuthCookieKey
	}
	if merged.ConfigPath == "" && fileOpts.ConfigPath != "" {
		merged.ConfigPath = fileOpts.ConfigPath
	}
	return &merged
}
