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

func main() {
	// Создаем графическую часть
    a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())
    win := a.NewWindow("Погода Ярославль")
    //win.SetFullScreen(true)
    win.CenterOnScreen()
	showInfo(win)
	win.ShowAndRun()

}

func showInfo(win fyne.Window ) {

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

	var time2 = time.Now().String()
	vBox.Append(widget.NewGroup(time2))

	for k, v := range data {
		val, isArray := isKey(k)
		if val {
			fmt.Printf("%s:\n", k)
			if isArray {
				for i := 0; i < len(v.([]interface{})); i++ {
					nmap := v.([]interface{})[i].(map[string]interface{})
					for kk, vv := range nmap {
						if strings.Contains(kk, "icon") {

							var resource, _ = fyne.LoadResourceFromURLString(fmt.Sprintf("http://openweathermap.org/img/wn/%s.png", vv)) // создаем статический ресурс, содержащий иконку погоды
							var icon = widget.NewIcon(resource)
							vBox.Append(icon)
						} else {
							if strings.Contains(kk, "description"){
								groupBox.Append(widget.NewVBox(
									widget.NewLabelWithStyle(fmt.Sprintf("%v", vv),
										fyne.TextAlignTrailing, fyne.TextStyle{Bold: true, Monospace: true})))
							}
						}
						fmt.Printf("\tthe %s is %v\n", kk, vv)
					}
				}

			} else {
				nmap := v.(map[string]interface{})
				for kk, vv := range nmap {
					if isTempVal(kk) {
						farenTemp := C2K(vv.(float64))
						if strings.Contains(kk, "temp_max") {
							groupBox1.Append(widget.NewVBox(
								widget.NewLabel(fmt.Sprintf("Максимальная температура %s", farenTemp))))
						} else if strings.Contains(kk, "temp_min") {
							groupBox1.Append(widget.NewVBox(
								widget.NewLabel(fmt.Sprintf("Минимальная температура %s", farenTemp))))
						} else {
							groupBox1.Append(widget.NewVBox(
								widget.NewLabel(fmt.Sprintf("Температура %s", farenTemp))))
						}
						fmt.Printf("\tthe %s is %s\n", kk, farenTemp)
					} else if isSunVal(kk) {
						sunTime := time.Unix(int64(vv.(float64)), 0)
						if strings.Contains(kk, "sunrise") {
							groupBox2.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Рассвет в  %v", sunTime),
									fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})))
						} else {
							groupBox2.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Закат в %v", sunTime),
									fyne.TextAlignCenter, fyne.TextStyle{Bold: true, Monospace: true})))
						}
					} else {
						if isSpeed(kk) {
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Скорость ветра (м/c):  %v", vv), fyne.TextAlignLeading, fyne.TextStyle{ Monospace: true})))
						}
						if isClouds(kk) {
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Затянутость неба:  %v ", vv.(float64)), fyne.TextAlignLeading, fyne.TextStyle{ Monospace: true})))
						}
						if isFeelsLike(kk) {
							var temp = C2K(vv.(float64))
							groupBox3.Append(widget.NewVBox(
								widget.NewLabelWithStyle(fmt.Sprintf("Температура ощущается как:  %v", temp), fyne.TextAlignLeading, fyne.TextStyle{ Monospace: true})))
						}
						fmt.Printf("\tthe %s is %v\n", kk, vv)

					}
				}
			}
		} else {
		}
	}
	vBox.Append(widget.NewHBox(groupBox2))
	vBox.Append(widget.NewHBox(groupBox1, groupBox))
	vBox.Append(widget.NewHBox(groupBox3))
	vBox.Append(widget.NewButton("Обновить", func() {

		win.Content().Refresh()
		showInfo(win)
		fmt.Println("Updated")
		win.Show()
	}))
	vBox.Resize(fyne.Size{Width: 900, Height: 600})

	input := widget.NewEntry()
	input.SetText("/home/monsterpc/result.txt")
	input.SetPlaceHolder("Путь для сохранения отчета")

	content := widget.NewVBox(input, widget.NewButton("Сохранить полный отчет", func() {
		log.Println("Content was:", input.Text)
		save(input.Text)
		fmt.Println("Saved")
	}))

	vBox.Append(widget.NewVBox(
		content))
	win.SetContent(vBox)

}

func save(fileName string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

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


	for k, v := range data {
		val, isArray := isKey(k)
		if val {
			_, err := f.WriteString(fmt.Sprintf("%s:\n", k))
			if err != nil {
				fmt.Println(err)
				f.Close()
				return
			}

			if isArray {
				for i := 0; i < len(v.([]interface{})); i++ {
					nmap := v.([]interface{})[i].(map[string]interface{})
					for kk, vv := range nmap {
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
						_, err := f.WriteString(fmt.Sprintf("\tthe %s is %s\n", kk, farenTemp))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}
					} else if isSunVal(kk) {
						sunTime := time.Unix(int64(vv.(float64)), 0)
						_, err := f.WriteString(fmt.Sprintf("\tthe %s at %v\n", kk, sunTime))
						if err != nil {
							fmt.Println(err)
							f.Close()
							return
						}

					} else {
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
}

func isKey(k string) (ok bool, isArray bool) {
	isArray, ok = weatherKeys[k]
	return ok, isArray
}

func C2K(kelvin float64) string {
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

func isFeelsLike(k string) bool {
	return strings.Contains(k, "feels_like")
}

func isClouds(k string) bool {
	return strings.Contains(k, "all")
}