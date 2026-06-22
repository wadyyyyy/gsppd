// Package api handles user's requests from web page
package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type ContentItem struct {
	ObjType     string `json:"obj_type"`
	ObjName     string `json:"obj_name"`
	Description string `json:"description"`
}

type BossPoint struct {
	Name string  `json:"name"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

func HandleContent(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		fidParam := r.URL.Query().Get("fid")
		if fidParam == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"error": "отсутствует параметр fid"}`))
			return
		}

		rows, err := db.Query("SELECT obj_type, obj_name, description FROM content WHERE area_fid = ?", fidParam)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "ошибка базы данных"}`))
			return
		}
		defer rows.Close()

		var items []ContentItem
		for rows.Next() {
			var item ContentItem
			if err := rows.Scan(&item.ObjType, &item.ObjName, &item.Description); err != nil {
				continue
			}
			items = append(items, item)
		}

		if items == nil {
			items = []ContentItem{}
		}

		json.NewEncoder(w).Encode(items)
	}
}

type FeatureCollection struct {
	Type     string    `json:"type"`
	Features []Feature `json:"features"`
}

type Feature struct {
	Type       string          `json:"type"`
	Properties Properties      `json:"properties"`
	Geometry   json.RawMessage `json:"geometry"`
}

type Properties struct {
	Fid      int    `json:"fid"`
	AreaName string `json:"area_name"`
}

func HandleAreas(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT fid, name, AsGeoJSON(geom) FROM areas")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "ошибка базы данных"}`))
			return
		}
		defer rows.Close()

		collection := FeatureCollection{
			Type:     "FeatureCollection",
			Features: []Feature{},
		}

		for rows.Next() {
			var fid int
			var name string
			var geomStr string
			if err := rows.Scan(&fid, &name, &geomStr); err != nil {
				continue
			}

			feature := Feature{
				Type: "Feature",
				Properties: Properties{
					Fid:      fid,
					AreaName: name,
				},
				Geometry: json.RawMessage(geomStr),
			}
			collection.Features = append(collection.Features, feature)
		}

		json.NewEncoder(w).Encode(collection)
	}
}

func HandleBosses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		rows, err := db.Query("SELECT name, x, y FROM points")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "ошибка базы данных"}`))
			return
		}
		defer rows.Close()

		var items []BossPoint
		for rows.Next() {
			var item BossPoint
			if err := rows.Scan(&item.Name, &item.X, &item.Y); err != nil {
				continue
			}
			items = append(items, item)
		}

		if items == nil {
			items = []BossPoint{}
		}

		json.NewEncoder(w).Encode(items)
	}
}
