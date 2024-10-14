package types

import "time"

type UserStore interface {
  CreateUser(User) error
  GetUserByEmail(string) (*User, error)
  GetUserByID(int) (*UserProfile, error)
}

type User struct {
  ID int `json:"id"`
  Username string `json:"username"`
  Email string `json:"email"`
  Password string `json:"password"`
  CreatedOn string `json:"created_at"`
  ModifiedOn string `json:"modified_on"`
}

type UserProfile struct {
  ID int `json:"id"`
  Username string `json:"username"`
  Email string `json:"email"`
  CreatedOn string `json:"created_at"`
  ModifiedOn string `json:"modified_on"`
}


type RegisterUserPayload struct {
  Username string `json:"username" validate:"required"`
  Email string `json:"email" validate:"required,email"`
  Password string `json:"password" validate:"required,min=6,max=130"`
}

type LoginUserPayload struct {
  Email string `json:"email" validate:"required"`
  Password string `json:"password" validate:"required,min=6,max=130"`
}

type RoomStore interface {
	CreateRoom(Room) (int, error)
	GetRoomById(int) (Room, error)
	AddMemberToRoom(int, string) error
	GetRoomMembers(int) ([]RoomMember, error)
	LeaveRoom(int, string) error
  GetRoomByRoomCode(string) (Room, error)
}

type GameStore interface {
	CreateGameSettings(GameSettings) (int64, error)
	GetGameSettings(int) (*GameSettings, error)
	CreateTeam(Team) (int, error)
	CreateGameState([]Team) (int, error)
	GetTeamByID(int32) (Team, error)
	CreateGame(Game) (int, error)
	GetGameById(int) (Game, error)
	GetCurrentPlayer(int) (string, error)
	UpdateCurrentPlayer(int) (string, error)
	GetCurrentRound(int) (Round, error)
	UpdateCurrentRound(int) (Round, error)
	UpsertUserWords(int, string, []string) (int, error)
	GetUserWords(int, string) ([]string, error)
	CreateGameWordsState(int, []string) (int, error)
	SetGameWordsState(int, int) error
	GetCurrentWord(int) (string, error)
}

type Room struct {
	ID            int       `json:"id"`
  RoomCode      string    `json:"room_code"`
	Password      string    `json:"password"`
	CreatedAt     time.Time `json:"created_at"`
	OwnerName     string    `json:"owner_name"`
	MaxMembers    int       `json:"max_members" validate:"required,min=1,max=20"`
	CurrentGameID int       `json:"game_id"`
}

type RoomMember struct {
	RoomID     int    `json:"room_id"`
	PlayerName string `json:"player_name"`
}

type CreateRoomPayload struct {
  RoomCode   string `json:"room_code" validate:"required"`
	Password   string `json:"password"`
	OwnerName  string `json:"owner_name" validate:"required"`
	MaxMembers int    `json:"max_members" validate:"required"`
}

type GetRoomByCodePayload struct {
  RoomCode string `json:"room_code" validate:"required"`
  Password string `json:"password"`
}

type Game struct {
	ID             int     `json:"id"`
	GameSettingsID int     `json:"game_settings_id"`
	GameStateID    int     `json:"game_state_id"`
	TeamsID        []int32 `json:"teams"`
	WordsStateID   int     `json:"words_state_id"`
}

type CreateGamePayload struct {
	ID             int     `json:"id"`
	GameSettingsID int     `json:"game_settings_id"`
	TeamsID        []int32 `json:"teams_id" validate:"required,min=1"`
}

type GameSettings struct {
	ID                     int64 `json:"id"`
	NumberOfWordsPerPlayer uint  `json:"number_of_words_per_player"`
	TimePerPlayer          uint  `json:"time_per_player"`
}

type CreateGameSettingsPayload struct {
	NumberOfWordsPerPlayer uint `json:"number_of_words_per_player" validate:"required,min=1,max=20"`
	TimePerPlayer          uint `json:"time_per_player" validate:"required,min=1,max=180"`
}

// type GameState struct {
// 	ID                 int   `json:"id"`
// 	RoundID            int   `json:"round_id"`
// 	PlayerIdOrder      []int `json:"player_id_oredr"`
// 	CurrentPlayerIndex int   `json:"current_player_index"`
// }

type Team struct {
	ID          int      `json:"id"`
	Name        string   `json:"name"`
	MembersName []string `json:"members"`
}

type CreateTeamsPayload struct {
	Name        string   `json:"name" validate:"required"`
	MembersName []string `json:"members" validate:"required"`
}

type GetTeamByMembersPayload struct {
	MembersID []int64 `json:"members" validate:"required"`
}

type GameIDPayload struct {
	GameID int `json:"game_id" validate:"required"`
}

type Round struct {
	ID        int    `json:"round_id" validate:"required"`
	RoundName string `json:"round_name" validate:"required"`
}

type UpsertWordsPayload struct {
	GameID     int      `json:"game_id" validate:"required"`
	PlayerName string   `json:"player_name" validate:"required"`
	Words      []string `json:"words" validate:"required"`
}

type GetUserWordsPayload struct {
	GameID     int    `json:"game_id" validate:"required"`
	PlayerName string `json:"player_name" validate:"required"`
}
