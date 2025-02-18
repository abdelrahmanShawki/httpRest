package data

import (
	"database/sql"
	"errors"
)

var ErrRecordNotFound = errors.New("record not found")
var ErrEditConflict = errors.New("edit conflict")

type Models struct {
	Movies MovieModel
	Users  UserModel
	Tokens TokenModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{db},
		Users:  UserModel{db},
		Tokens: TokenModel{db},
	}
}
