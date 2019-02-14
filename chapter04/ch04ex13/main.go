// Ch04ex13 queries the IMDb API at voku.xyz for movie titles and fetches the
// poster for the movie supplied as a command line argument.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apiBaseURL = "https://voku.xyz/imdb/api"

type movie struct {
	Name      string
	PosterURL string `json:"poster"`
	Error     string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s MOVIE\n", os.Args[0])
		os.Exit(1)
	}

	movie, err := getMovie(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get movie information: %s\n", err)
		os.Exit(1)
	}

	err = getPoster(movie)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not get poster: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Poster stored in %s.jpg\n", movie.Name)
}

func getMovie(name string) (*movie, error) {
	url := fmt.Sprintf("%s?q=%s", apiBaseURL, url.QueryEscape(name))
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(bodyBytes) == `{"Error":"No results found"}` ||
		strings.HasPrefix(string(bodyBytes), `"result not found.`) {
		return nil, errors.New("no results found")
	}

	var m movie
	if err := json.Unmarshal(bodyBytes, &m); err != nil {
		return nil, err
	}

	return &m, nil
}

func getPoster(m *movie) error {
	resp, err := http.Get(m.PosterURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.ContentLength == 0 {
		return errors.New("Could not fetch poster")
	}

	file, err := os.Create(m.Name + ".jpg")
	if err != nil {
		return err
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	file.Close()
	return nil
}
