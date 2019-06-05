package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"path"
)

type JokeResponse struct {
	Value struct {
		Author     string
		Categories []string
		Id         int
		Joke       string
	}
}

func JsonRequest(url string, v interface{}) error {
	resp, err := http.Get(url)
	if err == nil {
		defer resp.Body.Close()
		json.NewDecoder(resp.Body).Decode(v)
	}
	return err
}

func UrlJoin(base string, paths ...string) (*url.URL, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	p := append([]string{u.Path}, paths...)
	u.Path = path.Join(p...)
	return u, nil
}

func GetJokeResponse(joke_id string) (*JokeResponse, error) {
	url, url_err := UrlJoin("https://jokes-api.herokuapp.com", "/api/joke", joke_id)
	if url_err != nil {
		return nil, url_err
	}
	jokeResponse := new(JokeResponse)
	fmt.Println(url.String())
	json_err := JsonRequest(url.String(), jokeResponse)
	if json_err != nil {
		return nil, json_err
	}
	return jokeResponse, nil
}

func JokeHandler(w http.ResponseWriter, r *http.Request) {
	joke_id := mux.Vars(r)["id"]
	jokeResponse, err := GetJokeResponse(joke_id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
	templates.ExecuteTemplate(w, "jokes-template.html", jokeResponse.Value)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/jokes", JokeHandler)
	r.HandleFunc("/jokes/{id}", JokeHandler)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
