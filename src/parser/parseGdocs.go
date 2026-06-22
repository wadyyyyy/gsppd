// Package parser parses spreadsheet from excel to SQL
package parser

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func parseFIDs(fidStr string) []int {
	var fids []int
	parts := strings.Split(fidStr, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if strings.Contains(part, "-") {
			rangeParts := strings.Split(part, "-")
			if len(rangeParts) == 2 {
				start, err1 := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
				end, err2 := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
				if err1 == nil && err2 == nil && start <= end {
					for i := start; i <= end; i++ {
						fids = append(fids, i)
					}
				}
			}
		} else {
			if id, err := strconv.Atoi(part); err == nil {
				fids = append(fids, id)
			}
		}
	}
	return fids
}

func parseCoords(coordStr string) (float64, float64, error) {
	parts := strings.Split(coordStr, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("неверный формат координат: %s", coordStr)
	}

	xStr := strings.TrimSpace(parts[0])
	yStr := strings.TrimSpace(parts[1])
	if xStr == "" || yStr == "" {
		return 0, 0, fmt.Errorf("пустые координаты: %s", coordStr)
	}

	x, err := strconv.ParseFloat(xStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга X %q: %v", xStr, err)
	}
	y, err := strconv.ParseFloat(yStr, 64)
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга Y %q: %v", yStr, err)
	}

	return x, y, nil
}

func ParseGdocs(db *sql.DB, filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return fmt.Errorf("ошибка открытия excel: %v", err)
	}
	defer f.Close()

	stmt, err := db.Prepare(`
		INSERT INTO content (area_fid, obj_type, obj_name, description)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки sql-запроса: %v", err)
	}
	defer stmt.Close()

	rowsItems, err := f.GetRows("items")
	if err == nil {
		for i, row := range rowsItems {
			if i == 0 {
				continue
			}
			if len(row) < 6 {
				continue
			}

			itemName := row[1]
			itemDesc := row[2]
			areaFid := row[3]

			fullDesc := fmt.Sprintf("%s (Кол-во: %s). Место: %s", itemDesc, row[4], row[5])

			fids := parseFIDs(areaFid)
			for _, fid := range fids {
				_, err = stmt.Exec(fid, "item", itemName, fullDesc)
				if err != nil {
					fmt.Printf("Ошибка вставки предмета %s для зоны %d: %v\n", itemName, fid, err)
				}
			}
		}
		fmt.Println("Предметы успешно загружены!")
	} else {
		fmt.Println("Не удалось прочитать лист (item):", err)
	}

	rowsEnemies, err := f.GetRows("enemies")
	if err == nil {
		for i, row := range rowsEnemies {
			if i == 0 {
				continue
			}
			if len(row) < 5 {
				continue
			}

			enemyName := row[1]
			enemyDesc := row[2]
			areaFid := row[3]

			fullDesc := fmt.Sprintf("%s. Дроп: %s", enemyDesc, row[4])

			fids := parseFIDs(areaFid)
			for _, fid := range fids {
				_, err = stmt.Exec(fid, "enemy", enemyName, fullDesc)
				if err != nil {
					fmt.Printf("Ошибка вставки врага %s для зоны %d: %v\n", enemyName, fid, err)
				}
			}
		}
		fmt.Println("Враги успешно загружены!")
	} else {
		fmt.Println("Не удалось прочитать лист (enemies):", err)
	}

	return nil
}

func ParsePoints(db *sql.DB, filepath string) error {
	f, err := excelize.OpenFile(filepath)
	if err != nil {
		return fmt.Errorf("ошибка открытия excel: %v", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("в excel-файле нет листов")
	}
	sheetName := sheets[0]

	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("ошибка чтения листа %s: %v", sheetName, err)
	}

	stmt, err := db.Prepare(`
		INSERT INTO points (name, x, y, geom)
		VALUES (?, ?, ?, GeomFromText(?, 0))
	`)
	if err != nil {
		return fmt.Errorf("ошибка подготовки sql-запроса: %v", err)
	}
	defer stmt.Close()

	for i, row := range rows {
		if i == 0 {
			continue
		}
		if len(row) < 2 {
			continue
		}

		name := strings.TrimSpace(row[0])
		coordStr := strings.TrimSpace(row[1])
		if name == "" || coordStr == "" {
			continue
		}

		x, y, err := parseCoords(coordStr)
		if err != nil {
			fmt.Printf("Ошибка парсинга координат %q для %s: %v\n", coordStr, name, err)
			continue
		}

		wkt := fmt.Sprintf("POINT(%f %f)", x, y)
		_, err = stmt.Exec(name, x, y, wkt)
		if err != nil {
			fmt.Printf("Ошибка вставки точки %s: %v\n", name, err)
		}
	}

	fmt.Println("Точки успешно загружены!")
	return nil
}
