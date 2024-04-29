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

func (os *orderStore) GetOrders(userID int) ([]models.PizzaOrder, error) {
	rows, err := os.db.Query("SELECT order_id, user_id, sauce, cheese, main_topping, extra_topping, status, timestamp FROM orders WHERE user_id = ? ORDER BY timestamp DESC", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make([]models.PizzaOrder, 0)
	for rows.Next() {
		var order models.PizzaOrder
		err := rows.Scan(&order.OrderID, &order.UserID, &order.Sauce, &order.Cheese, &order.MainTopping, &order.ExtraTopping, &order.Status, &order.Timestamp)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (os *orderStore) GetOrderByID(userID int, orderID string) (*models.PizzaOrder, error) {
	var order models.PizzaOrder
	err := os.db.QueryRow("SELECT order_id, user_id, sauce, cheese, main_topping, extra_topping, status, timestamp FROM orders WHERE user_id = ? AND order_id = ?", userID, orderID).
		Scan(&order.OrderID, &order.UserID, &order.Sauce, &order.Cheese, &order.MainTopping, &order.ExtraTopping, &order.Status, &order.Timestamp)
	if err != nil {
		return nil, err
	}
	return &order, nil
}
