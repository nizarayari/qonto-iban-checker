package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

// BanksModel represents a table in postgres.
type BanksModel struct {
	gorm.Model
	BankCode string `gorm:"primary_key" json:"bank_code"`
	Bic      string `json:"bic"`
	Name     string `json:"name"`
}

// FormatsModel represents a table in postgres.
type FormatsModel struct {
	gorm.Model
	CountryCode    string `gorm:"primary_key" json:"country_code"`
	Country        string `json:"country"`
	IbanLength     int    `json:"iban_length"`
	BankCodeLength int    `json:"bank_code_length"`
}

func main() {
	setupDB()
	http.HandleFunc("/bic", handleIban)
	http.Handle("/favicon.ico", http.NotFoundHandler())

	err = http.ListenAndServe(":8080", nil)
	check(err)
}

func handleIban(w http.ResponseWriter, req *http.Request) {
	iban, ok := req.URL.Query()["iban"]

	if !ok || len(iban) < 1 {
		log.Println("Url Param 'iban' is missing")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jData, err := json.Marshal(getBic(iban[0]))
	check(err)
	w.Write(jData)

}

func setupDB() {
	connStr := "host=localhost user=nizarayari dbname=iban sslmode=disable password=postgres"

	db, err = gorm.Open("postgres", connStr)
	check(err)
	//defer db.Close()

	err = db.DB().Ping()
	check(err)

	db.DropTableIfExists(&BanksModel{})
	db.AutoMigrate(&BanksModel{})
	db.DropTableIfExists(&FormatsModel{})
	db.AutoMigrate(&FormatsModel{})

	db.CreateTable(&BanksModel{})
	db.CreateTable(&FormatsModel{})

	fillDbBanks()
	fillDbFormats()
}

func createTable(w http.ResponseWriter) {
	fmt.Fprintln(w, "CREATED TABLE ")
}

func fillDbBanks() {
	csvFile, _ := os.Open("banks.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		bank := BanksModel{
			BankCode: line[0],
			Bic:      line[1],
			Name:     line[2],
		}

		db.Create(&bank)

	}
}

func fillDbFormats() {
	csvFile, _ := os.Open("ibanRules.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		ibanLength, err := strconv.Atoi(line[2])
		check(err)
		bankCodeLength, err := strconv.Atoi(line[3])
		check(err)

		format := FormatsModel{
			Country:        line[0],
			CountryCode:    line[1],
			IbanLength:     ibanLength,
			BankCodeLength: bankCodeLength,
		}

		db.Create(&format)

	}
}

func check(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
