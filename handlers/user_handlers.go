package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"pizza/handlers/errorhandler"
	"pizza/models"
	"pizza/pubsub"
	"pizza/templs"

	"github.com/alexedwards/scs/v2"
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
	publisher      pubsub.Producer
	logger         *slog.Logger
}

func NewUserHandler(log *slog.Logger, userService UserService, orderStore OrderStore, publisher pubsub.Producer, sessionManager *scs.SessionManager, logger *slog.Logger, errorHandler *errorhandler.CentralErrorHandler) *UserHandler {
	return &UserHandler{
		userService:    userService,
		errorHandler:   errorHandler,
		sessionManager: sessionManager,
		ordersStore:    orderStore,
		publisher:      publisher,
		logger:         log,
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
	component := templs.LoginView(templs.NewLayoutData("Login", r))
	err := component.Render(context.Background(), w)

	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
}

func (h *UserHandler) HandleGetUserSignupForm(w http.ResponseWriter, r *http.Request) {
	component := templs.SignupView(templs.NewLayoutData("Signup", r))
	err := component.Render(context.Background(), w)

	if err != nil {
		h.errorHandler.HandleInternalServerError(w, r, err, "internal server error")
		return
	}
}
