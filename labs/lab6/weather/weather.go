package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Temperature Datatype
type Temperature float64

// Pressure Datatype
type Pressure float64

// Humidity Datatype
type Humidity float64

// Wind Datatype
type Speed float64

// Converts Kelvin to Fahrenheit
func (t Temperature) Fahrenheit() float64 {
	return (float64(t)-273.15)*(9.0/5.0) + 32.0
}

// Struct that stores temperature in Fahrenheit
// and summary
type Conditions struct {
	Summary     string
	Temperature Temperature
	Pressure    Pressure
	Humidity    Humidity
	Speed       Speed
}

// Create a data structure to represent these response fields.
type OWMResponse struct {
	Weather []struct {
		Main string
	}
	Main struct {
		Temp     Temperature
		Pressure Pressure
		Humidity Humidity
	}
	Wind struct {
		Speed Speed
	}
}

// Struct Client
// APIKey     -> Custom key provided by openweather to change stuff on the website
// BaserURL   -> "https://api.openweathermap.org"
// HTTPClient ->  HTTP client from the net/http package
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// Client constructor
func NewClient(key string) *Client {
	return &Client{
		APIKey:  key,
		BaseURL: "https://api.openweathermap.org",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Input -> City
// Output -> URL
func (c Client) FormatURL(location string) string {
	location = url.QueryEscape(location)
	return fmt.Sprintf("%s/data/2.5/weather?q=%s&appid=%s", c.BaseURL, location, c.APIKey)

}

// Makes GET request to the service with the URL,
// and calls a ParseReponse() helper method to parse the JSON data.
func (c *Client) GetWeather(location string) (Conditions, error) {
	URL := c.FormatURL(location)
	resp, err := c.HTTPClient.Get(URL)
	if err != nil {
		return Conditions{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return Conditions{}, fmt.Errorf("could not find location: %s ", location)
	}
	if resp.StatusCode != http.StatusOK {
		return Conditions{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Conditions{}, err
	}
	conditions, err := ParseResponse(data)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, nil
}

// Helper method to parse the JSON data.
func ParseResponse(data []byte) (Conditions, error) {
	var resp OWMResponse
	err := json.Unmarshal(data, &resp)
	if err != nil {
		return Conditions{}, fmt.Errorf("invalid API response %s: %w", data, err)
	}
	if len(resp.Weather) < 1 {
		return Conditions{}, fmt.Errorf("invalid API response %s: require at least one weather element", data)
	}
	conditions := Conditions{
		Summary:     resp.Weather[0].Main,
		Temperature: resp.Main.Temp,
		Pressure:    resp.Main.Pressure,
		Humidity:    resp.Main.Humidity,
		Speed:       resp.Wind.Speed,
	}
	return conditions, nil
}

func Get(location, key string) (Conditions, error) {
	c := NewClient(key)
	conditions, err := c.GetWeather(location)
	if err != nil {
		return Conditions{}, err
	}
	return conditions, nil
}

func RunCLI() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s LOCATION\n\nExample: %[1]s London,UK", os.Args[0])
		os.Exit(1)
	}
	location := os.Args[1]
	key := os.Getenv("OPENWEATHERMAP_API_KEY")
	if key == "" {
		fmt.Fprintln(os.Stderr, "Please set the environment variable OPENWEATHERMAP_API_KEY")
		os.Exit(1)
	}
	conditions, err := Get(location, key)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//fmt.Printf("%s %.1fº\n", conditions.Summary, conditions.Temperature.Fahrenheit())
	fmt.Println("Summary     -> ", conditions.Summary)
	fmt.Print("Temperature -> ", conditions.Temperature.Fahrenheit())
	fmt.Print("º\n")
	fmt.Println("Pressure    -> ", conditions.Pressure)
	fmt.Println("Humidity    -> ", conditions.Humidity)
	fmt.Println("Wind Speed  -> ", conditions.Speed)

}
