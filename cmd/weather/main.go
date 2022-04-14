package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)


type weatherProvider interface {
	temperature(city string) (float64, error)
}

type openWeatherMap struct{}
func (w openWeatherMap) temperature(city string)(float64, error) {
	apiKey := os.Getenv("OPENWEATHERMAP_API_KEY") 
	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)
	response, err := http.Get(url)
	if err != nil {
		return -1.0, err
	}
	defer response.Body.Close()
	var recievedData struct {
		Name string `json:"name`
		Main struct {
			Kelvin float64 `json:"temp"`
		} `json:"main"`
	}
	if err := json.NewDecoder(response.Body).Decode(&recievedData); err != nil {
		return -1.0, err
	}
	return recievedData.Main.Kelvin, nil
}

type visualCrossing struct{}
func (vc visualCrossing) temperature (city string) (float64, error) {
	apiKey := os.Getenv("VISUALCROSSING_API_KEY")
	url := fmt.Sprintf("https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/%s?unitGroup=metric&key=%s&contentType=json", city, apiKey)
	response, err := http.Get(url)
	if err != nil {
		return -1.0, err
	}
	defer response.Body.Close()
	var recievedData struct {
		CurrentConditions struct {
			Temp float64 `json:"temp"`
		} `json:"currentConditions"`
	}
	if err := json.NewDecoder(response.Body).Decode(&recievedData); err != nil {
		return -1.0, err
	}
	return recievedData.CurrentConditions.Temp + 273.0, nil
}

type weather struct {
	Temperature float64 `json:"temperature"`
}
func getWeather(w http.ResponseWriter, r *http.Request){
	params := r.URL.Query()
	if len(params["city"]) == 0 {
		http.Error(w, "invalid request", http.StatusInternalServerError)
		return
	}

	weatherProviders := []weatherProvider{visualCrossing{}, openWeatherMap{}}
	var totalTemp float64
	for _ , provider := range weatherProviders {
		temp, error := provider.temperature(params["city"][0])
		if error != nil {
			http.Error(w, error.Error(), http.StatusInternalServerError)
			return 
		}
		totalTemp += temp
	}
	totalTemp = totalTemp / float64(len(weatherProviders))
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(w).Encode(weather{totalTemp})
}



func main(){
	enverr := godotenv.Load()
	if (enverr != nil){
		log.Fatal("Error loading the .env file")
	}
	http.HandleFunc("/weather", getWeather)
	http.ListenAndServe(":8080", nil)
}