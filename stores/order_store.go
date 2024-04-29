package stores

import (
	"database/sql"
	"pizza/models"
	"time"
)

type orderStore struct {
	db *sql.DB
}

func NewOrderStore(db *sql.DB) *orderStore {
	return &orderStore{
		db: db,
	}
}

func (os *orderStore) InsertOrder(o models.PizzaOrder) error {
	stmt, err := os.db.Prepare("INSERT INTO orders (order_id, user_id, sauce, cheese, main_topping, extra_topping, status, timestamp) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(o.OrderID, o.UserID, o.Sauce, o.Cheese, o.MainTopping, o.ExtraTopping, "order_placed", time.Now())
	if err != nil {
		return err
	}

	return nil
}
