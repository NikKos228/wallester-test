package model

import (
	"database/sql"
	_ "database/sql"
	"log"
	_ "log"
	"time"
)

type Person struct {
	Firstname string
	Surname   string
	Email     string
	Gender    string
	Birthday  string
	Address   string
	Id        int
}

func birthdate(date string) string {
	dt, _ := time.Parse("2006-01-02T15:04:05Z", date)
	bDate := dt.Format("02.01.2006")
	return bDate
}

func (db *database) GetPerson() ([]*Person, error) {
	sqlStatement := `SELECT id, first_name, last_name, 
    birth_date, gender, email, address FROM people;`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	people := make([]*Person, 0)
	for rows.Next() {
		person := new(Person)
		err := rows.Scan(
			&person.Firstname,
			&person.Surname,
			&person.Birthday,
			&person.Gender,
			&person.Email,
			&person.Address)
		if err != nil {
			return nil, err
		}
		person.Birthday = birthdate(person.Birthday)
		people = append(people, person)
	}
	if err = rows.Err(); err != nil {
		return people, err
	}
	return people, nil
}

func (db *database) DeletePerson(id int) error {
	sqlStatement := `DELETE FROM people WHERE id=$1;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return nil
}

func (db *database) AddPerson(person Person) error {
	sqlStatement := `INSERT INTO people (first_name, last_name,
		 birth_date, gender, email, address) VALUES ($1, $2, $3, $4, $5, $6);`
	_, err := db.Exec(sqlStatement,
		person.Firstname,
		person.Surname,
		person.Birthday,
		person.Gender,
		person.Email,
		person.Address)
	if err != nil {
		return err
	}
	return nil
}

func (db *database) UpdatePerson(person Person) error {
	sqlStatement := `UPDATE customers SET first_name=$2, last_name=$3,
	 birth_date=$4, gender=$5, email=$6, address=$7 WHERE id=$1;`
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	_, err = tx.Exec(sqlStatement,
		person.Firstname,
		person.Surname,
		person.Birthday,
		person.Gender,
		person.Email,
		person.Address,
		person.Id)

	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func (db *database) GetByName(name, orderBy, order string, limit, skips int) ([]*Person, error) {
	sqlStatement := `SELECT id, first_name, last_name, 
	birth_date, gender, email, address FROM customers 
	WHERE first_name SIMILAR TO $1 OR last_name SIMILAR TO $1 
	ORDER BY ` + orderBy + ` ` + order + ` LIMIT $2 OFFSET $3;`
	rows, err := db.Query(sqlStatement, name, limit, skips)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	people := make([]*Person, 0)
	for rows.Next() {
		person := new(Person)
		err := rows.Scan(
			&person.Firstname,
			&person.Surname,
			&person.Birthday,
			&person.Gender,
			&person.Email,
			&person.Address,
			&person.Id)
		if err != nil {
			return nil, err
		}
		person.Birthday = birthdate(person.Birthday)
		people = append(people, person)
	}
	if err = rows.Err(); err != nil {
		return people, err
	}
	return people, nil
}

func (db *database) GetById(id int) (Person, error) {
	var person Person
	sqlStatement := `SELECT id, first_name, last_name, birth_date, gender, email, address FROM customers WHERE id=$1;`
	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(
		&person.Id,
		&person.Firstname,
		&person.Surname,
		&person.Birthday,
		&person.Gender,
		&person.Email,
		&person.Address)
	if err != nil && err != sql.ErrNoRows {
		log.Println(err.Error())
		return person, err
	}
	person.Birthday = birthdate(person.Birthday)
	return person, nil
}
