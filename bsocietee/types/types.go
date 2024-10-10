package types

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
