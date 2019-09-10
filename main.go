// go run main.go

package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"strconv"
)

import "./jokesources"

type JokeMap map[int]jokesources.Joke

var (
	r JokeMap
)

func Repository() JokeMap {
	if r == nil {
		r = make(JokeMap)
	}
	return r
}

func PresentJoke(w http.ResponseWriter, joke jokesources.Joke) {
	templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
	templates.ExecuteTemplate(w, "jokes-template.html", joke)
	Repository()[joke.Id] = joke
}

type Handler struct {
	withId func(http.ResponseWriter, *http.Request)
	random func(http.ResponseWriter, *http.Request)
}

func JokeHandlerFactory(retriever jokesources.JokeRetriever) Handler {
	return Handler {
		withId: func(w http.ResponseWriter, r *http.Request) {
			jokeId, err := strconv.Atoi(mux.Vars(r)["id"])
			if err != nil {
				return // should present something to the user here
			}
			if response, ok := Repository()[jokeId]; ok {
				PresentJoke(w, response)
			} else {
				joke, err := retriever.WithId(jokeId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				PresentJoke(w, joke)
			}
		},
		random: func(w http.ResponseWriter, r *http.Request) {
			joke, err := retriever.Random()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			PresentJoke(w, joke)
		},
	}
}

func main() {
	var handlers = JokeHandlerFactory(jokesources.IcndbRetriever)
	r := mux.NewRouter()
	r.HandleFunc("/jokes", handlers.random)
	r.HandleFunc("/jokes/{id}", handlers.withId)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
