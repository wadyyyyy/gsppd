// Package db initializes qslite database with spacialite extension
package db

import (
	"database/sql"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

func init() {
	sql.Register("spatialite", &sqlite3.SQLiteDriver{
		ConnectHook: func(conn *sqlite3.SQLiteConn) error {
			paths := []string{
				"mod_spatialite",                              // Поиск в стандартных путях системы (PATH / LD_LIBRARY_PATH)
				"mod_spatialite.so",                           // Linux
				"mod_spatialite.dylib",                        // macOS
				"mod_spatialite.dll",                          // Windows
				"/opt/homebrew/lib/mod_spatialite.dylib",      // macOS ARM (Apple Silicon)
				"/usr/local/lib/mod_spatialite.dylib",         // macOS Intel
				"/usr/lib/x86_64-linux-gnu/mod_spatialite.so", // Ubuntu / Debian
			}

			var err error
			for _, p := range paths {
				err = conn.LoadExtension(p, "sqlite3_modspatialite_init")
				if err == nil {
					return nil
				}
			}
			return fmt.Errorf("не удалось найти mod_spatialite. Последняя ошибка: %v", err)
		},
	})
}

func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("spatialite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ошибка загрузки библиотеки SpatiaLite: %v", err)
	}

	return db, nil
}
