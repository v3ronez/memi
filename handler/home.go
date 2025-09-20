package handler

import (
	"net/http"

	homeView "github.com/v3ronez/memi/view/home"
)

func HomeHandleIndex(w http.ResponseWriter, r *http.Request) error {
	return homeView.Index().Render(r.Context(), w)
}
