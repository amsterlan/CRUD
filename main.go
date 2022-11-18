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

	sliceEmployee := []Employee{}

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

		listFuncionary := Employee{id, name, email, salary}
		sliceEmployee = append(sliceEmployee, listFuncionary)

	}

	count := len(sliceEmployee)
	listIndexPage := IndexPage{count, sliceEmployee}

	tmpl.ExecuteTemplate(w, "Index", listIndexPage)

	defer result.Close()
}

type Employee struct {
	id     int
	name   string
	email  string
	salary float64
}

type IndexPage struct {
	count      int
	funcionary []Employee
}

func Show(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	result, err := db.Query("SELECT id, name, email, salary FROM employees WHERE id=?;", nId)
	if err != nil {
		panic(err.Error())
	}

	listFuncionary := Employee{}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}

		listFuncionary.id = id
		listFuncionary.name = name
		listFuncionary.email = email
		listFuncionary.salary = salary
	}

	tmpl.ExecuteTemplate(w, "Show", listFuncionary)

	defer db.Close()

}

func New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "New", nil)
}

func Edit(w http.ResponseWriter, r *http.Request) {
	db := dbConn()

	nId := r.URL.Query().Get("id")

	result, err := db.Query("SELECT * FROM employees WHERE id=?", nId)
	if err != nil {
		panic(err.Error())
	}

	listEmployee := Employee{}

	for result.Next() {
		var id int
		var name, email string
		var salary float64

		// Faz o Scan do SELECT
		err = result.Scan(&id, &name, &email, &salary)
		if err != nil {
			panic(err.Error())
		}

		listEmployee.id = id
		listEmployee.name = name
		listEmployee.email = email
		listEmployee.salary = salary
	}

	tmpl.ExecuteTemplate(w, "Edit", listEmployee)

	// Fecha a conexão com o banco de dados
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
		id := r.FormValue("uid")
		salary := r.FormValue("salary")

		insForm, err := db.Prepare("UPDATE employees SET name=?, email=?, salary=?, WHERE id=?")
		if err != nil {
			panic(err.Error())
		}

		insForm.Exec(name, email, id, salary)

		log.Println("UPDATE: Name: " + name + " |E-mail: " + email + "|Salario:" + salary)
	}

	defer db.Close()

	http.Redirect(w, r, "/", 301)
}

func Delete(w http.ResponseWriter, r *http.Request) {

	db := dbConn()

	nId := r.URL.Query().Get("id")

	delForm, err := db.Prepare("DELETE FROM employees WHERE id=?")
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
	resultado, err := db.Query("SELECT *  FROM employees ")
	if err != nil {
		panic(err.Error())
	}

	list := [][]string{}

	for resultado.Next() {
		var id int
		var name, email string
		var salary float64

		err = resultado.Scan(&id, &name, &email, &salary)
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
