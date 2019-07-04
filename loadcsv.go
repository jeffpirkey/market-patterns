package main

import (
	"encoding/csv"
	log "github.com/sirupsen/logrus"
	"market-patterns/model"
	"strconv"
	"time"
)

const (
	timeFormat = "2006-01-02"
)

func load(tsym string, r *csv.Reader) {

	vals, err := r.ReadAll()
	if err != nil {
		log.Fatalf("error reading csv due to %v", err)
	}

	ticker := Tickers.Find(tsym)
	for i, v := range vals {

		if i == 0 {
			// skip header line
			continue
		}

		date := convertTime(v[0])
		open := convertFloat(v[1])
		high := convertFloat(v[2])
		low := convertFloat(v[3])
		cl := convertFloat(v[4])
		volume := convertInt(v[5])

		if err != nil {
			log.Fatalf("error parsing csv time due ot %v", err)
		}

		p := model.Period{Date: date, Open: open, High: high, Low: low, Close: cl, Volume: volume}

		ticker.AddPeriod(&p)
	}
}

func convertFloat(v string) float64 {
	f, err := strconv.ParseFloat(v, 64)
	if err != nil {
		log.Errorf("unable to convert csv value to float due to %v", err)
	}
	return f
}

func convertInt(v string) int {
	f, err := strconv.Atoi(v)
	if err != nil {
		log.Errorf("unable to convert csv value to int due to %v", err)
	}
	return f
}

func convertTime(v string) time.Time {
	t, err := time.Parse(timeFormat, v)
	if err != nil {
		log.Errorf("unable to convert csv value to time due to %v", err)
	}
	return t
}