package data

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"httpRest/internal/validator"
	"time"
)

type Movie struct {
	ID        int64     `json:",omitempty"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title,omitempty"`
	Year      int32     `json:",omitempty"`
	Runtime   Runtime   `json:",omitempty"` //custom type
	Genres    []string  `json:"category,omitempty"`
	Version   int32
}

func ValidateMovie(v *validator.Validator, movie *Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= int32(time.Now().Year()), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(movie.Genres), "genres", "must not contain duplicate values")
}

type MovieModel struct {
	DB *sql.DB
}

func (m MovieModel) Insert(movie *Movie) error {
	query := `
			INSERT INTO movies (title, year, runtime, genres)
			VALUES ($1, $2, $3, $4)
			RETURNING id, created_at, version
				`
	args := []interface{}{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	//setting timeout for inserting
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
	if err != nil {
		return err
	}
	return nil
}

func (m MovieModel) Get(id int64) (*Movie, error) {

	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie Movie
	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrRecordNotFound
		} else {
			return nil, err
		}
	}
	return &movie, nil
}

func (m MovieModel) Update(movie *Movie) error {

	query := `UPDATE movies
			SET title = $1, 
				year = $2, 
				runtime = $3, 
				genres = $4, 
				version = version + 1
			WHERE id = $5 AND version = $6
			RETURNING version;`

	args := []interface{}{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	return nil
}

func (m MovieModel) Delete(id int64) error {
	query := `DELETE FROM movies
				WHERE id = $1`
	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowAff, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowAff == 0 {
		return ErrRecordNotFound
	}
	return nil
}

type MockMovieModel struct{}

func (m MockMovieModel) Insert(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Get(id int64) (*Movie, error) {
	return nil, nil
}

func (m MockMovieModel) Update(movie *Movie) error {
	return nil
}

func (m MockMovieModel) Delete(id int64) error {
	return nil
}
