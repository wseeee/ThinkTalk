package env

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	dir, err := os.Getwd()
	if err != nil {
		_ = godotenv.Load()
		return
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	_ = godotenv.Load()
}
