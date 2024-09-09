package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
)

const (
	ColorReset  = "\033[0m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
)

var city string
var lang string
var apiKey string = "6531c553d30d4297ad9155019240809"

type WeatherResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Region    string  `json:"region"`
		Country   string  `json:"country"`
		Lat       float64 `json:"lat"`
		Lon       float64 `json:"lon"`
		TzID      string  `json:"tz_id"`
		Localtime string  `json:"localtime"`
	} `json:"location"`
	Current struct {
		LastUpdated string  `json:"last_updated"`
		TempC       float64 `json:"temp_c"`
		TempF       float64 `json:"temp_f"`
		IsDay       int     `json:"is_day"`
		Condition   struct {
			Text string `json:"text"`
			Icon string `json:"icon"`
			Code int    `json:"code"`
		} `json:"condition"`
		WindMPH    float64 `json:"wind_mph"`
		WindKPH    float64 `json:"wind_kph"`
		WindDegree int     `json:"wind_degree"`
		WindDir    string  `json:"wind_dir"`
		PressureMB float64 `json:"pressure_mb"`
		PressureIN float64 `json:"pressure_in"`
		PrecipMM   float64 `json:"precip_mm"`
		PrecipIN   float64 `json:"precip_in"`
		Humidity   int     `json:"humidity"`
		Cloud      int     `json:"cloud"`
		FeelsLikeC float64 `json:"feelslike_c"`
		FeelsLikeF float64 `json:"feelslike_f"`
		WindChillC float64 `json:"windchill_c"`
		WindChillF float64 `json:"windchill_f"`
		HeatIndexC float64 `json:"heatindex_c"`
		HeatIndexF float64 `json:"heatindex_f"`
		DewPointC  float64 `json:"dewpoint_c"`
		DewPointF  float64 `json:"dewpoint_f"`
		VisKM      float64 `json:"vis_km"`
		VisMiles   float64 `json:"vis_miles"`
		UV         float64 `json:"uv"`
		GustMPH    float64 `json:"gust_mph"`
		GustKPH    float64 `json:"gust_kph"`
	} `json:"current"`
}

func getWeather(city, lang string) (WeatherResponse, error) {
	encodedCity := url.QueryEscape(city)
	apiUrl := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&lang=%s&aqi=no", apiKey, encodedCity, lang)

	fmt.Println("Request URL:", apiUrl)

	resp, err := http.Get(apiUrl)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("could not fetch weather data: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Status Code:", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Response Body:", string(body))
		return WeatherResponse{}, fmt.Errorf("city not found or API error, status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("could not read response body: %v", err)
	}

	var result WeatherResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return WeatherResponse{}, fmt.Errorf("could not unmarshal JSON: %v", err)
	}

	if result.Location.Name == "" {
		return WeatherResponse{}, fmt.Errorf("city not found")
	}

	return result, nil
}

func showWeather(response WeatherResponse, lang string) {
	switch lang {
	case "uz":
		fmt.Println("---------------------------------------------------------")
		fmt.Printf(ColorGreen+"Shahar: "+ColorReset+"%s\n", response.Location.Name)
		fmt.Printf(ColorGreen+"Hudud: "+ColorReset+"%s\n", response.Location.Region)
		fmt.Printf(ColorGreen+"Davlat: "+ColorReset+"%s\n", response.Location.Country)
		fmt.Printf(ColorGreen+"Harorat: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.TempC, response.Current.TempF)
		fmt.Printf(ColorGreen+"Ob-havo: "+ColorReset+"%s\n", response.Current.Condition.Text)
		fmt.Printf(ColorGreen+"Shamol: "+ColorReset+"%.1f MPH (%.1f KPH), %s\n", response.Current.WindMPH, response.Current.WindKPH, response.Current.WindDir)
		fmt.Printf(ColorGreen+"Bosim: "+ColorReset+"%.1f hPa\n", response.Current.PressureMB)
		fmt.Printf(ColorGreen+"Namlik: "+ColorReset+"%d%%\n", response.Current.Humidity)
		fmt.Printf(ColorGreen+"His qilinadigan harorat: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.FeelsLikeC, response.Current.FeelsLikeF)
		fmt.Printf(ColorGreen+"Ko'rinish: "+ColorReset+"%.1f km\n", response.Current.VisKM)
		fmt.Println("---------------------------------------------------------")

	case "ru":
		fmt.Println("---------------------------------------------------------")
		fmt.Printf(ColorGreen+"Город: "+ColorReset+"%s\n", response.Location.Name)
		fmt.Printf(ColorGreen+"Регион: "+ColorReset+"%s\n", response.Location.Region)
		fmt.Printf(ColorGreen+"Страна: "+ColorReset+"%s\n", response.Location.Country)
		fmt.Printf(ColorGreen+"Температура: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.TempC, response.Current.TempF)
		fmt.Printf(ColorGreen+"Погода: "+ColorReset+"%s\n", response.Current.Condition.Text)
		fmt.Printf(ColorGreen+"Ветер: "+ColorReset+"%.1f MPH (%.1f KPH), %s\n", response.Current.WindMPH, response.Current.WindKPH, response.Current.WindDir)
		fmt.Printf(ColorGreen+"Давление: "+ColorReset+"%.1f hPa\n", response.Current.PressureMB)
		fmt.Printf(ColorGreen+"Влажность: "+ColorReset+"%d%%\n", response.Current.Humidity)
		fmt.Printf(ColorGreen+"Ощущается как: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.FeelsLikeC, response.Current.FeelsLikeF)
		fmt.Printf(ColorGreen+"Видимость: "+ColorReset+"%.1f км\n", response.Current.VisKM)
		fmt.Println("---------------------------------------------------------")

	default:
		fmt.Println("---------------------------------------------------------")
		fmt.Printf(ColorGreen+"City: "+ColorReset+"%s\n", response.Location.Name)
		fmt.Printf(ColorGreen+"Region: "+ColorReset+"%s\n", response.Location.Region)
		fmt.Printf(ColorGreen+"Country: "+ColorReset+"%s\n", response.Location.Country)
		fmt.Printf(ColorGreen+"Temperature: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.TempC, response.Current.TempF)
		fmt.Printf(ColorGreen+"Weather: "+ColorReset+"%s\n", response.Current.Condition.Text)
		fmt.Printf(ColorGreen+"Wind: "+ColorReset+"%.1f MPH (%.1f KPH), %s\n", response.Current.WindMPH, response.Current.WindKPH, response.Current.WindDir)
		fmt.Printf(ColorGreen+"Pressure: "+ColorReset+"%.1f hPa\n", response.Current.PressureMB)
		fmt.Printf(ColorGreen+"Humidity: "+ColorReset+"%d%%\n", response.Current.Humidity)
		fmt.Printf(ColorGreen+"Feels Like: "+ColorReset+"%.2f°C (%.2f°F)\n", response.Current.FeelsLikeC, response.Current.FeelsLikeF)
		fmt.Printf(ColorGreen+"Visibility: "+ColorReset+"%.1f km\n", response.Current.VisKM)
		fmt.Println("---------------------------------------------------------")
	}
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "weather",
		Short: "Weather CLI",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Print("Enter language (en, ru, uz) [default: en]: ")
			fmt.Scanln(&lang)

			lang = strings.ToLower(lang)
			if lang != "ru" && lang != "uz" {
				lang = "en"
			}

			for {
				fmt.Print("Enter city name (or 'exit' to quit): ")
				fmt.Scanln(&city)

				if strings.ToLower(city) == "exit" {
					break
				}

				if city == "" {
					fmt.Println("City name cannot be empty.")
					continue
				}

				response, err := getWeather(city, lang)
				if err != nil {
					fmt.Println("Error:", err)
					continue
				}

				showWeather(response, lang)
			}
		},
	}

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
