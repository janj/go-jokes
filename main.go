package main

import (
	"fmt"
	"html/template"
	"net/http"
)

type Joke struct {
	Setup     string
	Punchline string
}

func JokesHandler(w http.ResponseWriter, r *http.Request) {
	joke := Joke{"What rhymes with orange?", "No it doesn't!"}
	templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
	if err := templates.ExecuteTemplate(w, "jokes-template.html", joke); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/jokes", JokesHandler)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
