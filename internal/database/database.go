package database

import (
	"engractice/internal/models"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type Database struct {
	sheet *sheets.Service
}

var (
	spreadsheetId = os.Getenv("SHEET_ID")    //"1_xKMjnfCG3ADEH5nz5JOqvsFsdQ7UVPmc2ZDBtpvoc8"
	rangeData     = os.Getenv("SHEET_RANGE") //"vocabulary!A2:E"
	//database = os.Getenv("BLUEPRINT_DB_DATABASE")
	credentialsFilePath = os.Getenv("CREADENTIALS_FILE_PATH") // "credentials.json"
)

func NewDatabase() *Database {
	service, err := sheets.NewService(nil, option.WithCredentialsFile(credentialsFilePath))
	if err != nil {
		log.Fatalf("Unable to create Sheets client: %v", err)
	}
	return &Database{
		sheet: service,
	}
}

func (t *Database) GetSpreadsheetData() ([]models.Vocabulary, error) {
	resp, err := t.sheet.Spreadsheets.Values.Get(spreadsheetId, rangeData).Do()
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve data from sheet: %v", err)
	}
	if len(resp.Values) == 0 {
		return nil, fmt.Errorf("no data found in the specified range")
	}
	words, err := t.parseSheetData(resp.Values)
	return words, err
}

func (t *Database) UpdateSpreadsheetData(words []models.Vocabulary) error {
	var vr sheets.ValueRange
	for _, word := range words {
		vr.Values = append(vr.Values, []interface{}{
			word.English,
			word.Vietnamese,
			word.MP3,
			word.Tag,
			word.Point,
		})
	}
	_, err := t.sheet.Spreadsheets.Values.Update(spreadsheetId, rangeData, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		log.Fatalf("Unable to update data in sheet: %v", err)
	}
	return nil
}

func (t *Database) parseSheetData(data [][]interface{}) ([]models.Vocabulary, error) {
	var words []models.Vocabulary
	for _, row := range data {
		if len(row) < 5 {
			continue
		}
		point := 0
		if p, ok := row[4].(string); ok {
			fmt.Sscanf(p, "%d", &point)
		}
		words = append(words, models.Vocabulary{
			English:    fmt.Sprintf("%v", row[0]),
			Vietnamese: fmt.Sprintf("%v", row[1]),
			MP3:        fmt.Sprintf("%v", row[2]),
			Tag:        fmt.Sprintf("%v", row[3]),
			Point:      point,
		})
	}
	if len(words) == 0 {
		panic("No valid words found in the sheet data")
	}
	return words, nil
}
