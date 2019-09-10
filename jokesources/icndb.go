package jokesources

import "strconv"

type icndbResponse struct {
	Value struct {
		Id         int `json:"id"`
		Joke       string `json:"joke"`
	} `json:"value"`
}

func getIcndbJoke(components ...string) (Joke, error) {
	response := new(icndbResponse)
	err := GetJsonResponse(response, "https://api.icndb.com", append([]string{"/jokes"}, components...)...)
	if err != nil {
		return *new(Joke), err
	}
	return Joke(response.Value), nil
}

var IcndbRetriever = JokeRetriever{
	Random: func () (Joke, error) {
		return getIcndbJoke("random")
	},
	WithId: func (jokeId int) (Joke, error) {
		return getIcndbJoke(strconv.Itoa(jokeId))
	},
}
