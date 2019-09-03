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

func UrlJoin(base string, paths ...string) (*url.URL, error) {
	u, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	p := append([]string{u.Path}, paths...)
	u.Path = path.Join(p...)
	return u, nil
}

type JokeRetriever struct {
	Random func() (Joke, error)
	WithId func(string) (Joke, error)
}

func HerokuJokesFactory() JokeRetriever {
	type HerokuResponse struct {
		Value struct {
			Author     string
			Id         string
			Joke       string
		}
	}

	getJoke := func(jokeId ...string) (Joke, error) {
		response := new(HerokuResponse)
		err := GetJsonResponse(response, "https://jokes-api.herokuapp.com", append([]string{"/api/joke"}, jokeId...)...)
		if err != nil {
			return *new(Joke), err
		}
		return response.Value, nil
	}

	return JokeRetriever{
		Random: func () (Joke, error) {
			return getJoke()
		},
		WithId: func (jokeId string) (Joke, error) {
			return getJoke(jokeId)
		},
	}
}

func GetJsonResponse(jsonResponse interface{}, urlBase string, urlComponents ...string) (error) {
	fullPath, err := UrlJoin(urlBase, urlComponents...)
	if err != nil {
		return err
	}
	fmt.Println("Getting from:", fullPath.String())
	resp, err := http.Get(fullPath.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(jsonResponse)
	return nil
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
	var handlers = JokeHandlerFactory(HerokuJokesFactory())
	r := mux.NewRouter()
	r.HandleFunc("/jokes", handlers.random)
	r.HandleFunc("/jokes/{id}", handlers.withId)
	http.Handle("/", r)
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", nil))
}
