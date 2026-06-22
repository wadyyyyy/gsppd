package main

import (
	"fmt"
	"net/http"
	"os"

	"gsppd/src/api"
	"gsppd/src/db"
	"gsppd/src/parser"
)

func main() {
	fmt.Println("Запускаем сборку базы данных...")
	dbPath := "../data/map_data.sqlite"
	os.Remove(dbPath)

	database, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Printf("Ошибка инициализации БД: %v\n", err)
		return
	}
	defer database.Close()

	schemaBytes, err := os.ReadFile("db/schema.sql")
	if err != nil {
		fmt.Printf("Ошибка чтения schema.sql: %v\n", err)
		return
	}

	_, err = database.Exec(string(schemaBytes))
	if err != nil {
		fmt.Printf("Ошибка создания таблиц: %v\n", err)
		return
	}
	fmt.Println("Таблицы успешно созданы.")

	err = parser.ParseGeoJSON(database, "../data/areas.geojson")
	if err != nil {
		fmt.Println("Сбой при парсинге GeoJSON:", err)
	}

	err = parser.ParseGdocs(database, "../data/ERSOTE_gsppd.xlsx")
	if err != nil {
		fmt.Println("Сбой при парсинге Excel:", err)
	}

	pointsPath := "../data/bosses.xlsx"
	if _, err := os.Stat(pointsPath); err == nil {
		err = parser.ParsePoints(database, pointsPath)
		if err != nil {
			fmt.Println("Сбой при парсинге точек:", err)
		}
	}

	fmt.Println("Сборка базы данных успешно завершена!")

	fmt.Println("Запускаем веб-сервер на http://localhost:8080 ...")
	fs := http.FileServer(http.Dir("../static"))
	http.Handle("/", fs)

	assets := http.FileServer(http.Dir("../assets"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	http.HandleFunc("/api/areas", api.HandleAreas(database))
	http.HandleFunc("/api/content", api.HandleContent(database))
	http.HandleFunc("/api/bosses", api.HandleBosses(database))
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
