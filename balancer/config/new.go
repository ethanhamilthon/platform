package config

import (
	"os"

	"github.com/joho/godotenv"
)

var (
	Mode      string
	isTesting bool
	HttpPort  string
	HttpsPort string
)

func getEnv(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func init() {
	Init()
}

func Init() {
	godotenv.Load()
	Mode = getEnv(os.Getenv("MODE"), "dev")
	isTesting = getEnv(os.Getenv("IS_TESTING"), "false") == "true"
	HttpPort = getEnv(os.Getenv("HTTP_PORT"), "80")
	HttpsPort = getEnv(os.Getenv("HTTPS_PORT"), "443")
	loadBrokerConfigs()
}

func IsTesting() bool {
	return isTesting
}

func SetTesting() {
	isTesting = true
}
