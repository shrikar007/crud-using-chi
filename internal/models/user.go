package models

import (
	"github.com/jmoiron/sqlx"
)

type User struct {
	Id         int    `json:"id" db:"id"`
	FirstName  string `json:"first_name" db:"first_name"`
	MiddleName string `json:"middle_name" db:"middle_name"`
	LastName   string `json:"last_name" db:"last_name"`
	Email      string `json:"email" db:"email"`
	Gender     string `json:"gender" db:"gender"`
	Mobile     string `json:"mobile" db:"mobile"`
	Password   string `json:"password" db:"password"`
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResource struct {
	DB *sqlx.DB
}

func NewUserModel(db *sqlx.DB) (u *UserResource) {
	u = &UserResource{}
	u.DB = db
	return
}
func (u *UserResource) GetUser(login Login) (user User, err error) {
	err = u.DB.Get(&user, "SELECT * FROM users WHERE email=? AND password=?", login.Email, login.Password)
	if err != nil {
		return
	}
	return
}

func (u *UserResource) StoreUser(user User) (err error) {
	_, err = u.DB.Exec(`INSERT INTO users (first_name,middle_name,password,last_name, email,gender,mobile) VALUES (?,?,?,?,?,?,?)`, user.FirstName, user.MiddleName, user.Password, user.LastName, user.Email, user.Gender, user.Mobile)
	if err != nil {
		return
	}
	return
}

func (u *UserResource) GetProfileById(Id string) (user User, err error) {
	err = u.DB.Get(&user, "SELECT * FROM users WHERE id=?", Id)
	if err != nil {
		return
	}
	return
}

func (u *UserResource) UpdateProfile(user User, Id string) (err error) {
	_, err = u.DB.Exec("UPDATE users SET first_name=?,middle_name=?,last_name=?,email=?,gender=?,mobile=?,password=? WHERE id = ?", user.FirstName, user.MiddleName, user.LastName, user.Email, user.Gender, user.Mobile, user.Password, Id)
	if err != nil {
		return
	}
	return
}

func (u *UserResource) GetProfiles() (user []User, err error) {
	err = u.DB.Select(&user, "SELECT * FROM users")
	if err != nil {
		return
	}
	return
}
