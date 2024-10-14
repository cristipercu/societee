package room

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/cristipercu/societee/bsocietee/service/auth"
	"github.com/cristipercu/societee/bsocietee/types"
	"github.com/cristipercu/societee/bsocietee/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type Handler struct {
	store     types.RoomStore
	userStore types.UserStore
}

func NewHandler(store types.RoomStore, userStore types.UserStore) *Handler {
	return &Handler{store: store, userStore: userStore}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/rooms", auth.WithJWTAuth(h.handleCreateRoom, h.userStore)).Methods(http.MethodPost)
  router.HandleFunc("/roomCode", auth.WithJWTAuth(h.handleGetRoomByRoomCode, h.userStore)).Methods(http.MethodPost)
	router.HandleFunc("/ws/rooms/{room_id}", auth.WithJWTAuth(h.wsHandler, h.userStore)).Methods(http.MethodGet)
  router.HandleFunc("/rooms/join", auth.WithJWTAuth(h.handleJoinRoom, h.userStore)).Methods(http.MethodPost)
  router.HandleFunc("/rooms/leave", auth.WithJWTAuth(h.handleLeaveRoom, h.userStore)).Methods(http.MethodPost)
  router.HandleFunc("/rooms/members/{room_id}", auth.WithJWTAuth(h.handleGetRoomMembers, h.userStore)).Methods(http.MethodGet)
}

// Create Room
//
//		@Summary		create a new room for players to connect to
//		@Description	if the room exists IDK what happens
//		@Tags			room
//		@Accept			json
//		@Produce		json
//		@Param			payload	body	types.CreateRoomPayload	true	"payload"
//	 @Success 200 {integer} room_id "Room ID information"
//
//	@Router			/rooms [post]
func (h *Handler) handleCreateRoom(w http.ResponseWriter, r *http.Request) {
  var payload types.CreateRoomPayload
  // TODO: return utils.WriteError on error, not empty return you dumb f...
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

	room_id, err := h.store.CreateRoom(types.Room{
    RoomCode:   payload.RoomCode,
		Password:   payload.Password,
		OwnerName:  payload.OwnerName,
		MaxMembers: payload.MaxMembers,
	})


	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}


	utils.WriteJSON(w, http.StatusCreated, map[string]int{"room_id": room_id})

}

func (h *Handler) handleJoinRoom(w http.ResponseWriter, r *http.Request) {
  var payload types.RoomMember

  if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

  err := h.store.AddMemberToRoom(payload.RoomID, payload.PlayerName)

  if err != nil {
    utils.WriteError(w, http.StatusInternalServerError, err)
    return
  }

  utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "added"})
}

func (h *Handler) handleLeaveRoom(w http.ResponseWriter, r *http.Request) {
  var payload types.RoomMember
  log.Println("sunt in leave")

  if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

  err := h.store.LeaveRoom(payload.RoomID, payload.PlayerName)

  if err != nil {
    utils.WriteError(w, http.StatusInternalServerError, err)
    return
  }

  utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "removed"})
}


func (h *Handler) handleGetRoomMembers(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room_id := vars["room_id"]
	room_idInt, err := strconv.Atoi(room_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the room id should be int"))
		return
	}

  members, err := h.store.GetRoomMembers(room_idInt)

  if err != nil {
    utils.WriteError(w, http.StatusBadRequest, err)
  }

  log.Println(members)
  utils.WriteJSON(w, http.StatusOK, members)
}


type ChatRoom struct {
	connections []*websocket.Conn
	mutex       sync.Mutex
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var chatRooms = make(map[int]*ChatRoom)

// Broadcast a message to connections within a room
func (room *ChatRoom) broadcast(messageType int, message []byte) {
	room.mutex.Lock()
	defer room.mutex.Unlock()

	for _, conn := range room.connections {
		if err := conn.WriteMessage(messageType, message); err != nil {
			// Handle potential disconnections
			log.Println("Error writing to connection:", err)
		}
	}
}


func (h *Handler) handleGetRoomByRoomCode(w http.ResponseWriter, r *http.Request) {
  var payload types.GetRoomByCodePayload

	if err := utils.ParseJSON(r, &payload); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		error := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", error))
		return
	}

  room, err := h.store.GetRoomByRoomCode(payload.RoomCode)
  if err != nil {
    utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("room not found"))
    return
  }

  if room.Password != payload.Password {
    utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("wrong password"))
    return
  }

  utils.WriteJSON(w, http.StatusOK, room)
}
//		Room websocket
//
//		@Summary		ws to handle message inside a room
//		@Description	the room_id should be a url query param
//		@Tags			room
//		@Accept			json
//		@Produce		json
//	    @Param        room_id    query    int     true     "ID of the room"
//
//		@Success		200
//
//		@Router			/ws/rooms/{room_id} [get]
func (h *Handler) wsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	room_id := vars["room_id"]
	// playerName := r.Header.Get("player_name")
	// if playerName == "" {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("please provide the player name"))
	// }
	// userID := auth.GetUserIDFromContext(r.Context())
	log.Println(room_id)
	room_idInt, err := strconv.Atoi(room_id)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("the room id should be int"))
		return
	}

	// db_room, err := h.store.GetRoomById(room_idInt)
	// if err != nil {
	// 	utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("room not found"))
	// 	return
	// }

  log.Println("sunt aici")
	// err = h.store.AddMemberToRoom(db_room.ID, playerName, db_room.MaxMembers)

	if err != nil {
		http.Error(w, "Error adding member to room", http.StatusInternalServerError)
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Error upgrading to WebSocket", http.StatusInternalServerError)
		return
	}

	// u, err := h.userStore.GetUserByID(userID)
	// if err != nil {
	// 	http.Error(w, "Error getting user", http.StatusInternalServerError)
	// 	return
	// }

	// defer h.store.LeaveRoom(db_room.ID, playerName)
	// handleWebSocketConnection(conn, room_idInt, playerName)
	handleWebSocketConnection(conn, room_idInt)
}

func handleWebSocketConnection(conn *websocket.Conn, room_idInt int) {
	// Ensure room exists
	if _, ok := chatRooms[room_idInt]; !ok {
		chatRooms[room_idInt] = &ChatRoom{}
	}
	room := chatRooms[room_idInt] // Access the existing room

	// Add connection to the room
	room.mutex.Lock()
	room.connections = append(room.connections, conn)
	room.mutex.Unlock()
	defer conn.Close() // Clean up the connection

	// new_member_join_message := "joined the room"

	// room.broadcast(websocket.TextMessage, []byte(new_member_join_message))

	// Read message loop
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			break // Handle connection closure
		}

		// Broadcast message in the room
		room.broadcast(messageType, message)
	}
}
