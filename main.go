// go run main.go

package main

import (
	"encoding/json"
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
}

type Handler struct {
	withId func(http.ResponseWriter, *http.Request)
	random func(http.ResponseWriter, *http.Request)
	randomApi func(http.ResponseWriter, *http.Request)
	withIdApi func(http.ResponseWriter, *http.Request)
}

func JokeHandlerFactory(retriever jokesources.JokeRetriever) Handler {
	jokeForIdParam := func(r *http.Request) (jokesources.Joke, error) {
		jokeId, err := strconv.Atoi(mux.Vars(r)["id"])
		if err != nil {
			return *new(jokesources.Joke), err
		}
		if response, ok := Repository()[jokeId]; ok {
			return response, nil
		} else {
			joke, err := retriever.WithId(jokeId)
			if err != nil {
				return *new(jokesources.Joke), err
			}
			Repository()[joke.Id] = joke
			return joke, nil
		}
	}

	return Handler {
		withId: func(w http.ResponseWriter, r *http.Request) {
			joke, err := jokeForIdParam(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			PresentJoke(w, joke)
		},
		random: func(w http.ResponseWriter, r *http.Request) {
			templates := template.Must(template.ParseFiles("templates/random-template.html"))
			templates.ExecuteTemplate(w, "random-template.html", nil)
		},
		randomApi: func(w http.ResponseWriter, r *http.Request) {
			joke, err := retriever.Random()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				Repository()[joke.Id] = joke
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode([]jokesources.Joke{joke})
			}
		},
		withIdApi: func(w http.ResponseWriter, r *http.Request) {
			joke, err := jokeForIdParam(r)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]jokesources.Joke{joke})
		},
	}
}

func main() {
	var handlers = JokeHandlerFactory(jokesources.IcndbRetriever)
	r := mux.NewRouter()
	r.PathPrefix("/templates/js").Handler(http.StripPrefix("/templates/js", http.FileServer(http.Dir("."+"/templates/js"))))
	r.HandleFunc("/jokes", handlers.random)
	r.HandleFunc("/jokes/{id}", handlers.withId)
	r.HandleFunc("/api/jokes", handlers.randomApi)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
