package services

import (
	"errors"
	"pizza/models"

	"golang.org/x/crypto/bcrypt"
)

type UserStore interface {
	Insert(params models.UserSignupParams) error
	GetUserByEmail(email string) (*models.User, error)
}

type OrderStore interface {
	GetOrders(userID int) ([]models.PizzaOrder, error)
	GetOrderByID(userID int, orderID string) (*models.PizzaOrder, error)
}

type userService struct {
	userStore  UserStore
	orderStore OrderStore
}

func NewUserService(userStore UserStore, orderStore OrderStore) *userService {

	return &userService{
		userStore:  userStore,
		orderStore: orderStore,
	}
}

func (m *userService) InsertUser(p models.UserSignupParams) error {

	return m.userStore.Insert(p)

}
func (m *userService) Authenticate(p models.UserLoginParams) (int, error) {
	user, err := m.userStore.GetUserByEmail(p.Email)
	if err != nil {
		if errors.Is(err, models.ErrUserNotFound) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(p.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, models.ErrInvalidCredentials
		} else {
			return 0, err
		}
	}
	return user.ID, nil
}

func (m *userService) Exists(id int) (bool, error) {
	return false, nil
}

func (m *userService) GetOrders(userID int) ([]models.PizzaOrder, error) {
	return m.orderStore.GetOrders(userID)
}

func (m *userService) GetOrderByID(userID int, orderID string) (*models.PizzaOrder, error) {
	return m.orderStore.GetOrderByID(userID, orderID)
}
