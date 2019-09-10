package jokesources

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
)

type Joke struct {
	Id         int
	Joke       string
}

type JokeRetriever struct {
	Random func() (Joke, error)
	WithId func(int) (Joke, error)
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
	//dump, err := httputil.DumpResponse(resp, true)
	//if err == nil {
	//	fmt.Println("Response dump:", string(dump))
	//}
	json.NewDecoder(resp.Body).Decode(jsonResponse)
	return nil
}
