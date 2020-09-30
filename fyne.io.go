package main

import "net/http"
import (
        "fyne.io/fyne/app"
    	"fyne.io/fyne/widget"
    	"fmt"
    	"log"
    	"encoding/json"

        // Shortening the import reference name seems to make it a bit easier
        owm "github.com/briandowns/openweathermap"
)

type WeatherInfo struct {
  List []WeatherListItem `json:list`
}



type WeatherListItem struct {
  Dt      int           `json:dt`
  Main    WeatherMain   `json:main`
  Weather []WeatherType `json:weather`
}

type WeatherMain struct {
  Temp      float32 `json:temp`
  FeelsLike float32 `json:feels_like`
  Humidity  int     `json:humidity`
}

type WeatherType struct {
  Icon string `json:icon`
}
// "99a2c3bc9d5f4bfde5eb78c75845393d"
func main() {
    a := app.New()
    w, err := owm.NewCurrent("F", "ru", "99a2c3bc9d5f4bfde5eb78c75845393d") // fahrenheit (imperial) with Russian output
    if err != nil {
        log.Fatalln(err)
    }

    w.CurrentByName("Yaroslavl")
    fmt.Println(w)
    win := a.NewWindow("Hello World")
    win.SetContent(widget.NewVBox(
        widget.NewLabel("Hello World!"),
        widget.NewButton("Quit", func() {
            a.Quit()
        }),
    ))
    win.ShowAndRun()
}

func getWeatherForecast(result interface{}) error {
	var url = fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=Voronezh&cnt=4&units=metric&appid=%s") // инициализируйте со своим ключом
	response, err := http.Get(url)
	if err != nil {
	fmt.Print(err)
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(result)
}
