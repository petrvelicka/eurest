package parser

import "time"

type EurestMenu struct {
	Date    time.Time
	Soup    string
	Main    []string
	Dessert string
}
