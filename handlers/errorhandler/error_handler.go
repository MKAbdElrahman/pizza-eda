package errorhandler

import (
	"log/slog"
	"net/http"
)

type CentralErrorHandler struct {
	logger *slog.Logger
}

func NewCentralErrorHandler(logger *slog.Logger) *CentralErrorHandler {
	return &CentralErrorHandler{
		logger: logger,
	}
}
func (h *CentralErrorHandler) HandleInternalServerError(w http.ResponseWriter, r *http.Request, err error, msg string) {
	method := r.Method
	uri := r.URL.RequestURI()
	h.logger.Error(err.Error(), slog.String("method", method), slog.String("uri", uri))
	http.Error(w, msg, http.StatusInternalServerError)
}

func (h *CentralErrorHandler) HandleBadRequestFromClient(w http.ResponseWriter, r *http.Request, err error, msg string) {
	http.Error(w, msg, http.StatusBadRequest)
}

func (h *CentralErrorHandler) HandleResourceNotFound(w http.ResponseWriter, r *http.Request, err error, msg string) {
	http.Error(w, msg, http.StatusNotFound)
}

func (h *CentralErrorHandler) HandleNotAuthorized(w http.ResponseWriter, r *http.Request, err error, msg string) {
	http.Error(w, msg, http.StatusUnauthorized)
}
