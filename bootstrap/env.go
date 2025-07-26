package bootstrap

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Env struct {
	ServerAddress string
	Port          string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPass        string
}

func NewEnv() *Env {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env file found, assuming production environment")
	}

	viper.AutomaticEnv()

	return &Env{
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		Port:          viper.GetString("PORT"),
		DBHost:        viper.GetString("DB_HOST"),
		DBPort:        viper.GetString("DB_PORT"),
		DBName:        viper.GetString("DB_NAME"),
		DBUser:        viper.GetString("DB_USER"),
		DBPass:        viper.GetString("DB_PASS"),
	}
}
