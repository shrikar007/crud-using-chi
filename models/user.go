package models

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

type LoginResponse struct {
	Token string `json:"token"`
}



