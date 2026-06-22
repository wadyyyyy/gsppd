# Interactive GIS Map Data Parser & API Server

A backend system and web API built in Go for parsing, processing, and serving geospatial and tabular data. This project processes GeoJSON polygons and Excel-based metadata, compiles them into a unified SQLite database, and serves the data via a RESTful API for interactive web mapping (applied to an Elden Ring DLC map).

This project highlights backend skills in data parsing, database design, and API development, matching real-world requirements for collecting and structuring complex datasets.

## Features

- **Geospatial Data Parsing:** Parses complex GeoJSON data (including Polygons and MultiPolygons with cutouts) into a relational database.
- **Tabular Data Processing:** Extracts and processes structured content from Excel files (`.xlsx`) using `excelize` to enrich geospatial areas with metadata.
- **Automated Database Assembly:** On startup, the system dynamically drops, initializes, and populates an SQLite database from the raw data sources.
- **RESTful API:** Serves static frontend assets and exposes API endpoints to retrieve map areas and associated content.
- **Data Visualization:** Includes a Python utility script utilizing `matplotlib` to quickly render and verify parsed GeoJSON shapes.

## Tech Stack

- **Backend:** Go (Golang) 1.25
- **Database:** SQLite3 (`github.com/mattn/go-sqlite3`)
- **Data Parsing:**
  - Excel: `github.com/xuri/excelize/v2`
  - JSON: Go standard library `encoding/json`
- **Scripting & Prototyping:** Python 3, Matplotlib (for GeoJSON visual preview)
- **GIS Software:** QGIS (for initial map data generation and `.qgz` project management)

## Project Structure

```text
gsppd/
├── data/                       # Processed data outputs
│   ├── ERSOTE_gsppd.xlsx       # Raw Excel data table
│   ├── areas.geojson           # Polygon geometry data
│   └── map_data.sqlite         # Generated SQLite database (created on runtime)
├── src/                        # Go source code
│   ├── api/                    # HTTP handlers and API logic
│   ├── db/                     # SQLite initialization and schema definitions
│   ├── parser/                 # Parsers for Excel and GeoJSON formats
│   └── main.go                 # Application entry point
├── static/                     # Frontend web assets (HTML/CSS/JS)
├── preview_geojson.py          # Python utility to plot GeoJSON geometry
└── gsppd.qgz                   # QGIS project file
```

## Getting Started

### Prerequisites

- Go (1.25+)
- Python 3.x (optional, for running the visualization script)
- GCC / CGO enabled (required for `go-sqlite3`)
- **SpatiaLite** extension for SQLite (required for parsing spatial data)
  - macOS: `brew install libspatialite`
  - Ubuntu/Debian: `sudo apt-get install libsqlite3-mod-spatialite`
  - Windows: Download the mod_spatialite binaries and add them to your PATH.

### Installation & Execution

1. **Clone the repository and install Go modules:**

   ```bash
   cd gsppd
   go mod tidy
   ```

2. **Run the Go Backend Server:**
   The application will automatically build the SQLite database from the GeoJSON and Excel files, then start the web server.

   ```bash
   cd src
   go run main.go
   ```

   The server will start at `http://localhost:8080`.

3. **Preview GeoJSON Data (Python):**
   If you want to visualize the raw map data without starting the web server, it is recommended to use a virtual environment:
   ```bash
   python3 -m venv .venv
   source .venv/bin/activate  # On Windows: .venv\Scripts\activate
   pip install matplotlib
   python preview_geojson.py
   ```

## API Endpoints

- `GET /` - Serves the interactive frontend map interface.
- `GET /api/areas` - Retrieves spatial geometry directly from the SQLite database (using SpatiaLite's `AsGeoJSON`) and serves it as a valid GeoJSON FeatureCollection.
- `GET /api/content` - Retrieves structured, parsed metadata from the SQLite database associated with the requested map areas.

## Relevancy to Data Engineering & Backend Roles

This project actively demonstrates:

- **Data Ingestion:** Reading from disparate data sources (JSON, Excel).
- **Data Normalization:** Merging spatial and tabular data into a strictly typed SQL schema.
- **Go Proficiency:** Implementing file I/O, routing, and database drivers in Go, making it an excellent showcase for high-load data parsing services.
