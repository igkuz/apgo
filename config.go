package apgo

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type APConfig struct {
	DB map[string]string
}

func NewConfig() *APConfig {
	env := os.Getenv("APP_ENV")

	if env == "" {
		log.Fatal("You should provide APP_ENV environmental variable for correct config loading")
	}

	err := godotenv.Load(".env." + env)

	if err != nil {
		log.Println("Error while loading config: ", err)
		log.Fatal("Config was not loaded, please check .env file for APP_ENV=", env)
	}

	return &APConfig{DB: getDbConfig()}
}

func (c *APConfig) GetDbString() string {
	return c.DB["user"] + ":" + c.DB["password"] + "@/" + c.DB["name"] + "?parseTime=true"
}

func getDbConfig() map[string]string {
	db := make(map[string]string)
	db["name"] = os.Getenv("DB_NAME")
	db["user"] = os.Getenv("DB_USER")
	db["password"] = os.Getenv("DB_PASSWORD")
	db["dialect"] = os.Getenv("DB_DIALECT")
	return db
}
