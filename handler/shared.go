package handler

import (
	"log/slog"
	"net/http"
)

func HandleTempl(h func(http.ResponseWriter, *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			slog.Error("Error to render page", "err", err)
			return
		}
	}
}
