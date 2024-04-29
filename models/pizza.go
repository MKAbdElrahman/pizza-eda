package models

import "time"

type Pizza struct {
	Sauce        string `schema:"sauce" json:"sauce"`
	Cheese       string `schema:"cheese" json:"cheese"`
	MainTopping  string `schema:"main_topping" json:"main_topping"`
	ExtraTopping string `schema:"extra_topping" json:"extra_topping"`
}

type PizzaOrder struct {
	OrderID   string `json:"order_id"`
	UserID    int    `json:"user_id"`
	Pizza     `json:"pizza"`
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}
