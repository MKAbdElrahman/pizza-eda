package handlers

import (
	"context"
	"fmt"
	"net/http"
	"pizza/models"
	"pizza/templs"
)

func HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Ok")
}

func HandleViewHome(w http.ResponseWriter, r *http.Request) {

	component := templs.HomeView(r)
	err := component.Render(context.Background(), w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func HandleViewMenu(w http.ResponseWriter, r *http.Request) {

	menu := models.PizzaMenu{
		Sauces:        []string{"Tomato", "Organic", "tomato", "Pesto", "Marinara", "Buffalo", "Hummus"},
		Cheeses:       []string{"Mozzarella", "Provolone", "Cheddar", "Ricotta", "Gouda", "Gruyere"},
		MainToppings:  []string{"Pepperoni", "Sausage", "Chicken", "Pork", "Minced meat", "Vegan meat"},
		ExtraToppings: []string{"Mushroom", "Onion", "Egg", "Ham", "Green pepper", "Fresh garlic"},
	}

	component := templs.MenuView(r, menu)
	err := component.Render(context.Background(), w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
