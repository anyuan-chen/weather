package main

import (
	"fmt"
	"net/http"
)

type weatherData struct {
	Name string `json:"name"`
	Main struct {
		Kelvin float64 `json:"temp"`
	} `json:"main"`
}

func queryOpenWeather(city string)(weatherData, error){
	apiCall := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid={API key}", city, )

	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?q={city name}&appid={API key}")
}

func getWeather(w http.ResponseWriter, r *http.Request){
	w.Write([]byte("hi"))
}

func main(){
	fmt.Println("hi")
	http.HandleFunc("/weather", getWeather)
	http.ListenAndServe(":8080", nil)
}