package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	//"net/url"
)

type RecommendationResponse struct {
    Recommendations []string `json:"recommendations"`
}


func recommendByCosineHandler(w http.ResponseWriter, req *http.Request) {
    movieTitle := req.URL.Query().Get("movie_title")
    numRecommendations := req.URL.Query().Get("num_recommendations")

    recommendations := recommend_by_cosine(movie_title, num_recommendations)

    response := RecommendationResponse{recommendations}
    jsonResponse, err := json.Marshal(response)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jsonResponse)
}

func main() {
	http.HandleFunc("/recommend_by_cosine", recommendByCosineHandler)
    fmt.Println("Starting Server on Port 8080...")
    http.ListenAndServe(":8080", nil)
}