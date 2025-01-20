package data

import "time"

type Movie struct {
	ID        int64     `json:",omitempty"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title,omitempty"`
	Year      int32     `json:",omitempty"`
	Runtime   int32     `json:",omitempty"` //custom type
	Genres    []string  `json:"category,omitempty"`
	Version   int32
	// time the movie information is updated
}
