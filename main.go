package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Movie struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Isbn     string             `json:"isbn"`
	Title    string             `json:"title"`
	Genre    string             `json:"genre"`
	Year     int                `json:"year"`
	Director *Director          `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI("mongodb+srv://riddhikakundu:<db-password>@cluster0.0vn1uh4.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("movieDB").Collection("movies")
}

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	query := r.URL.Query().Get("search")
	filter := bson.M{}
	if query != "" {
		filter = bson.M{"title": bson.M{"$regex": query, "$options": "i"}}
	}
	var movies []Movie
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		var movie Movie
		cursor.Decode(&movie)
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var movie Movie
	err := collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&movie)
	if err != nil {
		http.Error(w, "Movie not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(movie)
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	res, err := collection.InsertOne(context.TODO(), movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	movie.ID = res.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	var movie Movie
	_ = json.NewDecoder(r.Body).Decode(&movie)
	update := bson.M{
		"$set": movie,
	}
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	movie.ID = id
	json.NewEncoder(w).Encode(movie)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id, _ := primitive.ObjectIDFromHex(params["id"])
	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(bson.M{"message": "Movie deleted"})
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	fmt.Println("Server starting at port 8000...")
	log.Fatal(http.ListenAndServe(":8000", r))
}
