package main

func main() {

}

/*


func Newcsv(w http.ResponseWriter, r *http.Request) {
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

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(n); err != nil {
			panic(err)
			res = append(res, n)
		}

		defer db.Close()
	}

-------------------------------------------------------


type Funcionary struct {
	Number int `json:"number" form:"number"`
}

func Count(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	db.QueryRow("SELECT COUNT(*)FROM names")

	dbdata := Funcionary{}
	list := []Funcionary{}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(dbdata); err != nil {
		panic(err)

	}
	list = append(list, dbdata)
	defer db.Close()
}
}
	http.HandleFunc("/count", Count)
	http.HandleFunc("/csv", Newcsv)

-------------------------------------------

//fs := http.FileServer(http.Dir(uploadPath))
	//http.Handle("/files/", http.StripPrefix("/files", fs))
 use /upload for uploading files and /files/{fileName} for downloading")

	//log.Print("Server started on localhost:8080,
*/
