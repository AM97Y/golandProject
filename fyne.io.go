package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)


const (
	dumpRaw = false
	zip     = "150000,ru"
	api     =  "99a2c3bc9d5f4bfde5eb78c75845393d"
)

var (
	weatherKeys = map[string]bool{"main": false, "wind": false, "coord": false, "weather": true, "sys": false, "clouds": false}
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
	// Создаем графическую часть
    a := app.New()
    win := a.NewWindow("Hello World")
    win.SetContent(widget.NewVBox(
        widget.NewLabel("Hello World!"),
        widget.NewButton("Quit", func() {
            a.Quit()
        }),
    ))
    //win.ShowAndRun()
    // Читаем данные.
	urlString := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?zip=%s&APPID=%s", zip, api)
	u, err := url.Parse(urlString)
	res, err := http.Get(u.String())
	if err != nil {
		log.Fatal(err)
	}

	jsonBlob, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	
	var data map[string]interface{}

	if dumpRaw {
		fmt.Printf("blob = %s\n\n", jsonBlob)
	}

	err = json.Unmarshal(jsonBlob, &data)
	if err != nil {
		fmt.Println("error:", err)
	}

	if dumpRaw {
		fmt.Printf("%+v", data)
	}


	var vBox = widget.NewVBox()

	for k, v := range data {
		val, isArray := isKey(k)
		if val {
			fmt.Printf("%s:\n", k)
			if isArray {
				for i := 0; i < len(v.([]interface{})); i++ {
					nmap := v.([]interface{})[i].(map[string]interface{})
					for kk, vv := range nmap {
						fmt.Printf("\tthe %s is %v\n", kk, vv)
					}
				}
			} else {
				nmap := v.(map[string]interface{})
				for kk, vv := range nmap {
					if isTempVal(kk) {
						farenTemp := faren(vv.(float64))

						vBox.Append(widget.NewVBox(
							widget.NewLabel(fmt.Sprintf("\tthe %s is %f\n", kk, farenTemp))))

						fmt.Printf("\tthe %s is %f\n", kk, farenTemp)
					} else if isSunVal(kk) {
						sunTime := time.Unix(int64(vv.(float64)), 0)
						vBox.Append(widget.NewVBox(
							widget.NewLabel(fmt.Sprintf("\tthe %s at %v\n", kk, sunTime))))
						fmt.Printf("\tthe %s at %v\n", kk, sunTime)
					} else {
						vBox.Append(widget.NewVBox(
							widget.NewLabel(fmt.Sprintf("\tthe %s is %v\n", kk, vv))))
						fmt.Printf("\tthe %s is %v\n", kk, vv)
					}
				}
			}
		} else {
		}
	}
	win.SetContent(vBox)
	win.ShowAndRun()}

func isKey(k string) (ok bool, isArray bool) {
	isArray, ok = weatherKeys[k]
	return ok, isArray
}
// Надо С
func faren(kelvin float64) float64 {
	return 9.0/5.0*(kelvin-273.0) + 32.0
}

func isTempVal(k string) bool {
	return strings.Contains(k, "temp")
}

func isSunVal(k string) bool {
	return strings.Contains(k, "sun")
}
