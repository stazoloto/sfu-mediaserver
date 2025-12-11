package config

import "sync"

type Config struct {
	WebSocketPort  string
	EnableSSL      bool
	CertFile       string
	KeyFile        string
	AllowedOrigins []string
}

var (
	instance *Config
	once     sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		instance = &Config{
			WebSocketPort:  "8080",
			EnableSSL:      false,
			AllowedOrigins: []string{"*"},
		}
	})
	return instance
}
