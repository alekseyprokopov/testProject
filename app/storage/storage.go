package storage

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"testProject/app/response"
	"testProject/configs"
)

type Storage struct {
	db *sql.DB
}

func New(path string) (*Storage, error) {
	db, err := sql.Open("postgres", configs.SqlPath)

	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	return &Storage{
		db: db,
	}, nil
}
func (s *Storage) Save(user *response.UserData) error {
	personId, err := s.savePerson(user)
	if err != nil {
		return err
	}

	if err := s.saveLocation(user, personId); err != nil {
		return err
	}

	return nil
}

func (s *Storage) savePerson(user *response.UserData) (int64, error) {
	var lastId int64 = 0
	q := `INSERT INTO person(gender, title_name, first_name, last_name) VALUES ($1, $2, $3, $4) RETURNING id`

	//result, err := s.db.Exec(q, user.Gender, user.Name.Title, user.Name.First, user.Name.Last)
	if err := s.db.QueryRow(q, user.Gender, user.Name.Title, user.Name.First, user.Name.Last).Scan(&lastId); err != nil {
		return 0, fmt.Errorf("can't insert data in person table: %w", err)
	}

	return lastId, nil
}
func (s *Storage) saveLocation(user *response.UserData, personId int64) error {

	q := `INSERT INTO location(street_number,street_name,city,country,postcode,coordinates_latitude,coordinates_longitude, person_id)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	_, err := s.db.Exec(
		q,
		user.Location.Street.Number,
		user.Location.Street.Name,
		user.Location.City,
		user.Location.Country,
		user.Location.Postcode,
		user.Location.Coordinates.Latitude,
		user.Location.Coordinates.Longitude,
		personId,
	)
	if err != nil {
		log.Printf("ERRROR %v", err)

		return fmt.Errorf("can't insert data in location table: %w", err)
	}

	return nil
}
func (s *Storage) savePersonLocation(personId, locationId int64) error {
	q := `INSERT INTO person_location(person_id, location_id) VALUES (?, ?)`

	_, err := s.db.Exec(q, personId, locationId)
	if err != nil {
		return fmt.Errorf("can't insert data in person_location table: %w", err)
	}
	return nil
}

func (s *Storage) Close() {
	s.db.Close()
}
func (s *Storage) Init() error {

	personQuery := `CREATE TABLE IF NOT EXISTS person(
    	id SERIAL PRIMARY KEY,
    	gender VARCHAR(10),
    	title_name VARCHAR (10),
    	first_name VARCHAR (20),
    	last_name VARCHAR (20)
    )`

	locationQuery := `CREATE TABLE IF NOT EXISTS location(
    	id SERIAL PRIMARY KEY,
    	street_number INTEGER,
    	street_name VARCHAR(20),
    	city VARCHAR(20),
    	country VARCHAR(20),
    	postcode VARCHAR(20),
    	coordinates_latitude VARCHAR (20),
    	coordinates_longitude VARCHAR (20),
    	person_id INTEGER, 
        FOREIGN KEY (person_id) REFERENCES person(id)
    )`

	_, err := s.db.Exec(personQuery)
	if err != nil {
		return fmt.Errorf("can't create person table: %w", err)
	}

	_, err = s.db.Exec(locationQuery)
	if err != nil {
		return fmt.Errorf("can't create location table: %w", err)
	}

	return nil
}
