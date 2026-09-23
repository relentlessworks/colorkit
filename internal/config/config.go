package config

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

// Config holds all application configuration.
type Config struct {
	Addr   string
	Secret string
}

// Default returns the default configuration.
func Default() *Config {
	return &Config{
		Addr:   ":7100",
		Secret: "",
	}
}

// Load parses flags and env vars, layered as: defaults < env < flags.
func Load() *Config {
	c := Default()

	// Env vars
	if v := os.Getenv("COLORKIT_ADDR"); v != "" {
		c.Addr = v
	}
	if v := os.Getenv("COLORKIT_SECRET"); v != "" {
		c.Secret = v
	}

	// Flags
	flag.StringVar(&c.Addr, "addr", c.Addr, "listen address")
	flag.StringVar(&c.Secret, "secret", c.Secret, "auth token signing secret (auto-generated if empty)")
	flag.Parse()

	// Generate random secret if not provided
	if c.Secret == "" {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			panic(fmt.Sprintf("failed to generate secret: %v", err))
		}
		c.Secret = hex.EncodeToString(b)
	}

	return c
}
