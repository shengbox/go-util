package env

import "os"

func Get(key, defaultValue string) string {
	if env := os.Getenv(key); env == "" {
		return defaultValue
	} else {
		return env
	}
}
