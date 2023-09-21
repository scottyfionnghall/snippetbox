package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID int
	Name string
	Email string
	HashedPassword []byte
	Created time.Time
}

type UserModel struct{
	DB *sql.DB
}

// Method to insert a new user into database
func (m *UserModel) Insert(name, email, password string) error{
	return nil
}
// Method to authinticate the user. Verifes whether a user exists with the
// provided email address and password. Returns relevant user ID.
func (m *UserModel) Authenticate(email, password string) (int, error){
	return 0, nil
}
// Method to check if the user with a specific ID exists
func (m *UserModel) Exists (id int) (bool, error){
	return false, nil
}