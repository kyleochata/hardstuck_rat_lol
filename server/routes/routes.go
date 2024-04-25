package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

//dir returns the absolute path of the given environment file (envFile) in the Go module's root directory.
//It searches for the go.mod file in the current directory and its parent directories until it finds it. and then appends envFile to same directory of the go.mod file.
//Panics if it can't find go.mod file

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
func get_riot_key() string {
	LoadEnv(".env")
	apiKey := os.Getenv("RIOT_API_KEY")
	if apiKey == "" {
		log.Fatal("RIOT_API_KEY is not set or is empty")
	}
	return apiKey
}

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

func getAccountByNameAndTag(w http.ResponseWriter, r *http.Request) {

	apiKey := get_riot_key()

	summonerName := r.PathValue("summoner_name")
	summonerTag := r.PathValue("tag")

	if summonerName == "" || summonerTag == "" {
		http.Error(w, "Missing path variables", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf("https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%v/%v?api_key=%v", summonerName, summonerTag, apiKey)

	response, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		responseData, _ := io.ReadAll(response.Body)
		http.Error(w, fmt.Sprintf("HTTP error: %d, Response body: %s", response.StatusCode, responseData), http.StatusBadRequest)
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
		return
	}

	var account Account
	err = json.Unmarshal(responseData, &account)
	if err != nil {
		http.Error(w, "Error parsing JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)

}

// func getMatch(puuid string) []string{

// }

// func getUserMatchList(w http.ResponseWriter, r *http.Request){

// }

func RouteSetup(router *http.ServeMux) {
	router.HandleFunc("GET /api/hello", func(w http.ResponseWriter, r *http.Request) {
		s := "hello"
		b, _ := json.Marshal(s)
		w.Write(b)
	})
	router.HandleFunc("GET /api/hello/{user_id}/", func(w http.ResponseWriter, r *http.Request) {
		userId := r.PathValue("user_id")
		w.Write([]byte(fmt.Sprintf("hello %s", userId)))
	})

	router.HandleFunc("GET /api/lol/{summoner_name}/{tag}/", getAccountByNameAndTag)
	router.HandleFunc("GET /api/lol/champions", func(w http.ResponseWriter, r *http.Request) {
		response, err := http.Get("https://ddragon.leagueoflegends.com/cdn/14.8.1/data/en_US/champion/Tristana.json")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error fetching data: %v", err), http.StatusInternalServerError)
			return
		}
		defer response.Body.Close()

		if response.StatusCode != http.StatusOK {
			responseData, _ := io.ReadAll(response.Body)
			http.Error(w, fmt.Sprintf("HTTP error: %d, Response body: %s", response.StatusCode, responseData), http.StatusBadRequest)
			return
		}

		responseData, err := io.ReadAll(response.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading response body: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(responseData)
	})
	//router.HandleFunc("GET /api/lol/{summoner_name}/{tag}/matches", getUserMatches)
}
