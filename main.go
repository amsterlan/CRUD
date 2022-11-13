package main

import (
	csvmanager "CRUD/csv"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"text/template"
)

var tmpl = template.Must(template.ParseGlob("tmpl/*"))

type Names struct {
	Id    int    `json:"id" form:"id"`
	Name  string `json:"name" form:"name"`
	Email string `json:"email" form:"email"`
}

func dbConn() (db *sql.DB) {

	db, err := sql.Open("mysql", "root:alan@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err.Error())
	}
	return db
}

func Index(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	selDB, err := db.Query("SELECT *  FROM names ")
	if err != nil {
		panic(err.Error())
	}

	n := Names{}
	res := []Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Email = email

		res = append(res, n)
	}

	tmpl.ExecuteTemplate(w, "Index", res)

	defer db.Close()
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT id, name, email FROM names WHERE id=?;", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Email = email
	}

	tmpl.ExecuteTemplate(w, "Show", n)

	defer db.Close()

}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT * FROM names WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Names{}

	for selDB.Next() {
		var id int
		var name, email string

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Email = email
	}

	tmpl.ExecuteTemplate(w, "Edit", n)

	// Fecha a conexão com o banco de dados
	defer db.Close()
}

func Insert(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")

		insForm, err := db.Prepare("INSERT INTO names(name, email) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email)

		log.Println("INSERT: Name: " + name + " | E-mail: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Update(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	if r.Method == "POST" {

		name := r.FormValue("name")
		email := r.FormValue("email")
		id := r.FormValue("uid")

		insForm, err := db.Prepare("UPDATE names SET name=?, email=? WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, id)

		log.Println("UPDATE: Name: " + name + " |E-mail: " + email)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM names WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(nId)

	log.Println("DELETE")

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func DownCsv(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	selDB, err := db.Query("SELECT *  FROM names ")
	if err != nil {
		panic(err.Error())
	}

	n := Names{}
	res := []Names{}

	for selDB.Next() {
		var id int
		var name, email string

		err = selDB.Scan(&id, &name, &email)
		if err != nil {
			panic(err.Error())
		}

		n.Id = id
		n.Name = name
		n.Email = email

		res = append(res, n)

		u := fmt.Sprintf("%#v,%#v,%#v", id, name, email)

		list := [][]string{
			{u},
		}

		byteData, err := csvmanager.WriteAll(list)

		if err != nil {
			log.Fatalln(err)
		}
		w.Write([]byte(byteData))
	}
}

type Regs struct {
	Id int `json:"p" form:"p"`
}

func Registers(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	regs := Regs{}
	reg := []Regs{}

	rows, err := db.Query("SELECT COUNT(*)FROM test.names")
	if err != nil {
		panic(err.Error())
	}

	var lines int
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&lines); err != nil {
			log.Fatal(err)
			fmt.Printf("%s", lines)

			regs.Id = id

			reg = append(reg, regs)
		}
		fmt.Println(lines)
		tmpl.ExecuteTemplate(w, "Count", lines)
	}
}

func main() {
	http.HandleFunc(`/lines`, Registers)
	http.HandleFunc(`/csv`, DownCsv)
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
