package utils

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
)

func ParseJSON(r *http.Request, payload any) error {
  defer r.Body.Close()
  if r.Body == nil {
    return fmt.Errorf("request body is empty")
  }
  return json.NewDecoder(r.Body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, payload any) error {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  return json.NewEncoder(w).Encode(payload)
}

func WriteError(w http.ResponseWriter, status int, err error) {
  WriteJSON(w, status, map[string]string{"error": err.Error()})
}

var Validate = validator.New()

func GetTokenFromRequest(r *http.Request) string {
	tokenAuth := r.Header.Get("Authorization")
	tokenQuery := r.URL.Query().Get("token")

	if tokenAuth != "" {
		return tokenAuth
	}

	if tokenQuery != "" {
		return tokenQuery
	}

	return ""
}

func GetGameIDFromRequest(r *http.Request) (int, error) {
	gameIDHeader := r.Header.Get("gameId")
	gameIDQuery := r.URL.Query().Get("gameId")
	if gameIDQuery != "" {
		gameID, err := strconv.Atoi(gameIDQuery)
		return gameID, err
	}

	if gameIDHeader != "" {
		gameID, err := strconv.Atoi(gameIDHeader)
		return gameID, err
	}

	return 0, fmt.Errorf("gameId should be provided")
}

func GetUserNameFromRequest(r *http.Request) (string, error) {
	userNameQuery := r.URL.Query().Get("playerName")
	if userNameQuery != "" {
		return userNameQuery, nil
	}

	userNameHeader := r.Header.Get("playerName")
	if userNameHeader != "" {
		return userNameHeader, nil
	}

	return "", fmt.Errorf("userName should be provided")

}

func IsUniqueNaive(slice []interface{}) bool {
	for i := 0; i < len(slice); i++ {
		for j := i + 1; j < len(slice); j++ {
			if slice[i] == slice[j] {
				return false
			}
		}
	}
	return true
}

func ShuffleStrings(slice []string) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
