package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	_ "time"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var database *sql.DB

func GetPage(w http.ResponseWriter, pageGUID string) (Page, error) {
	thisPage := Page{}
	fmt.Println(White + "Getting " + Green + pageGUID + White + " page..." + Reset)

	// Get page from database by guid
	query := "SELECT page_guid, page_title, page_content, page_date FROM pages WHERE page_guid = $1"
	err := database.QueryRow(query, pageGUID).Scan(&thisPage.GUID, &thisPage.Title, &thisPage.RawContent, &thisPage.Date)

	if err != nil {
		http.Error(w, http.StatusText(404), http.StatusNotFound)
		log.Println(Red + "Couldn't get page " + pageGUID + Reset)
		log.Println(Red + err.Error() + Reset)
		return thisPage, err
	}

	fmt.Println(White + "Page_guid " + Green + pageGUID + White + " accessed" + Reset)

	// Get raw content
	thisPage.Content = template.HTML(thisPage.RawContent)

	return thisPage, err
}

func APIPage(w http.ResponseWriter, r *http.Request) {
	// Get GUID variable
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	// Get page by GUID
	thisPage, err := GetPage(w, pageGUID)
	if err != nil {
		return
	}

	// To json
	_, err = json.Marshal(thisPage)

	if err != nil {
		http.Error(w, Red+err.Error()+Reset, http.StatusInternalServerError)
		return
	}

	// Return page JSON
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, thisPage)
}

func APICommentPost(w http.ResponseWriter, r *http.Request) {

	var commentAdded bool

	err := r.ParseForm()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	guid := r.FormValue("guid")
	name := r.FormValue("name")
	email := r.FormValue("email")
	comments := r.FormValue("comments")

	// Insert comment and return the id of the newly inserted comment
	query := "INSERT INTO comments (comment_name, comment_email, comment_text, page_id) VALUES ($1, $2, $3, $4) RETURNING id"
	var id int
	err = database.QueryRow(query, name, email, comments, guid).Scan(&id)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	commentAdded = true

	// Prepare the JSON response
	resp := JSONResponse{Fields: make(map[string]string)}
	resp.Fields["id"] = strconv.Itoa(id) // Use Itoa to convert int to string
	resp.Fields["added"] = strconv.FormatBool(commentAdded)

	// Send the JSON response
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintln(w, string(jsonResp))
}

func ServePage(w http.ResponseWriter, r *http.Request) {
	// Get GUID
	vars := mux.Vars(r)
	pageGUID := vars["guid"]

	// Get page
	thisPage, err := GetPage(w, pageGUID)
	if err != nil {
		return
	}

	// Parse to html
	t, _ := template.ParseFiles("templates/blog.html")
	t.Execute(w, thisPage)
}

func main() {
	fmt.Println(Magenta + "Starting server..." + Reset)

	// Router
	routes := mux.NewRouter()
	routes.HandleFunc("/api/pages", APIPage).Methods("GET").Schemes("https")
	routes.HandleFunc("/api/pages/{guid:[0-9a-zA\\-]+}", APIPage).
		Methods("GET").
		Schemes("https")
	routes.HandleFunc("/api/comments", APICommentPost)
	routes.HandleFunc("/page/{guid:[0-9a-zA\\-]+}", ServePage)

	http.Handle("/", routes)

	fmt.Println(Magenta + "Connecting to database..." + Reset)

	// Connect to database
	dbConn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		DBHost, DBPort, DBUser, DBPass, DBDbase)
	db, err := sql.Open("postgres", dbConn)
	if err != nil {
		log.Println(Red + "Couldn't connect to database" + Reset)
		log.Println(Red + err.Error() + Reset)
		return
	}

	// Assign connected db to global database variable
	database = db
	fmt.Println(Cyan + "Successfully connected to database" + Reset)

	cert, err := tls.LoadX509KeyPair("cert.pem", "key.pem")
	if err != nil {
		log.Fatal(err)
	}

	// Run server
	fmt.Println("Server is running on " + Green + "http://localhost" + PORT + Reset)

	tlsConf := tls.Config{Certificates: []tls.Certificate{cert}}
	listener, err := tls.Listen("tcp", PORT, &tlsConf)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", PORT, err)
	}
	log.Fatal(http.Serve(listener, nil))
}
