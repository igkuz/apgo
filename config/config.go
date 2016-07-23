package config

import (
    "github.com/joho/godotenv"
    "os"
    "log"
)

type APConfig struct {
  DB    map[string]string
}

func NewConfig() (*APConfig) {
  env := os.Getenv("APP_ENV")

  if env == "" {
    log.Fatal("You should provide APP_ENV environmental variable for correct config loading")
  }

  err := godotenv.Load(".env." + env)

  if err != nil {
    log.Fatal("Config was not loaded, please check .env file for APP_ENV=", env)
  }

  return &APConfig{ DB: getDbConfig() }
}

func getDbConfig() (map[string]string) {
  db := make(map[string]string)
  db["name"] = os.Getenv("DB_NAME")
  return db
}

