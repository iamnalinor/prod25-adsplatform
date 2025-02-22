package config

import (
	"fmt"
	"os"
)

type Environment struct {
	ServerAddress string
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	OllamaHost    string
	OllamaModel   string
	MediaFsPath   string
	MediaBaseUrl  string
	RunningInCI   bool
}

func LoadEnvironment() Environment {
	return Environment{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		DBHost:        os.Getenv("POSTGRES_HOST"),
		DBPort:        os.Getenv("POSTGRES_PORT"),
		DBUser:        os.Getenv("POSTGRES_USERNAME"),
		DBPassword:    os.Getenv("POSTGRES_PASSWORD"),
		DBName:        os.Getenv("POSTGRES_DATABASE"),
		OllamaHost:    os.Getenv("OLLAMA_HOST"),
		OllamaModel:   os.Getenv("OLLAMA_MODEL"),
		MediaFsPath:   os.Getenv("MEDIA_FS_PATH"),
		MediaBaseUrl:  os.Getenv("MEDIA_BASE_URL"),
		RunningInCI:   os.Getenv("CI") == "true",
	}
}

func (c Environment) BuildDsn() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
