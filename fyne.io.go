package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
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
func getWeatherForecast(result interface{}) error {
	var url = fmt.Sprintf("https://api.openweathermap.org/data/2.5/forecast?q=Voronezh&cnt=4&units=metric&appid=%s", api) // инициализируйте со своим ключом
	response, err := http.Get(url)
	if err != nil {
		fmt.Print(err)
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(result)
}
func main() {
	// Создаем графическую часть
    a := app.New()
    win := a.NewWindow("Программа для просмотра погоды")
    win.SetContent(widget.NewVBox(
        widget.NewButton("Quit", func() {
            a.Quit()
        }),
    ))

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
			vBox.Append(widget.NewVBox(
				widget.NewLabel(fmt.Sprintf("\t%s\n", k))))

			if isArray {
				for i := 0; i < len(v.([]interface{})); i++ {
					nmap := v.([]interface{})[i].(map[string]interface{})
					for kk, vv := range nmap {
						if strings.Contains(kk, "icon") {
							var resource, _ = fyne.LoadResourceFromURLString(fmt.Sprintf("http://openweathermap.org/img/wn/%s.png", vv)) // создаем статический ресурс, содержащий иконку погоды
							var icon = widget.NewIcon(resource)
							vBox.Append(icon)
						}
						vBox.Append(widget.NewVBox(
							widget.NewLabel(fmt.Sprintf("\tthe %s is %v\n", kk, vv))))

						fmt.Printf("\tthe %s is %v\n", kk, vv)
					}
				}
			} else {
				nmap := v.(map[string]interface{})
				for kk, vv := range nmap {
					if isTempVal(kk) {
						farenTemp := faren(vv.(float64))

						//vBox.Append(widget.NewVBox(
						//	widget.NewLabel(fmt.Sprintf("\tthe %s is %f\n", kk, farenTemp))))

						fmt.Printf("\tthe %s is %f\n", kk, farenTemp)
					} else if isSunVal(kk) {
						sunTime := time.Unix(int64(vv.(float64)), 0)
						//vBox.Append(widget.NewVBox(
							//widget.NewLabel(fmt.Sprintf("\tthe %s at %v\n", kk, sunTime))))
						fmt.Printf("\tthe %s at %v\n", kk, sunTime)
					} else {
						//vBox.Append(widget.NewVBox(
							//widget.NewLabel(fmt.Sprintf("\tthe %s is %v\n", kk, vv))))
						fmt.Printf("\tthe %s is %v\n", kk, vv)
					}
				}
			}
		} else {
		}
	}
	vBox.Append(widget.NewButton("Закрыть", func() {
		a.Quit()
	}))
	win.SetContent(vBox)
	win.ShowAndRun()}

func isKey(k string) (ok bool, isArray bool) {
	isArray, ok = weatherKeys[k]
	return ok, isArray
}
// Надо С
func faren(kelvin float64) float64 {
	fmt.Printf("\tthe____________ %s is %v\n", kelvin-273.0, kelvin)
	return kelvin-273.0
}

func isTempVal(k string) bool {
	return strings.Contains(k, "temp")
}

func isSunVal(k string) bool {
	return strings.Contains(k, "sun")
}
