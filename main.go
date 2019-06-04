package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Joke struct {
	Text string
}

type JokeResponse struct {
	Value struct {
		Author     string
		Categories []string
		Id         int
		Joke       string
	}
}

func JokesHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("https://jokes-api.herokuapp.com/api/joke")
	if err == nil {
		defer resp.Body.Close()
		jokeResponse := new(JokeResponse)
		json.NewDecoder(resp.Body).Decode(jokeResponse)
		joke := Joke{jokeResponse.Value.Joke}
		templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
		templates.ExecuteTemplate(w, "jokes-template.html", joke)
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/jokes", JokesHandler)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
