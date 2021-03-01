package main

import (
	"./model"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var db *sql.DB

const (
	host     = "localhost"
	port     = 5432
	user     = ""
	password = ""
	dbname   = "crud"
)

type Env struct {
	db model.Store
}

const (
	MinAge = 18
	MaxAge = 60
	limit  = 2
)

type Page struct {
	Title   string
	Prev    int
	Next    int
	People  []*model.Person
	Page    int
	Search  string
	Order   string
	OrderBy string
}

func Time(start, end, check time.Time) bool {
	return check.After(start) && check.Before(end)
}

func isAgeValid(bday string) bool {
	now := time.Now()
	from := now.AddDate(-MaxAge, 0, 0)
	to := now.AddDate(-MinAge, 0, 0)
	date, _ := time.Parse("2006-01-02", bday)
	isValid := Time(from, to, date)
	return isValid
}

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := model.Init(psqlInfo)
	if err != nil {
		panic(err)
	}

	env := &Env{db}

	http.HandleFunc("/", env.PersonList)
	http.HandleFunc("/create", env.AddPerson)
	http.HandleFunc("/edit", env.EditPerson)
	http.HandleFunc("/search", env.SearchPerson)
	http.HandleFunc("/delete", env.DeletePerson)
	http.ListenAndServe(":8080", nil)
}

func (env *Env) PersonList(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	people, err := env.db.GetPerson()
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	files := []string{
		"templates/list.gohtml",
		"templates/main.gohtml",
		"templates/nav.gohtml",
	}
	templ, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Server error", 500)
		return
	}
	err = templ.Execute(w, people)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Server error", 500)
	}
}

func (env *Env) AddPerson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/add" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case "GET":
		files := []string{
			"templates/add.gohtml",
			"templates/main.gohtml",
			"templates/nav.gohtml",
		}
		templ, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = templ.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		person := model.Person{}
		person.Firstname = r.FormValue("firstName")
		person.Surname = r.FormValue("lastName")
		day := r.FormValue("day")
		month := r.FormValue("month")
		year := r.FormValue("year")
		bday := fmt.Sprintf("%s-%s-%s", year, month, day)

		if isAgeValid(bday) != true {
			message := fmt.Sprintf("The user should be older than %[1]d and younger than %[2]d!", MinAge, MaxAge)
			log.Println(message)
			http.Error(w, message, 500)
			return
		}
		person.Birthday = bday
		person.Gender = r.FormValue("gender")
		person.Email = r.FormValue("email")
		person.Address = r.FormValue("address")

		err := env.db.AddPerson(person)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Such user already exists.", 500)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (env *Env) EditPerson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/edit" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(r.FormValue("ID"))
	switch r.Method {
	case "GET":
		if err != nil {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}
		person, err := env.db.GetById(id)
		if err != nil {
			http.Error(w, http.StatusText(400), http.StatusBadRequest)
			return
		}
		files := []string{
			"templates/edit.gohtml",
			"templates/main.gohtml",
			"templates/nav.gohtml",
		}
		templ, err := template.ParseFiles(files...)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = templ.Execute(w, person)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	case "POST":
		if err := r.ParseForm(); err != nil {
			log.Println(w, "ParseForm() err: %v", err)
			return
		}
		person := model.Person{}
		person.Id, _ = strconv.Atoi(r.FormValue("id"))
		person.Firstname = r.FormValue("firstName")
		person.Surname = r.FormValue("lastName")
		birthDate := strings.Split(r.FormValue("birthDate"), ".")
		day := birthDate[0]
		month := birthDate[1]
		year := birthDate[2]
		person.Birthday = fmt.Sprintf("%s-%s-%s", year, month, day)
		if isAgeValid(person.Birthday) != true {
			msg := fmt.Sprintf("The user should be older than %[1]d and younger than %[2]d!", MinAge, MaxAge)
			log.Println(msg)
			http.Error(w, msg, 500)
			return
		}
		person.Gender = r.FormValue("gender")
		person.Email = r.FormValue("email")
		person.Address = r.FormValue("address")
		err := env.db.UpdatePerson(person)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Such user is already exists.", 500)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func (env *Env) SearchPerson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/search" {
		http.NotFound(w, r)
		return
	}
	var prev, next int
	page, err := strconv.Atoi(r.FormValue("page"))
	if err != nil {
		page = 1
	}
	if page < 1 {
		page = 1
	}
	next = page + 1
	prev = page - 1
	skips := limit * (page - 1)
	search := r.FormValue("q")
	orderBy := "id"
	orderByParam := r.FormValue("orderBy")
	if orderByParam != "" {
		orderBy = orderByParam
	}
	order := "ASC"
	param := r.FormValue("order")
	if param != "" {
		if param == "ASC" {
			order = "DESC"
		} else {
			order = "ASC"
		}
	}
	people, err := env.db.GetByName(search, orderBy, order, limit, skips)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	data := Page{
		Title:   "Search results",
		People:  people,
		OrderBy: orderBy,
		Prev:    prev,
		Next:    next,
		Search:  search,
		Page:    page,
		Order:   order,
	}
	files := []string{
		"templates/search.gohtml",
		"templates/main.gohtml",
		"templates/nav.gohtml",
	}
	tmpl, err := template.ParseFiles(files...)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}
}

func (env *Env) DeletePerson(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delete" {
		http.NotFound(w, r)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	err = env.db.DeletePerson(id)
	if err != nil {
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
