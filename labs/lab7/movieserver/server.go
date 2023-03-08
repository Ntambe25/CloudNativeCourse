// Package main implements a server for movieinfo service.
package main

import (
	"context"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/NaseerLodge/CloudNativeCourse/labs/lab7/movieapi"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

// server is used to implement movieapi.MovieInfoServer
type server struct {
	movieapi.UnimplementedMovieInfoServer
}

// Map representing a database
// Key is type string
// Value is a slice of type string
var moviedb = map[string][]string{"Pulp fiction": []string{"1994", "Quentin Tarantino", "John Travolta,Samuel Jackson,Uma Thurman,Bruce Willis"}}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) GetMovieInfo(ctx context.Context, in *movieapi.MovieRequest) (*movieapi.MovieReply, error) {

	//Get title of the movie and print it to the server
	title := in.GetTitle()
	log.Printf("Received: %v", title)
	reply := &movieapi.MovieReply{}

	//This if-else block checks whether the movie title is present in the moviedb map.
	//If it is not found, reply=empty.
	//If it is found, the movie's details are extracted from the database and filled into the movieapi.MovieReply object.
	if val, ok := moviedb[title]; !ok { // Title not present in database
		return reply, nil
	} else {

		//Converts year from string to integer
		if year, err := strconv.Atoi(val[0]); err != nil {
			//Error check
			reply.Year = -1
		} else {
			//If no error, then save year as int
			reply.Year = int32(year)
		}
		reply.Director = val[1]

		//Helps split the cast wherever there is a ,
		//cast... is used to append each of the element present in the
		//cast one by one
		cast := strings.Split(val[2], ",")
		reply.Cast = append(reply.Cast, cast...)

	}

	return reply, nil

}

// GetMovieInfo implements movieapi.MovieInfoServer
func (s *server) SetMovieInfo(ctx context.Context, in *movieapi.MovieData) (*movieapi.Status, error) {

	// Get movie data
	title := in.GetTitle()
	year := in.GetYear()
	director := in.GetDirector()
	cast := in.GetCast()

	reply := &movieapi.Status{}
	reply.Code = "fail"

	//Check the moviedb map to see if the movie is already set
	//If the title is not present enter it into the map
	if _, ok := moviedb[title]; !ok {

		//Create a new slice
		moviedata := make([]string, 0)

		//Append year, director and cast to the slice
		moviedata = append(moviedata, strconv.Itoa(int(year)))
		moviedata = append(moviedata, director)
		moviedata = append(moviedata, cast...)

		//The key will be the title
		//and the value will be the slice moviedata
		moviedb[title] = moviedata //add movie data into database
		reply.Code = "success"
	} else {
		return reply, errors.New("movie already exist in database")
	}
	return reply, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	movieapi.RegisterMovieInfoServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
