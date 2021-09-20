package transport

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type Post  struct {
	Title string `json:"title"`
	Id int `json:"id"`
	PostDate string `json:"postDate"`
}
const(
	host = "ec2-63-33-14-215.eu-west-1.compute.amazonaws.com"
	port = "5432"
	user = "efqfutmrzzunvf"
	password = "d9fbc191dd683c84603b840e4cd07b2ae93dcf4af8959a4230f2a5e0ccf1b1a0"
	dbname = "d4obgiuuk72amd"
)
func openConnection() *sql.DB {
	psqlConnStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=require", host, port,user,password,dbname)
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err!= nil{
		panic(err)
	}
	return db
}
func Router() http.Handler {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v1").Subrouter()
	s.HandleFunc("/posts", GetPostsHandler).Methods("GET", "OPTIONS")
	s.HandleFunc("/posts", AddPostHandler).Methods("POST", "OPTIONS")
	s.HandleFunc("/posts/{id}", GetPostHandler).Methods("GET", "OPTIONS")
	s.HandleFunc("/posts/{id}", EditPostHandler).Methods("PUT", "OPTIONS")
	s.HandleFunc("/posts/{id}", DeletePostHandler).Methods("DELETE", "OPTIONS")
	return logMiddleware(r)
}

var decoder = schema.NewDecoder()


func GetPostsHandler( w http.ResponseWriter, r *http.Request){
	db := openConnection()

	rows, err := db.Query(`SELECT * FROM posts`)
	if err!=nil{
		log.Fatal(err)
	}
	var posts []Post

	for rows.Next(){
		var post Post
		rows.Scan(&post.Title,&post.Id, &post.PostDate)
		posts = append(posts, post)
	}

	postsBytes, _ := json.MarshalIndent(posts, "", "\t")

	w.Header().Set("Content-Type", "application/json")
	w.Write(postsBytes)

	defer rows.Close()
	defer db.Close()
}
func GetPostHandler(w http.ResponseWriter, r *http.Request){
	db:=openConnection()
	params := mux.Vars(r)
	row := db.QueryRow(`SELECT * FROM posts WHERE posts.id=$1`,params["id"])
	var post Post
	row.Scan(&post.Title,&post.Id, &post.PostDate)
	w.Header().Set("Content-Type", "application/json")
	postBytes, _ := json.MarshalIndent(post, "", "\t")
	w.Write(postBytes)
	defer db.Close()
}
func EditPostHandler(w http.ResponseWriter, r *http.Request){
	db:=openConnection()
	var post Post
	err:= json.NewDecoder(r.Body).Decode(&post)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	params := mux.Vars(r)
	_, err = db.Exec(`UPDATE posts SET title = $1, postdate = $2 WHERE id = $3`,post.Title, post.PostDate,params["id"] )
	if err!=nil {
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
func DeletePostHandler(w http.ResponseWriter, r *http.Request){
	db:=openConnection()
	params := mux.Vars(r)
	_, err := db.Exec(`DELETE FROM posts WHERE id = $1`, params["id"] )
	if err!=nil{
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
func AddPostHandler(w http.ResponseWriter, r *http.Request){
	db:=openConnection()
	var post Post
	err := json.NewDecoder(r.Body).Decode(&post)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err = db.Exec(`INSERT INTO posts (title, postdate) VALUES ($1, $2)`, post.Title, post.PostDate)
	if err !=nil {
		w.WriteHeader(http.StatusBadRequest)
		panic(err)
	}
	w.WriteHeader(http.StatusOK)
	defer db.Close()
}
func logMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method":     r.Method,
			"url":        r.URL,
			"remoteAddr": r.RemoteAddr,
			"userAgent":  r.UserAgent(),
		}).Info("got a new request")
		h.ServeHTTP(w, r)
	})


}
