package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	mongodbEndpoint = "mongodb://172.17.0.2" // Find this from the Mongo container
)

var RWLock sync.RWMutex
var col *mongo.Collection
var ctx context.Context

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

type Post struct {
	ID        primitive.ObjectID `bson:"_id"`
	Clothing  string             `bson:"clothing"`
	Price     dollars            `bson:"price"`
	Tags      string             `bson:"tags"`
	CreatedAt time.Time          `bson:"created_at"`
}

func main() {
	// create a mongo client
	client, err := mongo.NewClient(
		options.Client().ApplyURI(mongodbEndpoint),
	)
	checkError(err)

	// Connect to mongo
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)

	// Disconnect
	defer client.Disconnect(ctx)

	col = client.Database("blog").Collection("posts")

	mux := http.NewServeMux()
	mux.HandleFunc("/list", list)
	mux.HandleFunc("/create", create)
	mux.HandleFunc("/delete", delete)
	mux.HandleFunc("/price", price)
	mux.HandleFunc("/update", update)
	log.Fatal(http.ListenAndServe(":8000", mux))

}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func list(w http.ResponseWriter, req *http.Request) {

	cursor, err := col.Find(ctx, bson.M{})
	if err != nil {
		log.Fatal(err)
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var episode bson.M
		if err = cursor.Decode(&episode); err != nil {
			log.Fatal(err)
		}
		fmt.Println(episode)
	}

}

func create(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")
	priceFloat, _ := strconv.ParseFloat(price, 32)

	fmt.Println("ITEM : ", item)
	fmt.Println("PRICE: ", priceFloat)

	res, err := col.InsertOne(ctx, &Post{
		ID:        primitive.NewObjectID(),
		Clothing:  item,
		Price:     dollars(priceFloat),
		Tags:      "clothing",
		CreatedAt: time.Now(),
	})

	if err == nil {
		fmt.Printf("inserted id: %s\n", res.InsertedID.(primitive.ObjectID).Hex())
	}

}

// Delete Function (In this case all items using DeleteMany)
func delete(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")
	result, err := col.DeleteMany(ctx, bson.M{"clothing": item})
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("DeleteMany removed %v document(s)\n", result.DeletedCount)
	}
	
}

// Price Function (Returns the price of an item)
func price(w http.ResponseWriter, req *http.Request) {

	item := req.URL.Query().Get("item")

	var post Post
	err := col.FindOne(ctx, bson.M{"clothing": item}).Decode(&post);
	
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
    	return
    }
	fmt.Fprintf(w, "Price of %s:%s\n", item, post.Price)
}


// Update Function 
// takes w as type "http.ResponseWriter" used to write response back to the client
// and req is request received by the server
func update(w http.ResponseWriter, req *http.Request) {
	
	// gets "item" and "price" from the HTTP request
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	// converts the price(string) to floating point number
	priceFloat, err := strconv.ParseFloat(price, 32)


	if err != nil {
		http.Error(w, "Invalid Price Value.", http.StatusNotFound)
		return
	}

	// Creates the new "bson.M" object with key-value pair
	filter := bson.M{"clothing": item}
	// Defining the update operation
	update := bson.M{"$set": bson.M{"price": priceFloat}}

	// Opdate operation is being performed using filter, update, and ctx
	result, err := col.UpdateOne(ctx, filter, update)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// If reult of update operation failed (i.e Modified Documents == 0), Erro is Printed
	if result.ModifiedCount == 0 {
		http.Error(w, fmt.Sprintf("No Such Item. %q", item), http.StatusNotFound)
		return
	}

	fmt.Println("Updated %d Document(s). ", result.ModifiedCount)
}