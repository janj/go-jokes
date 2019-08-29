// go run main.go

package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

type Joke struct {
	Author     string
	Id         int
	Joke       string
}

type JokeResponse struct {
	Value Joke
}

type JokeMap map[string]Joke

var (
	r JokeMap
)

func Repository() JokeMap {
	if r == nil {
		r = make(JokeMap)
	}
	return r
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

type JokeSource struct {
	baseUrl string
	pathToJoke string
	pathToRandomJoke string
}

var herokuSource = JokeSource {
	baseUrl: "https://jokes-api.herokuapp.com",
	pathToJoke: "/api/joke",
	pathToRandomJoke: "/api/joke",
}

type JokeHandlerMap struct {
	random func(w http.ResponseWriter, r *http.Request)
	withId func(w http.ResponseWriter, r *http.Request)
}

func GetJokeResponse(urlBase string, urlComponents ...string) (*JokeResponse, error) {
	url, url_err := UrlJoin(urlBase, urlComponents...)
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

func PresentJoke(w http.ResponseWriter, joke Joke) {
	templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
	templates.ExecuteTemplate(w, "jokes-template.html", joke)
	Repository()[strconv.Itoa(joke.Id)] = joke
}

func JokeHandlerFactory(jokeSource JokeSource) JokeHandlerMap {
	return JokeHandlerMap {
		withId: func(w http.ResponseWriter, r *http.Request) {
			joke_id := mux.Vars(r)["id"]
			if response, ok := Repository()[joke_id]; ok {
				PresentJoke(w, response)
			} else {
				jokeResponse, err := GetJokeResponse(jokeSource.baseUrl, jokeSource.pathToJoke, joke_id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				PresentJoke(w, jokeResponse.Value)
			}
		},
		random: func(w http.ResponseWriter, r *http.Request) {
			jokeResponse, err := GetJokeResponse(jokeSource.baseUrl, jokeSource.pathToRandomJoke)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			PresentJoke(w, jokeResponse.Value)
		},
	}
}

func main() {
	var jokeHandlers = JokeHandlerFactory(herokuSource)
	r := mux.NewRouter()
	r.HandleFunc("/jokes", jokeHandlers.random)
	r.HandleFunc("/jokes/{id}", jokeHandlers.withId)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
