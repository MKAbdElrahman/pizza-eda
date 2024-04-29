package models

type PizzaMenu struct {
	Sauces        []string `json:"sauces"`
	Cheeses       []string `json:"cheeses"`
	MainToppings  []string `json:"main_toppings"`
	ExtraToppings []string `json:"extra_toppings"`
}
