package game

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	"github.com/lib/pq"
	"github.com/cristipercu/societee/bsocietee/types"
	"github.com/cristipercu/societee/bsocietee/utils"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateGameSettings(gameSettings types.GameSettings) (int64, error) {

	rows, err := s.db.Query(`SELECT * FROM public.game_settings WHERE number_of_words_per_player = $1 AND time_per_player = $2 LIMIT 1`, gameSettings.NumberOfWordsPerPlayer, gameSettings.TimePerPlayer)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	if rows.Next() {
		var existingGameSettings types.GameSettings
		err = rows.Scan(&existingGameSettings.ID, &existingGameSettings.NumberOfWordsPerPlayer, &existingGameSettings.TimePerPlayer)
		if err != nil {
			return 0, err
		}
		return existingGameSettings.ID, nil
	} else {
		var lastID int64
		err := s.db.QueryRow(`INSERT INTO public.game_settings(number_of_words_per_player, time_per_player) VALUES ($1, $2) RETURNING id`, gameSettings.NumberOfWordsPerPlayer, gameSettings.TimePerPlayer).Scan(&lastID)
		if err != nil {
			return 0, nil
		}
		return lastID, nil
	}

}

func (s *Store) GetGameSettings(gameID int) (*types.GameSettings, error) {
	gameSettings := new(types.GameSettings)

	err := s.db.QueryRow(`SELECT game_settings.id, number_of_words_per_player, time_per_player 
	FROM public.game_settings JOIN public.game
	ON game_settings.id = game.game_settings_id
	WHERE game.id = $1`, gameID).Scan(&gameSettings.ID, &gameSettings.NumberOfWordsPerPlayer, &gameSettings.TimePerPlayer)

	if err != nil {
		return gameSettings, err
	}

	return gameSettings, nil
}

func (s *Store) CreateTeam(team types.Team) (int, error) {
	var teamID int
	err := s.db.QueryRow(`INSERT INTO public.teams(name, members) VALUES ($1, $2) RETURNING id`, team.Name, pq.Array(team.MembersName)).Scan(&teamID)
	if err != nil {
		return teamID, err
	}
	return teamID, nil
}

func (s *Store) GetTeamByID(teamID int32) (types.Team, error) {
	var team types.Team

	err := s.db.QueryRow(`SELECT id, name, members FROM public.teams WHERE id = $1`, teamID).Scan(&team.ID, &team.Name, pq.Array(&team.MembersName))
	if err != nil {
		return team, err
	}
	return team, nil
}

func (s *Store) CreateGame(game types.Game) (int, error) {
	var gameID int
	err := s.db.QueryRow(`INSERT INTO public.game(game_settings_id, game_state_id, teams_id) VALUES ($1, $2, $3) RETURNING id`, game.GameSettingsID, game.GameStateID, pq.Array(game.TeamsID)).Scan(&gameID)
	if err != nil {
		return gameID, err
	}
	return gameID, nil
}

func (s *Store) CreateGameState(teams []types.Team) (int, error) {
	// var gameStateID int

	var player_order []string

	maxIndex := 0

	// Find the maximum index among all slices
	for _, slice := range teams {
		if len(slice.MembersName) > maxIndex {
			maxIndex = len(slice.MembersName)
		}
	}

	// Interleave elements from the slices
	for i := 0; i < maxIndex; i++ {
		for _, slice := range teams {
			if i < len(slice.MembersName) {
				player_order = append(player_order, slice.MembersName[i])
			}
		}
	}

	var gameStateID int
	randomNumber := rand.Intn(10)

	err := s.db.QueryRow(`INSERT INTO public.game_state(round_id, player_order, current_player_index)
	VALUES (1, $1, $2) RETURNING id`, pq.Array(player_order), randomNumber).Scan(&gameStateID)
	if err != nil {
		return gameStateID, err
	}
	return gameStateID, nil
}

func (s *Store) GetCurrentPlayer(gameID int) (string, error) {
	var currentPlayer string
	err := s.db.QueryRow(`SELECT 
	CASE 
		WHEN current_player_index <= array_length(player_order, 1) THEN player_order[current_player_index]
		WHEN current_player_index % array_length(player_order, 1) = 0 THEN player_order[array_length(player_order, 1)] 
	ELSE player_order[current_player_index % array_length(player_order, 1)] 
	END AS current_player
	
	FROM public.game_state
	JOIN game ON game.game_state_id = game_state.id
	WHERE game.id = $1`, gameID).Scan(&currentPlayer)
	if err != nil {
		return currentPlayer, err
	}
	return currentPlayer, nil

}

func (s *Store) UpdateCurrentPlayer(gameID int) (string, error) {
	var currentPlayer string
	_, err := s.db.Exec(`UPDATE public.game_state
	SET current_player_index = current_player_index + 1 
	FROM game
	WHERE game.game_state_id = game_state.id
	AND game.id = $1`, gameID)

	if err != nil {
		return currentPlayer, err
	}

	currentPlayer, err = s.GetCurrentPlayer(gameID)
	if err != nil {
		return currentPlayer, err
	}
	return currentPlayer, nil
}

func (s *Store) GetCurrentRound(gameID int) (types.Round, error) {
	var currentRound types.Round

	err := s.db.QueryRow(`SELECT round_id, name
	FROM public.game_state
	JOIN game ON game.game_state_id = game_state.id
	JOIN game_rounds ON game_state.round_id = game_rounds.id
	WHERE game.id = $1`, gameID).Scan(&currentRound.ID, &currentRound.RoundName)
	if err != nil {
		return currentRound, err
	}
	return currentRound, nil
}

func (s *Store) UpdateCurrentRound(gameID int) (types.Round, error) {

	var currentRound types.Round

	currentRound, err := s.GetCurrentRound(gameID)
	if err != nil {
		return currentRound, nil
	}

	if currentRound.ID == 5 {
		return currentRound, nil
	}
	_, err = s.db.Exec(`UPDATE public.game_state
	SET round_id = round_id + 1
	FROM game
	WHERE game.game_state_id = game_state.id
	AND game.id = $1`, gameID)

	if err != nil {
		return currentRound, err
	}

	currentRound, err = s.GetCurrentRound(gameID)
	if err != nil {
		return currentRound, err
	}
	return currentRound, nil
}

func (s *Store) UpsertUserWords(gameID int, playerName string, words []string) (int, error) {
	var wordsID int

	gameSettings, err := s.GetGameSettings(gameID)

	if err != nil {
		return wordsID, err
	}

	nrOfWordsProvided := len(words)

	if nrOfWordsProvided > int(gameSettings.NumberOfWordsPerPlayer) {
		return wordsID, fmt.Errorf("to many words provided, max allowed based on game settings %v", gameSettings.NumberOfWordsPerPlayer)
	}

	log.Println("number of words ", gameSettings.NumberOfWordsPerPlayer)

	err = s.db.QueryRow(`INSERT INTO public.user_words(game_id, player_name, words)
	VALUES($1, $2, $3)
	ON CONFLICT (game_id, player_name) DO UPDATE
		SET words = $3
	RETURNING id`, gameID, playerName, pq.Array(words)).Scan(&wordsID)

	if err != nil {
		return wordsID, err
	}

	return wordsID, nil

}

func (s *Store) GetUserWords(gameID int, playerName string) ([]string, error) {
	var words []string

	err := s.db.QueryRow(`SELECT words FROM user_words WHERE game_id = $1 and player_name = $2;`, gameID, playerName).Scan(pq.Array(&words))

	if err != nil {
		return nil, err
	}

	return words, nil
}

func (s *Store) GetGameById(gameID int) (types.Game, error) {
	var game types.Game

	err := s.db.QueryRow(`SELECT id, game_settings_id, game_state_id, teams_id, 
	CASE WHEN game_words_state_id IS NULL THEN 0 ELSE game_words_state_id END AS game_words_state_id 
	FROM public.game WHERE id = $1;`, gameID).Scan(
		&game.ID, &game.GameSettingsID, &game.GameStateID, pq.Array(&game.TeamsID), &game.WordsStateID,
	)

	if err != nil {
		return game, err
	}

	return game, nil
}

func (s *Store) SetGameWordsState(gameID int, gameWordsStateID int) error {

	_, err := s.db.Exec(`UPDATE public.game SET game_words_state_id = $1 WHERE id = $2`, gameWordsStateID, gameID)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) CreateGameWordsState(gameID int, words []string) (int, error) {
	var gameWordsStateID int
	fmt.Println(words)
	utils.ShuffleStrings(words)

	fmt.Println(words)
	err := s.db.QueryRow(`
	INSERT INTO public.game_words_state(words_order, current_word_index)
	VALUES($1, 1) RETURNING id`, pq.Array(words)).Scan(&gameWordsStateID)

	if err != nil {
		return 0, err
	}

	return gameWordsStateID, nil

}

func (s *Store) GetCurrentWord(gameID int) (string, error) {
	var currentPlayer string
	err := s.db.QueryRow(`SELECT 
	CASE 
		WHEN current_word_index <= array_length(words_order, 1) THEN words_order[current_word_index]
	ELSE 'error_no_more_words'
	END AS current_word
	FROM public.game_words_state
	JOIN game ON game.game_words_state_id = game_words_state.id
	WHERE game.id = $1`, gameID).Scan(&currentPlayer)
	if err != nil {
		return currentPlayer, err
	}
	return currentPlayer, nil

}
