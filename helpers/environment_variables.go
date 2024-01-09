package helpers

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func LoadEnvironmentVariables() {
	currentDir, _ := os.Getwd()
	envFilePath := filepath.Join(currentDir, ".env")
	godotenv.Load(envFilePath)
}
