package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mongodb/database"
	"mongodb/middleware"
	"mongodb/models"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "dev")

func getHash(password []byte) string {
	hash, err := bcrypt.GenerateFromPassword(password, 14)
	if err != nil {
		log.Println(err)
	}
	return string(hash)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Typye", "application/json")
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)

	//hashing the password
	user.Password = getHash([]byte(user.Password))
	if err != nil {
		fmt.Println("error while decoding json body")
	}
	// validate := validator.New()
	// validationErr := validate.Struct(user)
	// if validationErr != nil {
	// 	w.WriteHeader(http.StatusUnauthorized)
	// 	fmt.Fprint(w, "validation error")
	// 	return

	// }
	//inserting into database
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	//checking if email exists
	count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	if count > 0 {
		http.Error(w, "email  already exists", http.StatusInternalServerError)
		return
	}
	//checking if phone number already exists
	count, err = userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	if count > 0 {
		http.Error(w, " phone already exists", http.StatusInternalServerError)
		return
	}

	count, err = userCollection.CountDocuments(ctx, bson.M{"username": user.UserName})
	defer cancel()
	if err != nil {
		log.Panic(err)
		return
	}
	if count > 0 {
		http.Error(w, "username already exists", http.StatusInternalServerError)
		return
	}

	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("insert document with id:", result.InsertedID)
	json.NewEncoder(w).Encode(result.InsertedID)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var user models.User
	var dbuser models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "error while decoding login details", http.StatusInternalServerError)
		return
	}

	var ctx, _ = context.WithTimeout(context.Background(), 100*time.Second)
	dbfinderr := userCollection.FindOne(ctx, bson.M{"username": user.UserName}).Decode(&dbuser)
	if dbfinderr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"` + err.Error() + `"}`))
		return
	}

	fmt.Println(user.Password, dbuser.Password)

	passErr := bcrypt.CompareHashAndPassword([]byte(dbuser.Password), []byte(user.Password))
	if passErr != nil {
		fmt.Printf("error druing comparing password: %v", passErr)
		return
	}

	fmt.Println(user.Password, dbuser.Password)
	fmt.Println(passErr)

	//defer cancel()

	if dbuser.UserName == "" {
		fmt.Println("username is not found")

	}

	tokenString, _ := middleware.GenerateToken(user.UserName)
	json.NewEncoder(w).Encode(tokenString)

}

var notes []models.Notes

func CreateNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var note models.Notes
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		http.Error(w, "error while decoding", http.StatusInternalServerError)
	}
	notes = append(notes, note)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	result, err := userCollection.InsertOne(ctx, note)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("insert document with id:", result.InsertedID)
	json.NewEncoder(w).Encode(result.InsertedID)

}

func GetNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application-json")
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, "missing authorozation header")
		return
	}
	tokenString = tokenString[len("Bearer "):]

	err := middleware.VerifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "invalid token")
		return
	}

	fmt.Fprint(w, "welcome to protected")
	json.NewEncoder(w).Encode(notes)
}
