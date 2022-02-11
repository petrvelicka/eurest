package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Config struct {
	Url            string
	TelegramToken  string
	TelegramUrl    string
	TelegramChatId int64
}

func ParseConfig(fname string) Config {
	content, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	var payload Config
	err = json.Unmarshal(content, &payload)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	return payload
}
