package data

import (
	"database/sql"
	"errors"
	"rxrz/greenlight/internal/validator"
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	//it's necessary to start the key names with capital letters so that they can be exported
	//and the compiler can see them.
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"` //hide it from the endpoint response
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"` //hide it from the endpoint if empty
	Runtime   Runtime   `json:"runtime,omitempty"`
	Genres    []string  `json:"genres,omitempty"` // if you do not want to use keyname: ",omitempty"
	Version   int32     `json:"version"`
}

func ValidateMovie(v *validator.Validator, m *Movie) {
	v.Check(m.Title != "", "title", "must be provided")
	v.Check(len(m.Title) <= 500, "title", "must not be more than 500 bytes long")
	v.Check(m.Year != 0, "year", "must be provided")
	v.Check(m.Year >= 1888, "year", "must be greater than 1888")
	v.Check(m.Year <= int32(time.Now().Year()), "year", "must not be in the future")
	v.Check(m.Runtime != 0, "runtime", "must be provided")
	v.Check(m.Runtime > 0, "runtime", "must be a positive integer")
	v.Check(m.Genres != nil, "genres", "must be provided")
	v.Check(len(m.Genres) >= 1, "genres", "must contain at least 1 genre")
	v.Check(len(m.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.Unique(m.Genres), "genres", "must not contain duplicate values")
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
	return m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

// Add a placeholder method for fetching a specific record from the movies
// table.
func (m MovieModel) Get(id int64) (*Movie, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1
	`
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	return &movie, nil

}

// Add a placeholder method for updating a specific record in the movies
// table.
func (m MovieModel) Update(movie *Movie) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the movies
// table.
func (m MovieModel) Delete(id int64) error {
	return nil
}

/*
Note that the string directive will only work on struct fields which
have int*, uint*, float* or bool types. For any other type of struct
field it will have no effect.
*/

// func (m Movie) MarshalJSON() ([]byte, error) {
// 	var runtime string
// 	if m.Runtime != 0 {
// 		runtime = fmt.Sprintf("%d mins", m.Runtime)
// 	}
// 	aux := struct {
// 		ID      int64    `json:"id"`
// 		Title   string   `json:"title"`
// 		Year    int32    `json:"year,omitempty"`
// 		Runtime string   `json:"runtime,omitempty"` // This is a string.
// 		Genres  []string `json:"genres,omitempty"`
// 		Version int32    `json:"version"`
// 	}{
// 		// Set the values for the anonymous struct.
// 		ID:      m.ID,
// 		Title:   m.Title,
// 		Year:    m.Year,
// 		Runtime: runtime,
// 		Genres:  m.Genres,
// 		Version: m.Version,
// 	}
// 	return json.Marshal(aux)
// }
