package main

import (
	"encoding/json"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
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
	a.Settings().SetTheme(theme.DarkTheme())
    win := a.NewWindow("Погода Ярославль")
    win.SetContent(widget.NewVBox(
        widget.NewButton("Quit", func() {
            a.Quit()
        }),

    ))

	f, err := os.Create("result.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Println(l, "bytes written successfully")


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
	var groupBox = widget.NewVBox()
	var groupBox1 = widget.NewVBox()
	var groupBox2 = widget.NewVBox()
	var groupBox3 = widget.NewVBox()

	var time_2 = time.Now().String()
	vBox.Append(widget.NewGroup(time_2))

	for k, v := range data {
		val, isArray := isKey(k)
		if val {
			_, err := f.WriteString(fmt.Sprintf("%s:\n", k))
			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}
			fmt.Printf("%s:\n", k)
			//vBox.Append(widget.NewVBox(
			//	widget.NewLabel(fmt.Sprintf("%s", k))))

			if isArray {

				for i := 0; i < len(v.([]interface{})); i++ {
					nmap := v.([]interface{})[i].(map[string]interface{})
					for kk, vv := range nmap {
						if strings.Contains(kk, "icon") {

							var resource, _ = fyne.LoadResourceFromURLString(fmt.Sprintf("http://openweathermap.org/img/wn/%s.png", vv)) // создаем статический ресурс, содержащий иконку погоды
							var icon = widget.NewIcon(resource)
							vBox.Append(icon)
						} else {
							groupBox.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("the %s is %v", kk, vv),
									fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Monospace: true})))
						}

						fmt.Printf("\tthe %s is %v\n", kk, vv)
						_, err := f.WriteString(fmt.Sprintf("\tthe %s is %v\n", kk, vv))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}
					}
				}

			} else {

				nmap := v.(map[string]interface{})
				for kk, vv := range nmap {
					if isTempVal(kk) {
						farenTemp := C2K(vv.(float64))

						groupBox1.Append(widget.NewVBox(
							widget.NewLabel(fmt.Sprintf("the %s is %s", kk, farenTemp))))

						fmt.Printf("\tthe %s is %s\n", kk, farenTemp)
						_, err := f.WriteString(fmt.Sprintf("\tthe %s is %s\n", kk, farenTemp))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}
					} else if isSunVal(kk) {
						sunTime := time.Unix(int64(vv.(float64)), 0)
						groupBox2.Append(widget.NewVBox(
							widget.NewLabelWithStyle(fmt.Sprintf("the %s at %v", kk, sunTime),
								fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})))
						fmt.Printf("\tthe %s at %v\n", kk, sunTime)
						_, err := f.WriteString(fmt.Sprintf("\tthe %s at %v\n", kk, sunTime))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}

					} else {
						if isSpeed(kk) {
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Скрость:  %v", vv), fyne.TextAlignCenter, fyne.TextStyle{ Monospace: true})))
						}
						if isClouds(kk) {
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Небо:  %v", vv), fyne.TextAlignLeading, fyne.TextStyle{ Monospace: true})))
						}
						if isFeels_like(kk) {
							var temp = C2K(vv.(float64))
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Температура чуствуется как:  %v", temp), fyne.TextAlignLeading, fyne.TextStyle{ Monospace: true})))
						}
						fmt.Printf("\tthe %s is %v\n", kk, vv)
						_, err := f.WriteString(fmt.Sprintf("\tthe %s is %v\n", kk, vv))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}
					}
				}
			}
		} else {
		}
	}
	vBox.Append(widget.NewHBox(groupBox2))
	vBox.Append(widget.NewHBox(groupBox, groupBox1))
	vBox.Append(widget.NewHBox(groupBox3))
	vBox.Append(widget.NewButton("Закрыть", func() {
		a.Quit()
	}))

	//vBox.Append(widget.NewButton("Обновить", func() {
	//	a.Quit()
	//	main()
	//}))
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	win.SetContent(vBox)
	win.ShowAndRun()}

func isKey(k string) (ok bool, isArray bool) {
	isArray, ok = weatherKeys[k]
	return ok, isArray
}

func C2K(kelvin float64) string {
	// fmt.Printf("\tthe____________ %s is %v\n", kelvin-273.0, kelvin)
	return fmt.Sprintf("%.2f °C", kelvin-273.0)
}

func isTempVal(k string) bool {
	return strings.Contains(k, "temp")
}

func isSunVal(k string) bool {
	return strings.Contains(k, "sun")
}
func isSpeed(k string) bool {
	return strings.Contains(k, "speed")
}

func isFeels_like(k string) bool {
	return strings.Contains(k, "feels_like")
}

func isClouds(k string) bool {
	return strings.Contains(k, "all")
}