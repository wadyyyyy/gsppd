// Package parser parses different data to SQL
package parser

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type FeatureCollection struct {
	Features []Feature `json:"features"`
}

type Feature struct {
	Properties struct {
		Fid      int    `json:"fid"`
		AreaName string `json:"area_name "`
	} `json:"properties"`
	Geometry json.RawMessage `json:"geometry"`
}

func ParseGeoJSON(db *sql.DB, filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	byteValue, _ := io.ReadAll(file)

	var collection FeatureCollection
	if err := json.Unmarshal(byteValue, &collection); err != nil {
		return fmt.Errorf("ошибка парсинга JSON: %v", err)
	}

	stmt, err := db.Prepare(`
		INSERT INTO areas (fid, name, geom)
		VALUES (?, ?, GeomFromGeoJSON(?))
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки запроса: %v", err)
	}
	defer stmt.Close()

	for _, feature := range collection.Features {
		geomStr := string(feature.Geometry)
		_, err := stmt.Exec(feature.Properties.Fid, feature.Properties.AreaName, geomStr)
		if err != nil {
			fmt.Printf("Ошибка вставки зоны %s: %v\n", feature.Properties.AreaName, err)
			continue
		}
	}

	fmt.Println("GeoJSON успешно загружен в БД!")
	return nil
}
