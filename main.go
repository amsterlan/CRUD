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

	n := Funcionario{}
	SlcFuncionario := []Funcionario{}
	IndxPg := IndexPage{}

	db := dbConn()
	resultado, err := db.Query("SELECT *  FROM funcionarios ")
	if err != nil {
		panic(err.Error())
	}

	for resultado.Next() {
		var id int
		var name, email string
		var salario float64

		err = resultado.Scan(&id, &name, &email, &salario)
		if err != nil {
			panic(err.Error())
		}

		n.id = id
		n.name = name
		n.email = email
		n.salario = salario

		SlcFuncionario = append(SlcFuncionario, n)

	}

	count := len(SlcFuncionario)
	IndxPg.funcionarios = SlcFuncionario
	IndxPg.count = count
	tmpl.ExecuteTemplate(w, "Index", IndxPg)
	defer db.Close()
}

type Funcionario struct {
	id      int
	name    string
	email   string
	salario float64
}

type IndexPage struct {
	count        int
	funcionarios []Funcionario
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	selDB, err := db.Query("SELECT id, name, email, salario FROM funcionarios WHERE id=?;", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Funcionario{}

	for selDB.Next() {
		var id int
		var name, email string
		var salario float64

		err = selDB.Scan(&id, &name, &email, &salario)
		if err != nil {
			panic(err.Error())
		}

		n.id = id
		n.name = name
		n.email = email
		n.salario = salario
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

	selDB, err := db.Query("SELECT * FROM funcionarios WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	n := Funcionario{}

	for selDB.Next() {
		var id int
		var name, email string
		var salario float64

		// Faz o Scan do SELECT
		err = selDB.Scan(&id, &name, &email, &salario)
		if err != nil {
			panic(err.Error())
		}

		n.id = id
		n.name = name
		n.email = email
		n.salario = salario
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
		salario := r.FormValue("salario")

		insForm, err := db.Prepare("INSERT INTO funcionarios (name, email, salario) VALUES(?,?,?)")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, salario)

		log.Println("INSERT: Name: " + name + " | E-mail: " + email + "| Salario: " + salario)
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
		salario := r.FormValue("salario")

		insForm, err := db.Prepare("UPDATE funcionarios SET name=?, email=?, salario=?, WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, id, salario)

		log.Println("UPDATE: Name: " + name + " |E-mail: " + email + "|Salario:" + salario)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM funcionarios WHERE id=?")
	if err != nil {
		panic(err.Error())
	}

	delForm.Exec(nId)

	log.Println("DELETE")

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

/*
func Indexadd(r *IndexPage) {
	r.count += len(listFuncionarios)
}*/

func DownCsv(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	resultado, err := db.Query("SELECT *  FROM funcionarios ")
	if err != nil {
		panic(err.Error())
	}

	list := [][]string{}

	for resultado.Next() {
		var id int
		var name, email string
		var salario float64

		err = resultado.Scan(&id, &name, &email, &salario)
		if err != nil {
			panic(err.Error())
		}

		funcionario := []string{strconv.Itoa(id), name, email, strconv.FormatFloat(salario, 'f', 2, 64)}

		list = append(list, funcionario)
	}

	byteData, err := csvmanager.WriteAll(list)

	if err != nil {
		log.Fatalln(err)
	}
	w.Write([]byte(byteData))
}

/*
	  func Registers(w http.ResponseWriter, r *http.Request) {

		db := dbConn()

		rows, err := db.Query("SELECT COUNT(*)FROM funcionarios")
		if err != nil {
			panic(err.Error())
		}
		regs := Regs{}
		reg := []Regs{}
		var lines int
		defer rows.Close()

		for rows.Next() {
			var registro int
			if err := rows.Scan(&lines); err != nil {
				log.Fatal(err)
				fmt.Printf("%s", lines)

				regs.Registers = registro

				reg = append(reg, regs)
			}
			fmt.Println(lines)
			tmpl.ExecuteTemplate(w, "index", regs)
		}

}
*/

func main() {

	//var lenghList = len(listFuncionarios)
	//http.HandleFunc(`/lines`, Registers)
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
