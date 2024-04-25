package routes

import (
	"encoding/json"
	"fmt"
	"hardstuck_rat_lol_server/helper"
	"io"
	"log"
	"net/http"
	"sync"
)

type Account struct {
	Puuid    string `json:"puuid"`
	GameName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}
type SummonerStore struct {
	mu        sync.RWMutex
	summoners map[string]Account
}

var GlobalSummonerStore = &SummonerStore{
	summoners: make(map[string]Account),
}

func (s *SummonerStore) SaveSummoner(account Account) {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := account.GameName + "-" + account.TagLine // Composite key
	s.summoners[key] = account
}

func (s *SummonerStore) GetSummoner(gameName, tagLine string) (Account, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	key := gameName + "-" + tagLine
	summoner, found := s.summoners[key]
	if !found {
		log.Printf("Summoner with name %s and tag %s not found", gameName, tagLine)
		return Account{}, false
	}
	return summoner, true
}

func getAccountByNameAndTag(w http.ResponseWriter, r *http.Request) {
	apiKey := helper.Get_riot_key()
	summonerName := r.PathValue("summoner_name")
	summonerTag := r.PathValue("tag")

	if summonerName == "" || summonerTag == "" {
		http.Error(w, "Missing path variables", http.StatusBadRequest)
		return
	}

	url := fmt.Sprintf(`https://americas.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%v/%v?api_key=%v`, summonerName, summonerTag, apiKey)

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

	//save context for use for other riot api calls
	GlobalSummonerStore.SaveSummoner(account)

	x, _ := GlobalSummonerStore.GetSummoner(summonerName, summonerTag)
	log.Println(x)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(account)

}

// func getMatch(puuid string) []string{

// }

func getUserMatchList(w http.ResponseWriter, r *http.Request) {
	summonerName := r.PathValue("summoner_name")
	summonerTag := r.PathValue("tag")

	summoner, found := GlobalSummonerStore.GetSummoner(summonerName, summonerTag)
	log.Println(summoner.Puuid, summoner.GameName, summoner.TagLine)
	if !found {
		fmt.Println("user from session not pulled, getumatch", summoner)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summoner.Puuid)

}

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
	router.HandleFunc("GET /api/lol/{summoner_name}/{tag}/matches", getUserMatchList)
}
