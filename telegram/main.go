package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/petrvelicka/eurest/parser"
)

type Message struct {
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

func SendMessage(url string, message *Message) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}(response.Body)
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send successful request. Status was %q", response.Status)
	}
	return nil
}

func main() {
	config_path := "config.json"

	if len(os.Args) == 2 {
		config_path = os.Args[1]
	}

	config := ParseConfig(config_path)

	SendMessage(config.TelegramUrl+config.TelegramToken+"/sendMessage", &Message{ParseMode: "HTML", ChatID: config.TelegramChatId, Text: parser.GetMenuStringHTML(time.Now(), config.Url)})
}
