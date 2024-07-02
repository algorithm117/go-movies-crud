package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// using structs and slices to store the data instead of external database
type Movie struct {
	ID       string    `json:"id"`
	Isbn     string    `json:"isbn"`
	Title    string    `json:"title"`
	Director *Director `json:"director"`
}

type Director struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	// we want to send the data in a json format so we make sure the writer has that header specified to json
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}

func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// access the id needed to delete specific movie
	params := mux.Vars(r)
	fmt.Println(params)
	for index, item := range movies {
		if item.ID == params["id"] {
			// movies[:index] = all movies up until index of item to be deleted ( so deleted item will not be in this movies slice ), and then movies[index+1:] is all the movies after the movie we want to delete
			// appeneded to the movies from the first parameter movies[:index]
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}

	// return updated list of movies
	json.NewEncoder(w).Encode(movies)
}

func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for _, movie := range movies {
		if movie.ID == params["id"] {
			json.NewEncoder(w).Encode(movie)
			return
		}
	}
}

func createMovie(w http.ResponseWriter, r *http.Request) {
	// to create a movie you need an object of Movie and Director
	// you will need to populate the fields of the two objects with values being sent over by the Request
	// you will need to append the movie to your global list of movies
	// return the new list of movies in json format as a response

	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	// decode the body of the request for the data needed to create the new movie. Store the decoded data into movie variable we created
	_ = json.NewDecoder(r.Body).Decode(&movie)
	// pick a number between 1 and 1000000 and also convert the integer to a string value
	movie.ID = strconv.Itoa(rand.Intn(1000000))
	movies = append(movies, movie)

	// let the user know the movie has been created
	json.NewEncoder(w).Encode(movie)
}

func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var movie Movie
	for index, movie := range movies {
		if movie.ID == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			break
		}
	}
	json.NewDecoder(r.Body).Decode(&movie)
	movie.ID = params["id"]
	movies = append(movies, movie)
	json.NewEncoder(w).Encode(movie)
}

func main() {
	r := mux.NewRouter()

	// initial movies added for testing application
	// & is the address and the * is used to retrieve the value stored at the address
	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
	movies = append(movies, Movie{ID: "2", Isbn: "340940", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

	fmt.Printf("Starting server at port 8080\n")
	log.Fatal(http.ListenAndServe(":8000", r))

}
