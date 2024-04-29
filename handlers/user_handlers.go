package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"pizza/handlers/errorhandler"
	"pizza/handlers/middleware"
	"pizza/models"
	"pizza/templs"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/gorilla/schema"
)

type UserService interface {
	InsertUser(models.UserSignupParams) error
	Authenticate(models.UserLoginParams) (int, error)
	GetOrders(userID int) ([]models.PizzaOrder, error)
	GetOrderByID(userID int, orderID string) (*models.PizzaOrder, error)
}

type Publisher interface {
	Publish(topic string, key, value []byte) error
}

type OrderStore interface {
	InsertOrder(o models.PizzaOrder) error
}

type UserHandler struct {
	userService    UserService
	errorHandler   *errorhandler.CentralErrorHandler
	sessionManager *scs.SessionManager
	ordersStore    OrderStore
	publisher      Publisher
}

func NewUserHandler(userService UserService, orderStore OrderStore, publisher Publisher, sessionManager *scs.SessionManager, logger *slog.Logger, errorHandler *errorhandler.CentralErrorHandler) *UserHandler {
	return &UserHandler{
		userService:    userService,
		errorHandler:   errorHandler,
		sessionManager: sessionManager,
		ordersStore:    orderStore,
		publisher:      publisher,
	}
}

func (h *UserHandler) HandlePostedSignup(w http.ResponseWriter, r *http.Request) {
	var p models.UserSignupParams

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&p, r.PostForm)
	if err != nil {
		http.Error(w, "Failed to decode form", http.StatusBadRequest)
		return
	}

	err = h.userService.InsertUser(p)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			h.errorHandler.HandleBadRequestFromClient(w, r, err, "Email already registered")
		} else {
			h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		}
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (h *UserHandler) HandlePostedLogin(w http.ResponseWriter, r *http.Request) {
	var p models.UserLoginParams

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	decoder := schema.NewDecoder()
	err = decoder.Decode(&p, r.PostForm)
	if err != nil {
		http.Error(w, "Failed to decode form", http.StatusBadRequest)
		return
	}

	id, err := h.userService.Authenticate(p)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			h.errorHandler.HandleBadRequestFromClient(w, r, err, "invalid credentials")

		} else {
			h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		}
		return
	}

	err = h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}

	h.sessionManager.Put(r.Context(), "authenticatedUserID", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) HandlePostedLogout(w http.ResponseWriter, r *http.Request) {

	err := h.sessionManager.RenewToken(r.Context())
	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}

	h.sessionManager.Remove(r.Context(), "authenticatedUserID")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *UserHandler) HandleGetUserLoginForm(w http.ResponseWriter, r *http.Request) {
	component := templs.LoginView(r)
	err := component.Render(context.Background(), w)

	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
}

func (h *UserHandler) HandleGetUserSignupForm(w http.ResponseWriter, r *http.Request) {
	component := templs.SignupView(r)
	err := component.Render(context.Background(), w)

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

	fmt.Println(o)

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
	// Redirect
	// http.Redirect(w, r, "/get_order?order_id="+orderID, http.StatusSeeOther)

}

func (h *UserHandler) HandleGetUserOrders(w http.ResponseWriter, r *http.Request) {
	userIDPathValue := r.PathValue("userID")

	userID, err := strconv.Atoi(userIDPathValue)
	if err != nil {
		h.errorHandler.HandleBadRequestFromClient(w, r, err, "bad user id")
		return
	}

	orders, err := h.userService.GetOrders(userID)
	if err != nil {
		http.Error(w, "Failed to get orders", http.StatusInternalServerError)
		return
	}

	component := templs.OrdersView(r, orders)
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

	response, err := json.Marshal(order)
	if err != nil {
		http.Error(w, "Failed to marshal orders", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
