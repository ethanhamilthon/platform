package config

import "os"

var (
	DbPath string
	Mode   string
)

func init() {
	Init()
}

func Init() {
	Mode = loadEnv(os.Getenv("MODE"), "dev")
	DbPath = loadEnv(os.Getenv("DB_PATH"), "./data/main.db")

	messageInit()
}

func loadEnv(envValue string, defaultVal string) string {
	if envValue == "" {
		return defaultVal
	}
	return envValue
}
