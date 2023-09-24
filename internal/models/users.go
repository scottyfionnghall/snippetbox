package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	InserStmt *sql.Stmt
	AuthStmt  *sql.Stmt
	ExistStmt *sql.Stmt
	DB        *sql.DB
}

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

func NewUserModel(db *sql.DB) (*UserModel, error) {
	insertSmt, err := db.Prepare(`INSERT INTO users (name, email, hashed_password, created)
	VALUES(?,?,?,UTC_TIMESTAMP())`)
	if err != nil {
		return nil, err
	}
	authStmt, err := db.Prepare("SELECT id, hashed_password FROM users WHERE email = ?")
	if err != nil {
		return nil, err
	}
	existStmt, err := db.Prepare("SELECT EXISTS(SELECT true FROM users WHERE id = ?)")
	if err != nil {
		return nil, err
	}
	return &UserModel{InserStmt: insertSmt, AuthStmt: authStmt, ExistStmt: existStmt, DB: db}, nil
}

func (u *UserModel) CloseAll() error {
	err := u.InserStmt.Close()
	if err != nil {
		return err
	}
	err = u.AuthStmt.Close()
	if err != nil {
		return err
	}
	return nil
}

// Method to insert a new user into database
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	// Use the Exec() method
	_, err = m.InserStmt.Exec(name, email, string(hashedPassword))
	if err != nil {
		// If this returns an error, we use the errors.As() to check
		// whether the error has the type *mysql.MySQLError. If it does,
		// the error releates to user entering already existing email
		// (which has code 1062), the return an ErrDuplicateEmail
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

// Method to authinticate the user. Verifes whether a user exists with the
// provided email address and password. Returns relevant user ID.
func (m *UserModel) Authenticate(email, password string) (int, error) {
	var id int
	var hashedPassword []byte
	err := m.AuthStmt.QueryRow(email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	// Check whether the hashed password and plain-text password provided match.
	// If the don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return id, nil
}

// Method to check if the user with a specific ID exists
func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	err := m.ExistStmt.QueryRow(id).Scan(&exists)
	return exists, err
}
