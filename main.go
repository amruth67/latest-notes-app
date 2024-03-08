package main

import (
	"fmt"
	"mongodb/controllers"
	"mongodb/database"
	"net/http"

	"github.com/gorilla/mux"
)

//var client *mongo.Client

func main() {
	database.DBinstance()
	r := mux.NewRouter()
	r.HandleFunc("/Signup", controllers.SignupHandler).Methods("POST")
	r.HandleFunc("/Login", controllers.LoginHandler).Methods("POST")
	r.HandleFunc("/CreateNotes", controllers.CreateNotes).Methods("POST")
	r.HandleFunc("/GetNotes", controllers.GetNotes).Methods("GET")
	http.ListenAndServe(":8000", r)
	fmt.Println("starting the server at 8000")
}
