package weather

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Tests response
func TestParseResponse(t *testing.T) {
	t.Parallel()
	// Reads from the JSON file
	data, err := ioutil.ReadFile("testdata/weather_data.json")
	if err != nil {
		t.Fatal(err)
	}
	want := Conditions{
		//Conditions from weather_data
		Summary:     "Clouds",
		Temperature: 281.33,
		Pressure:    1000,
		Humidity:    90,
		Speed:       3.09,
	}
	got, err := ParseResponse(data)
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestParseResponseEmpty(t *testing.T) {
	t.Parallel()
	_, err := ParseResponse([]byte{})
	if err == nil {
		t.Fatal("want error parsing empty response, got nil")
	}
}

func TestParseResponseInvalid(t *testing.T) {
	t.Parallel()
	data, err := ioutil.ReadFile("testdata/weather_invalid_data.json")
	if err != nil {
		t.Fatal(err)
	}
	_, err = ParseResponse(data)
	if err == nil {
		t.Fatal("want error parsing invalid response, got nil")
	}
}

func TestFormatURL(t *testing.T) {
	t.Parallel()
	c := NewClient("dummyAPIKey")
	location := "Paris,FR"
	want := "https://api.openweathermap.org/data/2.5/weather?q=Paris%2CFR&appid=dummyAPIKey"
	got := c.FormatURL(location)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFormatURLSpaces(t *testing.T) {
	t.Parallel()
	c := NewClient("dummyAPIKey")
	location := "Wagga Wagga,AU"
	want := "https://api.openweathermap.org/data/2.5/weather?q=Wagga+Wagga%2CAU&appid=dummyAPIKey"
	got := c.FormatURL(location)
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestSimpleHTTP(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	client := ts.Client()
	resp, err := client.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	want := http.StatusOK
	got := resp.StatusCode
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestGetWeather(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/weather_data.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	defer ts.Close()
	c := NewClient("dummyAPIkey")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	want := Conditions{
		Summary:     "Clouds",
		Temperature: 281.33,
		Pressure:    1000,
		Humidity:    90,
		Speed:       3.09,
	}
	got, err := c.GetWeather("Paris,FR")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestFahrenheit(t *testing.T) {
	t.Parallel()
	input := Temperature(274.15)
	want := 33.8
	got := input.Fahrenheit()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestPressure(t *testing.T) {
	t.Parallel()
	// Tests added functions, that is to say tests if the different parameters meant to be returned are correct
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/weather_data.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	defer ts.Close()

	c := NewClient("dummyAPIkey")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	// If Pressure is not equal to 1000, it throws an error
	input, _ := c.GetWeather("Paris,FR")
	want := Pressure(1000)
	got := input.Pressure
	if !cmp.Equal(want, got) {
		fmt.Println("Pressure not working, is equal to", want)
		t.Error(cmp.Diff(want, got))
	}
}

func TestHumidity(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/weather_data.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	// Tests if weather.go successfully sends information regarding Humidity
	// Otherwise throws error
	defer ts.Close()
	c := NewClient("dummyAPIkey")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	input, _ := c.GetWeather("Paris,FR")
	want := Humidity(90)
	got := input.Humidity
	if !cmp.Equal(want, got) {
		fmt.Println("Humidity not working, is equal to", want)
		t.Error(cmp.Diff(want, got))
	}
}

func TestWindspeed(t *testing.T) {
	t.Parallel()
	ts := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/weather_data.json")
		if err != nil {
			t.Fatal(err)
		}
		defer f.Close()
		io.Copy(w, f)
	}))
	// Tests if returns correct information about windspeed.
	defer ts.Close()
	c := NewClient("dummyAPIkey")
	c.BaseURL = ts.URL
	c.HTTPClient = ts.Client()
	input, _ := c.GetWeather("Paris,FR")
	want := Speed(3.09)
	got := input.Speed
	if !cmp.Equal(want, got) {
		fmt.Println("Windspeed not working, is equal to", want)
		t.Error(cmp.Diff(want, got))
	}
}
