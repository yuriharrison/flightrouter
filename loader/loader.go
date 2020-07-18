package loader

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	db "github.com/yuriharrison/bexs-test/flightsdb"
)

const bufferSize = 4096

// ImportFlightsFromFile import flights into FlightsDB
func ImportFlightsFromFile(fileName string, db *db.FlightsDB) error {
	f, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	err = LoadFlights(f, db)
	if err != nil {
		return err
	}
	return nil
}

// LoadFlights load flights into FlightsDB
func LoadFlights(file io.Reader, db *db.FlightsDB) error {
	csvReader := csv.NewReader(file)
	for {
		row, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		if len(row) != 3 {
			return fmt.Errorf("Expected 3 columns in CSV got \"%v\"", len(row))
		}
		origCode, destCode, priceStr := row[0], row[1], row[2]
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			err = db.Add(origCode, destCode, price)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
