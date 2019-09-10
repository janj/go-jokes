package jokesources

import "strconv"

type herokuResponse struct {
	Value struct {
		Id         int
		Joke       string
	}
}

func getHerokuJoke(components ...string) (Joke, error) {
	response := new(herokuResponse)
	err := GetJsonResponse(response, "https://jokes-api.herokuapp.com", append([]string{"/api/joke"}, components...)...)
	if err != nil {
		return *new(Joke), err
	}
	return response.Value, nil
}

var HerokuRetriever = JokeRetriever{
	Random: func () (Joke, error) {
		return getHerokuJoke()
	},
	WithId: func (jokeId int) (Joke, error) {
		return getHerokuJoke(strconv.Itoa(jokeId))
	},
}
