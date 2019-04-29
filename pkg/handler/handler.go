package handler

import (
	"encoding/json"

	"github.com/kfreiman/imageservice/pkg/service"

	"net/http"
)

// HTTPHandler basically provides access to Service's methods over HTTP.
type HTTPHandler interface {
	ModifyEndpoint() http.HandlerFunc
}

type handler struct {
	service   service.Service
	extractor Extractor
}

type (
	// BadRequestError means that the request is incorrectly formatted
	BadRequestError struct {
		error
	}

	// ValidationError means that the request contains invalid data
	ValidationError struct {
		error
	}
)

// NewHTTPHandler creates implementation of HTTPHandler.
func NewHTTPHandler(service service.Service, extractor Extractor) HTTPHandler {
	return &handler{service, extractor}
}

func errHTTP(w http.ResponseWriter, err error) {
	m := map[string]interface{}{
		"message": err.Error(),
	}
	w.Header().Set("Content-Type", "application/json")

	switch err.(type) {
	case *BadRequestError:
		w.WriteHeader(http.StatusBadRequest)
	case *ValidationError:
		w.WriteHeader(http.StatusUnprocessableEntity)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(m)
}
