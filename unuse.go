package main

/*
import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path"
)



 func Down(w http.ResponseWriter, r *http.Request) {

	records := [][]string{
		{"first_name", "last_name", "occupation"},
		{"John", "Doe", "gardener"},
		{"Lucy", "Smith", "teacher"},
		{"Brian", "Bethamy", "programmer", "%d"},
	}

	f, err := os.Create("users.csv")
	defer f.Close()

	if err != nil {

		log.Fatalln("failed to open file", err)
	}

	wr := csv.NewWriter(f)

	err = wr.WriteAll(records)

	if err != nil {
		log.Fatal(err)
	}
}

const maxUploadSize = 2 * 1024 * 1024 // 2 mb
const uploadPath = "./tmpl"

func uploadFileHandler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			t, _ := template.ParseFiles("upload.tmpl")
			t.Execute(w, nil)
			return
		}
		if err := r.ParseMultipartForm(maxUploadSize); err != nil {
			fmt.Printf("Could not parse multipart form: %v\n", err)
			renderError(w, "CANT_PARSE_FORM", http.StatusInternalServerError)
			return
		}

		// parse and validate file and post parameters
		file, fileHeader, err := r.FormFile("uploadFile")
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}
		defer file.Close()
		// Get and print out file size
		fileSize := fileHeader.Size
		fmt.Printf("File size (bytes): %v\n", fileSize)
		// validate file size
		if fileSize > maxUploadSize {
			renderError(w, "FILE_TOO_BIG", http.StatusBadRequest)
			return
		}
		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			renderError(w, "INVALID_FILE", http.StatusBadRequest)
			return
		}

		// check file type, detectcontenttype only needs the first 512 bytes
		detectedFileType := http.DetectContentType(fileBytes)
		switch detectedFileType {
		case "image/jpeg", "image/jpg":
		case "image/gif", "image/png":
		case "application/pdf":
		case "file/csv":
		case "file/txt":
			break
		default:
			renderError(w, "INVALID_FILE_TYPE", http.StatusBadRequest)
			return
		}
		fileName := randToken(12)
		fileEndings, err := mime.ExtensionsByType(detectedFileType)
		if err != nil {
			renderError(w, "CANT_READ_FILE_TYPE", http.StatusInternalServerError)
			return
		}
		newPath := filepath.Join(uploadPath, fileName+fileEndings[0])
		fmt.Printf("FileType: %s, File: %s\n", detectedFileType, newPath)

		// write file
		newFile, err := os.Create(newPath)
		if err != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		defer newFile.Close() // idempotent, okay to call twice
		if _, err := newFile.Write(fileBytes); err != nil || newFile.Close() != nil {
			renderError(w, "CANT_WRITE_FILE", http.StatusInternalServerError)
			return
		}
		w.Write([]byte("SUCCESS"))
	})
}

func renderError(w http.ResponseWriter, message string, statusCode int) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(message))
}

func randToken(len int) string {
	b := make([]byte, len)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
 /---------------------------------------------

func DownloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return
}

//-----------------------------------------------------------------------


func main() {

}

type Count struct {
	NumberCount int
}

func Download(url, filename string) (err error) {
	fmt.Println("Downloading ", url, " to ", filename)

	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return
}

func donwload() {
	pUrl := flag.String("url", "", "/kk")
	flag.Parse()
	url := *pUrl
	if url == "" {
		fmt.Fprintf(os.Stderr, "Error: empty URL!\n")
		return
	}

	filename := path.Base(url)
	fmt.Println("Checking if " + filename + " exists ...")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := download(url, filename)
		if err != nil {
			panic(err)
		}
		fmt.Println(filename + " saved!")
	} else {
		fmt.Println(filename + " already exists!")
	}
}

func download(url string, filename string) interface{} {

}

type ShoppingRecord struct {
	Vegetable string
	Fruit     string
}

func createShoppingList(data [][]string) []ShoppingRecord {
	var shoppingList []ShoppingRecord
	for i, line := range data {
		if i > 0 { // omit header line
			var rec ShoppingRecord
			for j, field := range line {
				if j == 0 {
					rec.Vegetable = field
				} else if j == 1 {
					rec.Fruit = field
				}
			}
			shoppingList = append(shoppingList, rec)
		}
	}
	return shoppingList
}

func List() {
	s := gin.New()
	s.GET("/newcsv", func(c *gin.Context) {
		resp := map[string]Names{}
		c.JSON(http.StatusOK, resp)
	})
	log.Fatal(http.ListenAndServe(":9000", s))
}

func Read(w http.ResponseWriter, r *http.Request) {
	// open file
	f, err := os.Open("data.csv")
	if err != nil {
		log.Fatal(err)
	}

	// remember to close the file at the end of the program
	defer f.Close()

	// read csv values using csv.Reader
	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// convert records to array of structs
	shoppingList := createShoppingList(data)

	// print the array
	fmt.Printf("%+v\n", shoppingList)
}

/* func Csv(w http.ResponseWriter, r *http.Request) {
	/*	db := dbConn()

		names := make([]Names, 0)
		rows, err := db.Query("SELECT *  FROM names")
		defer rows.Close()

		if err != nil {
			return
		}

		for rows.Next() {
			var n Names
			rows.Scan(&n.Id, &n.Name, &n.Email)
			names = append(names, n)
		}
		if err = rows.Err(); err != nil {
			return
		}

	b := &bytes.Buffer{}
	wr := json2csv.NewCSVWriter(b)
	j, _ := os.ReadFile("users.json")
	var x []map[string]interface{}

	// unMarshall json
	err := json.Unmarshal(j, &x)
	if err != nil {
		log.Fatal(err)
	}

	// convert json to CSV
	csv, err := json2csv.JSON2CSV(x)
	if err != nil {
		log.Fatal(err)
	}
	err = wr.WriteCSV(csv)
	if err != nil {
		log.Fatal(err)
	}
	wr.Flush()
	got := b.String()

	//Following line prints CSV
	println(got)

	// create file and append if you want
	createFileAppendText("output.csv", got)
}

func createFileAppendText(filename string, text string) {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.WriteString(text); err != nil {
		panic(err)
	}
}


//------------------------------------------------------------------------------------

func (n *Names) GetList(names []Names, err error) {
	db := dbConn()

	names = make([]Names, 0)
	rows, err := db.Query("SELECT *  FROM names")
	defer rows.Close()

	if err != nil {
		return
	}

	return
}

func GetNames(c *gin.Context) {
	var n Names
	names, err := n.GetList()
	if err != nil {
		log.Fatalln()
	}
	c.JSON(http.StatusOK, gin.H{
		"names": names,
	})
	log.Printf("JSON", n)
}



//-----------------------------------------------

db := dbConn()

rows, err := db.Query("SELECT COUNT(*)FROM test.names")
if err != nil {
	panic(err.Error())
}

m := Count{}
res := []Count{}

res = append(res, m)
tmpl.ExecuteTemplate(w, "Index", res)


}

*/
