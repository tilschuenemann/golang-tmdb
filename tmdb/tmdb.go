package tmdb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Genre struct {
	Name string `json:"name"`
	Id   int    `json:"id"`
}

type GenreCollection struct {
	Genre []Genre `json:"genres"`
}

type Collection struct {
	Id           int    `json:"id,,omitempty"`
	Name         string `json:"name,,omitempty"`
	PosterPath   string `json:"poster_path,,omitempty"`
	BackdropPath string `json:"backdrop_path,,omitempty"`
}

type ProductionCompany struct {
	Id            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	Iso_3166_1 string `json:"iso_3166_1"`
	Name       string `json:"name"`
}

type SpokenLanguage struct {
	EnglishName string `json:"englishName"`
	Iso_639_1   string `json:"iso_639_1"`
	En          string `json:"en"`
}

type MovieDetail struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection Collection          `json:"belongs_to_collection,omitempty"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	Id                  int                 `json:"id"`
	ImdbId              string              `json:"imdb_id"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float64             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountry   []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         json.Number         `json:"vote_average"`
	Votecount           int                 `json:"vote_count"`
}

type ResultPage struct {
	Page         int           `json:"page"`
	Results      []MovieDetail `json:"results"`
	TotalPages   int           `json:"total_pages"`
	TotalResults int           `json:"total_results"`
}

func getAccessToken() string {
	token, exists := os.LookupEnv("TMDB_ACCESS_TOKEN")
	if !exists {
		log.Fatal("TMDB_ACCESS_TOKEN not set!")
	}
	return token

}

func GetGenres() GenreCollection {
	url := "https://api.themoviedb.org/3/genre/movie/list?language=en"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", getAccessToken()))
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var results GenreCollection
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	return results
}

func GetMovieDetail(tmdbId int, resultChannel chan MovieDetail) {
	url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d", tmdbId)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", getAccessToken()))
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var results MovieDetail
	err = json.Unmarshal(body, &results)
	if err != nil {
		log.Fatal(err)
	}

	resultChannel <- results

}

func SearchMovie(query string, year string, resultChannel chan int) {
	url := "https://api.themoviedb.org/3/search/movie"

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", getAccessToken()))
	req.Header.Add("accept", "application/json")

	q := req.URL.Query()
	q.Add("query", query)
	q.Add("include_adult", fmt.Sprintf("%t", true))
	q.Add("language", "en-US")
	q.Add("year", year)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var results ResultPage
	err = json.Unmarshal(body, &results)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	// TODO sort by popularity first
	resultChannel <- results.Results[0].Id
}
