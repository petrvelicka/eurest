package parser

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

func downloadCSV(url string) [][]string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("error downloading csv")
	}

	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)

	csvReader := csv.NewReader(buf)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("error parsing csv")
	}
	return data
}

func buildURL(code string) string {
	return "https://docs.google.com/spreadsheets/d/" + code + "/export?format=csv"
}

func ParseCSV(code string) []EurestMenu {
	url := buildURL(code)
	log.Println("downloading from " + url)
	csv := downloadCSV(url)

	offset := 5
	menu := []EurestMenu{}

	for day := 0; day < 5; day += 1 {
		date, err := time.Parse("2.1.2006", csv[offset+1][0])
		if err != nil {
			log.Fatal("error parsing date")
		}

		menu = append(menu, EurestMenu{
			Date:    date,
			Soup:    csv[offset][2],
			Main:    []string{csv[offset+2][2], csv[offset+3][2], csv[offset+4][2]},
			Dessert: csv[offset+1][2],
		})

		offset += 7
	}

	return menu
}

func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func GetDay(day time.Time, menu []EurestMenu) (EurestMenu, error) {
	for _, dayMenu := range menu {
		if DateEqual(day, dayMenu.Date) {
			if dayMenu.Soup != "" {
				return dayMenu, nil
			}
			break
		}
	}
	return EurestMenu{}, errors.New("date not found")
}

func GetTodayMenu(url string) (EurestMenu, error) {
	menu := ParseCSV(url)

	fmt.Println(menu)

	today := time.Now()
	todayMenu, err := GetDay(today, menu)
	if err != nil {
		return EurestMenu{}, errors.New("day not found. Nejspíš jsou prázdniny...")
	}

	return todayMenu, nil
}

func GetMenuString(day time.Time, url string) string {
	menu, err := GetDay(day, ParseCSV(url))
	if err != nil {
		return err.Error()
	}
	return menu.Date.Format("2006-01-02") + ":\n\n**Soup**: " + menu.Soup + "\n\n**Meal 1**: " + menu.Main[0] + "\n\n**Meal 2**: " + menu.Main[1] + "\n\n**Meal 3**: " + menu.Main[2] + "\n\n**Dessert**: " + menu.Dessert
}
