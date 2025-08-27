package config

import "os"

type Config struct {
	Addr string
}

func FromEnv() Config {
	addr := os.Getenv("ADDR")
	if addr == ""{
		addr = ":8080"
	}
	return Config{Addr: addr}
}
