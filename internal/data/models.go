package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflit")
)

type Models struct {
	Users UserModel
}

func NewModel(conn *sql.DB) Models {
	return Models{
		Users: UserModel{DB: conn},
	}
}
