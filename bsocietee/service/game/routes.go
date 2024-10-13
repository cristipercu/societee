package game

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/cristipercu/societee/bsocietee/types"
	"github.com/cristipercu/societee/bsocietee/utils"
)

type Handler struct {
	store     types.GameStore
	userStore types.UserStore
	roomStore types.RoomStore
}

func NewHandler(store types.GameStore, userStore types.UserStore, roomStore types.RoomStore) *Handler {
	return &Handler{store: store, userStore: userStore, roomStore: roomStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/games", h.handleCreateGame).Methods(http.MethodPost)
	router.HandleFunc("/games/settings", h.handleCreateSettings).Methods(http.MethodPost)
	router.HandleFunc("/games/settings", h.handleGetGameSettings).Methods(http.MethodGet)
	router.HandleFunc("/games/teams", h.handleCreateTeam).Methods(http.MethodPost)
	router.HandleFunc("/games/gameState/currentPlayer", h.handleGetCurentPlayer).Methods(http.MethodGet)
	router.HandleFunc("/games/gameState/currentPlayer", h.handleUpdateCurentPlayer).Methods(http.MethodPut)
	router.HandleFunc("/games/gameState/currentRound", h.handleGetCurrentRound).Methods(http.MethodGet)
	router.HandleFunc("/games/gameState/currentRound", h.handleUpdateCurrentRound).Methods(http.MethodPut)
	router.HandleFunc("/games/words/add", h.handleUpsertWords).Methods(http.MethodPost)
	router.HandleFunc("/games/words/get", h.handleGetUserWords).Methods(http.MethodGet)
	router.HandleFunc("/games/startGame", h.handleStartGame).Methods(http.MethodPost)
	router.HandleFunc("/games/words/getCurrentWord", h.handleGetCurrentWord).Methods(http.MethodGet)
}

func (h *Handler) handleCreateGame(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateGamePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	teams := []types.Team{}

	uniqueTeamsID := make([]interface{}, len(payload.TeamsID))
	for i, id := range payload.TeamsID {
		uniqueTeamsID[i] = id
	}

	if !utils.IsUniqueNaive(uniqueTeamsID) {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("teams id must be unique"))
		return
	}

	for _, teamid := range payload.TeamsID {
		team, err := h.store.GetTeamByID(teamid)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error getting team: %d", teamid))
			return
		}
		teams = append(teams, team)
	}

	gameStateID, err := h.store.CreateGameState(teams)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	game := types.Game{
		GameSettingsID: payload.GameSettingsID,
		GameStateID:    gameStateID,
		TeamsID:        payload.TeamsID,
	}

	game.ID, err = h.store.CreateGame(game)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, game)

}

// Create Game Settings
//
//	@Summary		create a game settings entry in the database
//	@Description	if the game settings already exist, return the existing entry
//	@Tags			game
//	@Accept			json
//	@Produce		json
//
//	@Param			payload	body		types.CreateGameSettingsPayload	true	"payload"
//
//	@Success		200		{object}	types.GameSettings
//	@Router			/games/settings [post]
func (h *Handler) handleCreateSettings(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateGameSettingsPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	settingsID, err := h.store.CreateGameSettings(types.GameSettings{
		NumberOfWordsPerPlayer: payload.NumberOfWordsPerPlayer,
		TimePerPlayer:          payload.TimePerPlayer,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.GameSettings{ID: settingsID, NumberOfWordsPerPlayer: payload.NumberOfWordsPerPlayer, TimePerPlayer: payload.TimePerPlayer})
}

func (h *Handler) handleGetGameSettings(w http.ResponseWriter, r *http.Request) {
	var payload types.GameIDPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	gameSetting, err := h.store.GetGameSettings(payload.GameID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, gameSetting)

}

func (h *Handler) handleCreateTeam(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateTeamsPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	teamID, err := h.store.CreateTeam(types.Team{
		Name:        payload.Name,
		MembersName: payload.MembersName,
	})

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, types.Team{ID: teamID, Name: payload.Name, MembersName: payload.MembersName})

}

func (h *Handler) handleGetCurentPlayer(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	gameID := queryParams.Get("game_id")

	if gameID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game_id is required"))
		return
	}

	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game_id must be an integer"))
		return
	}

	currentPlayer, err := h.store.GetCurrentPlayer(gameIDInt)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"current_player": currentPlayer})
}

func (h *Handler) handleUpdateCurentPlayer(w http.ResponseWriter, r *http.Request) {
	var payload types.GameIDPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	currentPlayer, err := h.store.UpdateCurrentPlayer(payload.GameID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"current_player": currentPlayer})

}

func (h *Handler) handleGetCurrentRound(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	gameID := queryParams.Get("game_id")

	if gameID == "" {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game_id param missing"))
		return
	}

	gameIDInt, err := strconv.Atoi(gameID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("game_id should be int"))
		return
	}

	currentRound, err := h.store.GetCurrentRound(gameIDInt)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, currentRound)
}

func (h *Handler) handleUpdateCurrentRound(w http.ResponseWriter, r *http.Request) {
	var payload types.GameIDPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	currentRound, err := h.store.UpdateCurrentRound(payload.GameID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, currentRound)

}

func (h *Handler) handleUpsertWords(w http.ResponseWriter, r *http.Request) {
	var payload types.UpsertWordsPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	// userID := auth.GetUserIDFromContext(r.Context())

	wordsID, err := h.store.UpsertUserWords(payload.GameID, payload.PlayerName, payload.Words)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]int{"words_id": wordsID})
}

func (h *Handler) handleGetUserWords(w http.ResponseWriter, r *http.Request) {
	gameID, err := utils.GetGameIDFromRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("gameId error %v", err))
		return
	}

	playerName, err := utils.GetUserNameFromRequest(r)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("playerName error %v", err))
		return
	}

	words, err := h.store.GetUserWords(gameID, playerName)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string][]string{"words": words})

}

func (h *Handler) handleStartGame(w http.ResponseWriter, r *http.Request) {
	var payload types.GameIDPayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

	game, err := h.store.GetGameById(payload.GameID)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error getting the game %v", err))
		return
	}

	gameSetting, err := h.store.GetGameSettings(payload.GameID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not get game settings"))
		return
	}

	var allPlayers []string
	for _, teamID := range game.TeamsID {
		team, err := h.store.GetTeamByID(teamID)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not get team %v", teamID))
			return
		}
		allPlayers = append(allPlayers, team.MembersName...)
	}

	var allWords []string
	for _, playerName := range allPlayers {
		userWords, err := h.store.GetUserWords(payload.GameID, playerName)
		if err != nil {
			utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not get words for user %s", playerName))
			return
		}
		if len(userWords) != int(gameSetting.NumberOfWordsPerPlayer) {
			utils.WriteJSON(w, http.StatusTeapot, map[string]string{"message": "have some tea, users are still adding the words"})
			return
		}
		allWords = append(allWords, userWords...)

	}

	game.WordsStateID, err = h.store.CreateGameWordsState(payload.GameID, allWords)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not create game words state %v", err))
		return
	}

	err = h.store.SetGameWordsState(payload.GameID, game.WordsStateID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not set game words state %v", err))
		return
	}

	_, err = h.store.UpdateCurrentRound(payload.GameID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("error on updating the round %v", err))
		return
	}

	currentPlayer, err := h.store.GetCurrentPlayer(payload.GameID)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("error on getting the current player %v", err))
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"current player": currentPlayer})

}

func (h *Handler) handleGetCurrentWord(w http.ResponseWriter, r *http.Request) {
	gameID, err := utils.GetGameIDFromRequest(r)

	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("gameId error %v", err))
		return
	}

	currentWord, err := h.store.GetCurrentWord(gameID)

	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, fmt.Errorf("could not get the current word"))
		return
	}

	if currentWord == "error_no_more_words" {
		utils.WriteError(w, http.StatusTeapot, fmt.Errorf("no more words, have a tea and get to the next round"))
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"currentWord": currentWord})

}
