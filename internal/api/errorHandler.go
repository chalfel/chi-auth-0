package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/chalfel/chi-auth-0/pkg/exceptions"
)

type Response struct {
	StatusCode int
	Data       any
}
type InternalHandler func(w http.ResponseWriter, r *http.Request) (Response, error)

type ErrorResponse struct {
	Message    string `json:"message"`
	Constraint string `json:"constraint"`
}

func errorHandler(cb InternalHandler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp, err := cb(w, r)
		if err != nil {
			var errorBody any
			var status int
			var e *exceptions.Exception

			if errors.As(err, &e) {
				status = http.StatusUnprocessableEntity
				errorBody = err
			} else {
				errorBody = ErrorResponse{
					Message:    "Internal Server Error",
					Constraint: "internal_server_error",
				}

				status = http.StatusInternalServerError
			}

			body, err := json.Marshal(errorBody)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			w.WriteHeader(status)
			w.Write(body)
			return
		}

		body, err := json.Marshal(resp.Data)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(resp.StatusCode)
		w.Write(body)
	})
}
