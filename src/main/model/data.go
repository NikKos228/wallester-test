package model

import (
	"database/sql"
	_ "github.com/jackc/pgx"
)

type Store interface {
	AddPerson(person Person) error
	GetPerson() ([]*Person, error)
	GetByName(name, orderBy, order string, limit, skips int) ([]*Person, error)
	DeletePerson(id int) error
	UpdatePerson(person Person) error
	GetById(id int) (Person, error)
}

type database struct {
	*sql.DB
}

func Init(dsname string) (*database, error) {
	db, err := sql.Open("postgres", dsname)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &database{db}, nil
}
