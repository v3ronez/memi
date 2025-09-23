package handler

import (
	"net/http"

	"github.com/v3ronez/memi/view/auth"
)

func HandleLoginIndex(w http.ResponseWriter, r *http.Request) error {
	return auth.Login().Render(r.Context(), w)
}

