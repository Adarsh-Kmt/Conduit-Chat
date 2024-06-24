package util

import (
	"encoding/json"
	"net/http"

	"github.com/Adarsh-Kmt/chatapp/types"
)

type HttpFunc func(http.ResponseWriter, *http.Request) *types.APIError

func MakeHttpHandlerFunc(f HttpFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if APIError := f(w, r); APIError != nil {

			WriteJSON(w, APIError.ErrorStatus, map[string]string{"error": APIError.Error})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, body any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(body)
}
