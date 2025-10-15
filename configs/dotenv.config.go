package configs

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbDriver string `mapstructure:"DB_DRIVER"`
	DbHost   string `mapstructure:"DB_HOST"`
	DbUser   string `mapstructure:"DB_USER"`
	DbPass   string `mapstructure:"DB_PASS"`
	DbName   string `mapstructure:"DB_NAME"`
	DbPort   string `mapstructure:"DB_PORT"`

	AppName string `mapstructure:"APP_NAME"`
	AppEnv  string `mapstructure:"APP_ENV"`
	ApiPort int    `mapstructure:"API_PORT"`
	ApiHost string `mapstructure:"API_HOST"`
}

// LoadEnv loads environment variables from a .env file
func LoadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}
	return nil
}

// GetEnv returns the value of an environment variable
func GetEnv(key string) string {
	LoadEnv()
	return os.Getenv(key)
}
