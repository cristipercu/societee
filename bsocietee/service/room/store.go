package room

import (
	"database/sql"
	"fmt"

	"github.com/cristipercu/societee/bsocietee/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) CreateRoom(room types.Room, playerName string) (int, error) {
	var room_id int
	err := s.db.QueryRow(`INSERT INTO public.rooms (room_code, password, owner_name, max_members)
	VALUES ($1, $2, $3, $4) RETURNING id`, room.RoomCode, room.Password, room.OwnerName, room.MaxMembers).Scan(&room_id)
	if err != nil {
		return room_id, err
	}

	s.AddMemberToRoom(room.ID, playerName, room.MaxMembers)

	return room_id, nil
}

func (s *Store) GetRoomById(room_id int) (types.Room, error) {
	var room types.Room

	err := s.db.QueryRow(`SELECT id, room_code, password, owner_name, max_members, created_at
	FROM public.rooms WHERE id = $1`, room_id).Scan(&room.ID, &room.RoomCode, &room.Password, &room.OwnerName, &room.MaxMembers, &room.CreatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}

func (s *Store) GetRoomByRoomCode(roomCode string) (types.Room, error) {
	var room types.Room

	err := s.db.QueryRow(`SELECT id, room_code, password, owner_name, max_members, created_at
	FROM public.rooms WHERE room_code = $1`, roomCode).Scan(&room.ID, &room.RoomCode, &room.Password, &room.OwnerName, &room.MaxMembers, &room.CreatedAt)
	if err != nil {
		return room, err
	}
	return room, nil
}


func (s *Store) AddMemberToRoom(roomID int, playerName string, max_members int) error {

	current_members, err := s.GetRoomMembers(roomID)
	if err != nil {
		return err
	}

	if len(current_members) >= max_members {
		return fmt.Errorf("room is full")
	}

	for _, member := range current_members {
		if member.PlayerName == playerName {
			return nil
		}
	}

	_, err = s.db.Exec(`INSERT INTO public.room_members (room_id, player_name) VALUES ($1, $2)`, roomID, playerName)

	if err != nil {
		return err
	}
	return nil
}

func (s *Store) GetRoomMembers(roomID int) ([]types.RoomMember, error) {
	var members []types.RoomMember
	rows, err := s.db.Query(`SELECT * FROM public.room_members WHERE room_id = $1`, roomID)
	if err != nil {
		return members, err
	}
	for rows.Next() {
		var member types.RoomMember
		if err := rows.Scan(&member.RoomID, &member.PlayerName); err != nil {
			return members, err
		}
		members = append(members, member)
	}
	return members, nil
}

func (s *Store) LeaveRoom(roomID int, playerName string) error {
	_, err := s.db.Exec(`DELETE FROM public.room_members WHERE room_id = $1 AND player_name = $2`, roomID, playerName)
	if err != nil {
		return err
	}
	return nil
}
