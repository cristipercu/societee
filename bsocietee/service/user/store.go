package user

import (
	"database/sql"
	"fmt"

	"github.com/cristipercu/societee/bsocietee/types"
)

type Store struct {
  db *sql.DB
}

func NewStore (db *sql.DB) *Store {
  return &Store{
    db: db,
  }
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
  user := new(types.User)

  err := s.db.QueryRow(`SELECT id, username, email, password, created_on, modified_on
    FROM public.users WHERE email = $1`, email).Scan(
      &user.ID, &user.Username, &user.Email, &user.Password, &user.CreatedOn, &user.ModifiedOn, 
    )

  if err != nil {
    return nil, err
  }

  if user.ID == 0 {
    return nil, fmt.Errorf("user not found")
  }

  return user, nil
}


func (s *Store) GetUserByID (id int) (*types.UserProfile, error) {
  user := new(types.UserProfile)

  err := s.db.QueryRow(`SELECT id, username, email, created_on, modified_on
    FROM public.users WHERE id = $1`, id).Scan(
      &user.ID, &user.Username, &user.Email, &user.CreatedOn, &user.ModifiedOn, 
    )

  if err != nil {
    return nil, err
  }

  if user.ID == 0 {
    return nil, fmt.Errorf("user not found")
  }

  return user, nil
}


func (s *Store) CreateUser(user types.User) error {
  _, err := s.db.Exec(`INSERT INTO public.users(username, email, password) VALUES ($1, $2, $3)`, user.Username, user.Email, user.Password)

  if err != nil {
    return err
  }
  return nil
}


