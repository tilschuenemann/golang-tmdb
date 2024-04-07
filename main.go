package main

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"strings"

	"tilschuenenmann.com/golang-tmdb/tmdb"
)

type SearchTuple struct {
	Query string
	Year  string
}

func main() {
	movie_dir_path := "your_path_here"
	matching_regex := `^(\d{4})\s(.*)$`
	searchTuples := GetMoviesFromDir(movie_dir_path, matching_regex)

	genreCollection := tmdb.GetGenres()
	WriteGenres(genreCollection)

	tmddbIds := GetTmdbIds(searchTuples)
	movieDetails := GetMovieDetails(tmddbIds)
	WriteMovieDetails(movieDetails)
}

func GetTmdbIds(searchTuples []SearchTuple) []int {
	resultChannel := make(chan int)
	for _, st := range searchTuples {
		go tmdb.SearchMovie(st.Query, st.Year, resultChannel)
	}

	movieCount := len(searchTuples)
	results := make([]int, movieCount)

	for i := range movieCount {
		results[i] = <-resultChannel
	}
	return results
}

func GetMovieDetails(tmdbIds []int) []tmdb.MovieDetail {
	resultChannel := make(chan tmdb.MovieDetail)
	for _, t := range tmdbIds {
		go tmdb.GetMovieDetail(t, resultChannel)
	}

	movieCount := len(tmdbIds)
	results := make([]tmdb.MovieDetail, movieCount)

	for i := range movieCount {
		results[i] = <-resultChannel
	}
	return results
}

func WriteMovieDetails(moviedetails []tmdb.MovieDetail) {
	file, err := os.Create("moviedetails.jsonl")
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}
	defer file.Close()
	for _, m := range moviedetails {
		line, err := json.Marshal(m)
		if err != nil {
			log.Printf("Error encoding JSON: %s", err)
			return
		}
		line = append(line, '\n')

		if _, err := file.Write(line); err != nil {
			log.Printf("Error writing to file: %s", err)
			return
		}
	}

}

func WriteGenres(genreCollection tmdb.GenreCollection) {
	file, err := os.Create("genres.jsonl")
	if err != nil {
		log.Fatalf("Error creating file: %s", err)
	}
	defer file.Close()

	line, err := json.Marshal(genreCollection)
	if err != nil {
		log.Printf("Error encoding JSON: %s", err)
		return
	}
	line = append(line, '\n')

	if _, err := file.Write(line); err != nil {
		log.Printf("Error writing to file: %s", err)
		return
	}

}

func GetMoviesFromDir(movie_dir_path string, matching_regex string) []SearchTuple {
	files, err := os.ReadDir(movie_dir_path)
	if err != nil {
		log.Fatal(err)
	}

	var matches []SearchTuple

	pat, _ := regexp.Compile(matching_regex)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		fname := file.Name()

		matched := pat.MatchString(fname)
		if !matched {
			continue
		}

		parts := strings.SplitN(fname, " ", 2)

		matches = append(matches, SearchTuple{
			Query: parts[1],
			Year:  parts[0],
		})

	}

	return matches
}
