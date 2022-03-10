package main

type Language struct {
	MenuFor       string
	DateFormat    string
	Soup          string
	Meal          string
	Dessert       string
	EnjoyYourMeal string
}

var LanguageEnglish = Language{
	MenuFor:       "Menu for",
	DateFormat:    "2006/01/02",
	Soup:          "Soup",
	Meal:          "Meal",
	Dessert:       "Dessert",
	EnjoyYourMeal: "Enjoy your meal!",
}

var LanguageCzech = Language{
	MenuFor:       "Menu pro",
	DateFormat:    "02.01.2006",
	Soup:          "Polévka",
	Meal:          "Jídlo",
	Dessert:       "Dezert",
	EnjoyYourMeal: "Dobrou chuť!",
}
