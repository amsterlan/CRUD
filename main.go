package main

import (
	csvmanager "CRUD/csv"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

var tmpl = template.Must(template.ParseGlob("tmpl/*"))

func dbConn() (db *sql.DB) {

	db, err := sql.Open("mysql", "root:alan@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {

	employees := []Employee{}

	db := dbConn()
	result, err := db.Query("SELECT *  FROM employees ")
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}

		registeredEmployee := Employee{id, name, email, salary}
		employees = append(employees, registeredEmployee)

	}

	listIndexPage := IndexPage{len(employees), employees}

	tmpl.ExecuteTemplate(w, "Index", listIndexPage)

	defer result.Close()
}

type Employee struct {
	Id     int
	Name   string
	Email  string
	Salary float64
}

type IndexPage struct {
	Count     int
	Employees []Employee
}

func Show(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}
	db := dbConn()

	getId := r.URL.Query().Get("id")

	result, err := db.Query("SELECT * FROM employees WHERE id=?", getId)
	if err != nil {
		panic(err.Error())
	}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}

		employee = Employee{id, name, email, salary}

	}

	tmpl.ExecuteTemplate(w, "Show", employee)

	defer db.Close()

}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	getId := r.URL.Query().Get("id")

	result, err := db.Query("SELECT * FROM employees WHERE id=?", getId)
	if err != nil {
		panic(err.Error())
	}

	employee := Employee{}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		// Faz o Scan do SELECT
		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}
		employee = Employee{id, name, email, salary}

	}

	tmpl.ExecuteTemplate(w, "Edit", employee)

	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")
		salary := r.FormValue("salary")

		insForm, err := db.Prepare("INSERT INTO employees (name, email, salary) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, salary)

		log.Println("INSERT: Name: " + name + " | E-mail: " + email + "| Salary: " + salary)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")
		salary := r.FormValue("salary")
		id := r.FormValue("id")

		insForm, err := db.Prepare("UPDATE employees SET name=?, email=?, salary=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, salary, id)

		log.Println("UPDATE: Name: " + name + " |E-mail: " + email + "|Salary:" + salary)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	getId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM employees WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(getId)

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func DownloadCsv(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	result, err := db.Query("SELECT *  FROM employees ")
	if err != nil {
		panic(err.Error())
	}

	list := [][]string{}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}

		employee := []string{strconv.Itoa(id), name, email, strconv.FormatFloat(salary, 'f', 2, 64)}

		list = append(list, employee)
	}

	byteData, err := csvmanager.WriteAll(list)

	if err != nil {
		log.Fatalln(err)
	}
	w.Write([]byte(byteData))
}

func main() {
	http.HandleFunc(`/csv`, DownloadCsv)
	http.HandleFunc("/", Index)
	http.HandleFunc("/show", Show)
	http.HandleFunc("/new", New)
	http.HandleFunc("/edit", Edit)
	http.HandleFunc("/insert", Insert)
	http.HandleFunc("/update", Update)
	http.HandleFunc("/delete", Delete)

	http.ListenAndServe(":9000", nil)
	log.Println("Server started on: http://localhost:9000")
}
