package models

import (
	"time"
)

type User struct {
	ID             int
	Name           string
	Email          string
	Address        string
	Phone          string
	HashedPassword []byte
	Created        time.Time
}

type UserSignupParams struct {
	Name     string `schema:"username"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
	Address  string `schema:"address"`
	Phone    string `schema:"phone"`
}



type UserLoginParams struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}
