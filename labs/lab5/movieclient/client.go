// Package main imlements a client for movieinfo service
package main

import (
	"context"
	"log"
	"os"
	"time"

	//"gitlab.com/arunravindran/cloudnativecourse/lab5-grpc/movieapi"
	"github.com/Ntambe25/CloudNativeCourse/labs/lab5/movieapi"
	"google.golang.org/grpc"
)

const (
	address      = "localhost:50051"
	defaultTitle = "Pulp fiction"
)

func main() {
	// Setting up a connection to the gRPCserver.

	// grpc.Dial function used to create a new connection
	// grpc.WithInsecure specifies connection should not use TLS encryption
	// grpc.WithBlock specifies function should block until connection is established
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	// If connection is not established, ERROR is printed
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	// defer ensures that the connection is closed before Main func exits
	defer conn.Close()

	// Creates new client 
	//movieapi.NewMovieInfoClient takes connection as an input and
	// returns a client that can be used to make requests to the server
	c := movieapi.NewMovieInfoClient(conn)

	// Contact the server and print out its response.
	title := defaultTitle

	// If Cmd Args. > 1, 'title' is set to value of 1st Arg.
	if len(os.Args) > 1 {
		title = os.Args[1]
	}

	// Timeout if server doesn't respond
	// Creates a context to manage lifecycle of the request
	// If server does not respond within timeout period, request is cancelled
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Call the GetMovieInfo on gRPC server's MovieInfo service
	// Passes 'ctx' 'MovieRequest' msg. with title 
	r, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: title})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	// IF call is successful, server's response is returned in 'r'
	log.Printf("Movie Info for %s %d %s %v", title, r.GetYear(), r.GetDirector(), r.GetCast())

	//Giving the input MoviData to SetMovie Info with the title, year, director and Cast
	status, err := c.SetMovieInfo(ctx, &movieapi.MovieData{Title: "Top Gun Maverick", Year: 2022, Director: "Joseph Kosinski", Cast: []string{"Tom Cruise, Miles Teller, Jennifer Connelly"}})

	//Error Check
	if err != nil {
		log.Fatalf("could not set movie info: %v", err)
		log.Fatalf("SetMovieInfo Status: %v", status)
	}

	//Aftter setting the MovieData, now using GetMoveInfo to output the details
	r1, err := c.GetMovieInfo(ctx, &movieapi.MovieRequest{Title: "Top Gun Maverick"})
	if err != nil {
		log.Fatalf("could not get movie info: %v", err)
	}
	log.Printf("Movie Info for Top Gun Maverick: %d %s %v", r1.GetYear(), r1.GetDirector(), r1.GetCast())
}
