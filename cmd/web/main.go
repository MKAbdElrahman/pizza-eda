package main

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"pizza/handlers/errorhandler"
	"pizza/pubsub"
	"pizza/services"
	"pizza/stores"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/alice"

	"pizza/handlers"
	"pizza/handlers/middleware"

	_ "github.com/go-sql-driver/mysql"
)

var sessionManager = scs.New()

func main() {

	// LOGGER
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// DATABASE
	dsn := os.Getenv("USER_DB_DSN_FOR_SERVER")
	dbConn, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()
	logger.Info("connected to database")

	// SESSION MANAGER

	sessionManager.Store = mysqlstore.New(dbConn)
	sessionManager.Lifetime = 12 * time.Hour

	// MUX
	mux := http.NewServeMux()

	errorHandler := errorhandler.NewCentralErrorHandler(logger)
	userStore := stores.NewUserStore(dbConn)
	orderStore := stores.NewOrderStore(dbConn)

	producer, err := pubsub.NewConfluentProducer("localhost:9092")

	if err != nil {
		panic(err)
	}

	dynamic := alice.New(sessionManager.LoadAndSave, middleware.SetAuthenticatedUserInContext(sessionManager, userStore, errorHandler))

	protected := dynamic.Append(middleware.RequireLoggingInFirst(sessionManager))

	standard := alice.New(chimiddleware.Logger)

	// static files
	fileServer := http.FileServer(http.Dir("./static_assets/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))
	userService := services.NewUserService(userStore, orderStore)
	userHandler := handlers.NewUserHandler(logger, userService, orderStore, producer, sessionManager, logger, errorHandler)

	mux.Handle("GET /{$}", dynamic.ThenFunc(handlers.HandleViewHome))
	mux.Handle("GET /menu", dynamic.ThenFunc(handlers.HandleViewMenu))

	mux.HandleFunc("GET /health-check", handlers.HandleHealthCheck)

	mux.Handle("GET /user/login", dynamic.ThenFunc(userHandler.HandleGetUserLoginForm))
	mux.Handle("POST /user/login", dynamic.ThenFunc(userHandler.HandlePostedLogin))
	mux.Handle("GET /user/signup", dynamic.ThenFunc(userHandler.HandleGetUserSignupForm))
	mux.Handle("POST /user/signup", dynamic.ThenFunc(userHandler.HandlePostedSignup))

	mux.Handle("POST /order", protected.ThenFunc(userHandler.HandleCreateOrder))

	mux.Handle("GET /user/{userID}/orders", protected.ThenFunc(userHandler.HandleGetUserOrders))
	mux.Handle("GET /user/{userID}/orders/{orderID}", protected.ThenFunc(userHandler.HandleGetUserOrder))

	mux.Handle("/user/logout", protected.ThenFunc(userHandler.HandlePostedLogout))

	// SERVER
	addr := os.Getenv("SERVER_ADDR")
	logger.Info("starting server", slog.String("addr", addr))
	serv := http.Server{
		Handler:      standard.Then(mux),
		Addr:         addr,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Fatal(serv.ListenAndServe())
}
