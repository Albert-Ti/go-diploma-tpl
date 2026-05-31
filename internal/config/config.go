package config

import (
	"flag"
	"os"
)

var Envs struct {
	RunAddr           string
	DatabaseURI       string
	JWTSecretKey      string
	AccrualSystemAddr string
}

func ParseFlag() {
	flag.StringVar(&Envs.RunAddr, "a", "localhost:8080", "address and port to run server")
	flag.StringVar(&Envs.DatabaseURI, "d", "", "connection string to DB")
	flag.StringVar(&Envs.AccrualSystemAddr, "r", "", "address of the accrual calculation system")
	Envs.JWTSecretKey = "default_secret_key"

	flag.Parse()

	if envRunAddr := os.Getenv("RUN_ADDRESS"); envRunAddr != "" {
		Envs.RunAddr = envRunAddr
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		Envs.DatabaseURI = envDatabaseURI
	}

	if envJWTSecretKey := os.Getenv("JWT_SECRET_KEY"); envJWTSecretKey != "" {
		Envs.JWTSecretKey = envJWTSecretKey
	}

	if envAccrualSystemAddr := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddr != "" {
		Envs.AccrualSystemAddr = envAccrualSystemAddr
	}
}
