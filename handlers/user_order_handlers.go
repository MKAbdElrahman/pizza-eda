package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pizza/handlers/middleware"
	"pizza/models"
	"pizza/templs"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

func (h *UserHandler) HandleGetUserOrders(w http.ResponseWriter, r *http.Request) {
	userIDPathValue := r.PathValue("userID")

	userID, err := strconv.Atoi(userIDPathValue)
	if err != nil {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "bad user id")
		return
	}

	// user should only view own orders
	if userID != middleware.GetUserIDFromAuthenticatedContext(r.Context()) {
		h.errorHandler.HandleNotAuthorized(w, r, err, "unauthorized")
		return
	}

	orders, err := h.userService.GetOrders(userID)
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	component := templs.OrdersView(templs.NewLayoutData("Orders", r), orders)
	err = component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
}

func (h *UserHandler) HandleGetUserOrder(w http.ResponseWriter, r *http.Request) {
	userIDPathValue := r.PathValue("userID")

	userID, err := strconv.Atoi(userIDPathValue)
	if err != nil {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "bad user id")
		return
	}

	orderID := r.PathValue("orderID")
	order, err := h.userService.GetOrderByID(userID, orderID)
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	flash := h.sessionManager.PopString(r.Context(), "flash")

	layoutData := templs.NewLayoutData("Order", r)
	layoutData.Flash = flash

	component := templs.SingleOrdersView(layoutData, *order)
	err = component.Render(context.Background(), w)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}

}

func (h *UserHandler) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}

	var p models.Pizza

	decoder := schema.NewDecoder()
	err = decoder.Decode(&p, r.PostForm)
	if err != nil {
		http.Error(w, "Failed to decode form", http.StatusBadRequest)
		return
	}
	// Generate order unique ID
	orderID := uuid.New().String()[:8]

	userID, ok := middleware.GetUserIDFromContext(r.Context())

	if !ok {
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
	}

	o := models.PizzaOrder{
		OrderID:   orderID,
		UserID:    userID,
		Pizza:     p,
		Timestamp: time.Now(),
		Status:    "order_placed",
	}

	// Add order to DB
	err = h.ordersStore.InsertOrder(o)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
	// Produce to Kafka topic

	orderAsBytes, err := json.Marshal(o)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
	err = h.publisher.Publish("pizza-ordered", []byte(o.OrderID), orderAsBytes)
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}

	h.sessionManager.Put(r.Context(), "flash", "Order created!")

	http.Redirect(w, r, fmt.Sprintf("/user/%d/orders/%s", o.UserID, o.OrderID), http.StatusSeeOther)

}
