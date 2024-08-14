package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type Rate struct {
	Usdbrl struct {
		Code       string `json:"code"`
		Codein     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type DollarRateResponse struct {
	Bid string `json:"bid"`
}

func main() {
	sqliteDb, err := sql.Open("sqlite3", "./dollar_rate.db")
	if err != nil {
		log.Fatal(err)
	}
	db = sqliteDb
	initDb()
	defer db.Close()
	http.HandleFunc("/", getDollarRate)
	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func initDb() {
	const create string = `
  CREATE TABLE IF NOT EXISTS dollar_rates (
  id INTEGER NOT NULL PRIMARY KEY,
  time DATETIME NOT NULL,
  rate varchar(255) NOT NULL
  );`
	_, err := db.Exec(create)
	if err != nil {
		log.Fatal(err)
	}
}

func saveRateInDb(rate string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		return fmt.Errorf("timeout Saving in the database")
	default:
		const insert string = `INSERT INTO dollar_rates (time, rate) VALUES (?, ?);`

		stmt, err := db.Prepare(insert)
		if err != nil {
			return fmt.Errorf("error preparing statement: %v", err)
		}
		defer stmt.Close()
		_, err = stmt.Exec(time.Now(), rate)
		if err != nil {
			return fmt.Errorf("error executing statement: %v", err)
		}
		return nil
	}
}

func getDollarRate(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()

	select {
	case <-ctx.Done():
		http.Error(w, "Timeout getting the dollar rate", http.StatusRequestTimeout)
		return
	default:
		response, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		rateType := Rate{}
		err = json.Unmarshal(body, &rateType)
		if err != nil {
			log.Fatal(err)
		}

		err = saveRateInDb(rateType.Usdbrl.Bid)
		if err != nil {
			http.Error(w, "Timeout saving in the database", http.StatusRequestTimeout)
			return
		}

		res, err := json.Marshal(DollarRateResponse{Bid: rateType.Usdbrl.Bid})
		if err != nil {
			http.Error(w, "Error marshaling response", http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(res)
	}
}
