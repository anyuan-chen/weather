package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func queryOpenWeather(city string)(weatherData, error){
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY") 
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	response, err := http.Get(url)
	if err != nil {
		return weatherData{}, err
	}
	defer response.Body.Close()
	var recievedData weatherData
	if err := json.NewDecoder(response.Body).Decode(&recievedData); err != nil {
		return weatherData{}, err
	}
	return recievedData, nil
}

func getWeather(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	if len(params["city"]) == 0 {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}
	openWeatherData, openWeatherError := queryOpenWeather(params["city"][0])
	if openWeatherError != nil {
		http.Error(w, openWeatherError.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(openWeatherData)
}

func main(){
	enverr := godotenv.Load()
	if (enverr != nil){
		log.Fatal("Error loading the .env file")
	}
	http.HandleFunc("/weather", getWeather)
	http.ListenAndServe(":8080", nil)
}