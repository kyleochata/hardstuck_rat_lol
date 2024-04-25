package helper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// dir returns the absolute path of the given environment file (envFile) in the Go module's root directory.
// It searches for the go.mod file in the current directory and its parent directories until it finds it. and then appends envFile to same directory of the go.mod file.
// Panics if it can't find go.mod file
func dir(envFile string) string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		goModPath := filepath.Join(currentDir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			break
		}

		parent := filepath.Dir(currentDir)
		if parent == currentDir {
			panic(fmt.Errorf("go.mod not found"))
		}
		currentDir = parent
	}

	return filepath.Join(currentDir, envFile)
}
func LoadEnv(envFile string) {
	err := godotenv.Load(dir(envFile))
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
func Get_riot_key() string {
	LoadEnv(".env")
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY is not set or is empty")
	}
	return apiKey
}
