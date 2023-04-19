package entity

import "golang.org/x/crypto/bcrypt"

type User struct {
	ID       string
	Login    string
	Password string
	UserBalance
}

type UserBalance struct {
	Balance float32
	Spent   float32
}

func (u *User) IsValidPassword() bool {
	return u.Password != "" && len(u.Password) > 4
}

func (u *User) IsValidLogin() bool {
	return u.Password != "" && len(u.Login) > 4
}

func (u *User) IsEqual(other User) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(other.Password))
	if err != nil {
		return false
	}

	return u.Login == other.Login
}
