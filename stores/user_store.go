package stores

import (
	"database/sql"
	"errors"
	"strings"

	"pizza/models"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
)

type userStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *userStore {
	return &userStore{
		db: db,
	}
}

func (m *userStore) Insert(params models.UserSignupParams) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password,address,phone,created)
VALUES(?, ?, ?,?,?,UTC_TIMESTAMP())`

	_, err = m.db.Exec(stmt, params.Name, params.Email, string(hashedPassword), params.Address, params.Phone)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users_uc_email") {
				return models.ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (m *userStore) GetUserByEmail(email string) (*models.User, error) {
	var user models.User

	query := "SELECT id, name, email, hashed_password, address, phone, created FROM users WHERE email = ?"
	err := m.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.HashedPassword, &user.Address, &user.Phone, &user.Created)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, models.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (m *userStore) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.db.QueryRow(stmt, id).Scan(&exists)
	return exists, err
}
