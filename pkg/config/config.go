package config

import (
	"os"
)

const HostnameEnvVar = "HOSTNAME"
const PodNameEnvVar = "POD_NAME"

var LogLevel = os.Getenv("LOG_LEVEL")

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
