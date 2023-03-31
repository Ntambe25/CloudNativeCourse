package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Recommendation struct {
	Title       string
	Description string
}

func main() {
	http.HandleFunc("/recommend", handleRecommendation)
    fmt.Println("Listening on port 8080...")
    http.ListenAndServe(":8080", nil)
}

func handleRecommendation(w http.ResponseWriter, r *http.Request) {
    var genre string
    if r.Method == "POST" {
        genre = r.FormValue("genre")
    } else {
        genre = r.URL.Query().Get("genre")
    }

    // Make API call to machine learning model container
    mlURL := "http://ml-container:5000/recommend?genre=" + url.QueryEscape(genre)
    response, err := http.Get(mlURL)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer response.Body.Close()

	var recommendations []Recommendation
    err = json.NewDecoder(response.Body).Decode(&recommendations)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    responseJSON, err := json.Marshal(recommendations)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(responseJSON)
}