package data

import "time"

type Movie struct {
	ID        int64
	CreatedAt time.Time `json:"-"`
	Title     string
	Year      int32    `json:",omitempty"`
	Runtime   int32    `json:",omitempty"` //custom type
	Genres    []string `json:"category"`
	Version   int32
	// time the movie information is updated
}
