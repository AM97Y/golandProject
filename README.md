# goland yar
Данное приложение выводит погоду в Ярославле и позволяет сохранять расширенные отчеты по указаному пути.

ЗАПУСК

go get fyne.io/fyne

go get github.com/briandowns/openweathermap

go build

go run golandProject

ЛИБО

./gradlew build
Результат в папке .gogradle



Так же по желанию можно добавить иконку.
go get github.com/akavel/rsrc
rsrc -ico YOUR_ICON_FILE_NAME.ico
go build
