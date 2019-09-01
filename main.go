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
)

type Joke struct {
	Author     string
	Id         string
	Joke       string
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

type HerokuResponse struct {
	Value struct {
		Author     string
		Id         string
		Joke       string
	}
}

type JokeRetriever struct {
	random func() (Joke, error)
	withId func(string) (Joke, error)
}

var herokuJokes = JokeRetriever{
	random: func () (Joke, error) {
		response := new(HerokuResponse)
		err := GetJsonResponse(response, "https://jokes-api.herokuapp.com", "/api/joke")
		if err != nil {
			return *new(Joke), err
		}
		return response.Value, nil
	},
	withId: func (jokeId string) (Joke, error) {
		response := new(HerokuResponse)
		err := GetJsonResponse(response, "https://jokes-api.herokuapp.com", "/api/joke", jokeId)
		if err != nil {
			return *new(Joke), err
		}
		return response.Value, nil
	},
}

func GetJsonResponse(jsonResponse interface{}, urlBase string, urlComponents ...string) (error) {
	url, urlErr := UrlJoin(urlBase, urlComponents...)
	if urlErr != nil {
		return urlErr
	}
	fmt.Println("Getting from:", url.String())
	return JsonRequest(url.String(), jsonResponse)
}

func PresentJoke(w http.ResponseWriter, joke Joke) {
	templates := template.Must(template.ParseFiles("templates/jokes-template.html"))
	templates.ExecuteTemplate(w, "jokes-template.html", joke)
	Repository()[joke.Id] = joke
}

type Handler struct {
	withId func(http.ResponseWriter, *http.Request)
	random func(http.ResponseWriter, *http.Request)
}

func JokeHandlerFactory(retriever JokeRetriever) Handler {
	return Handler {
		withId: func(w http.ResponseWriter, r *http.Request) {
			jokeId := mux.Vars(r)["id"]
			if response, ok := Repository()[jokeId]; ok {
				PresentJoke(w, response)
			} else {
				joke, err := retriever.withId(jokeId)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				PresentJoke(w, joke)
			}
		},
		random: func(w http.ResponseWriter, r *http.Request) {
			joke, err := retriever.random()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			PresentJoke(w, joke)
		},
	}
}

func main() {
	var handlers = JokeHandlerFactory(herokuJokes)
	r := mux.NewRouter()
	r.HandleFunc("/jokes", handlers.random)
	r.HandleFunc("/jokes/{id}", handlers.withId)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
