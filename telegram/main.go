package main

import (
	"bytes"
	"encoding/json"
	"errors"
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

type MessageResponse struct {
	Ok     bool          `json:"ok"`
	Result MessageResult `json:"result"`
}

type MessageResult struct {
	MessageId int64 `json:"message_id"`
}

type PinMessage struct {
	MessageId           int64 `json:"message_id"`
	ChatId              int64 `json:"chat_id"`
	DisableNotification bool  `json:"disable_notification"`
}

func SendMessage(url string, message *Message) (int64, error) {
	payload, err := json.Marshal(message)
	if err != nil {
		return -1, err
	}
	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return -1, err
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			log.Println("failed to close response body")
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return -1, fmt.Errorf("failed to send successful request. Status was %q", response.Status)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(response.Body)
	jsonString := string(buf.String())
	var j MessageResponse
	err = json.Unmarshal([]byte(jsonString), &j)
	if err != nil {
		panic(err)
	}
	return j.Result.MessageId, nil
}

func SendPinMessage(url string, pinMessage *PinMessage) error {
	payload, err := json.Marshal(pinMessage)
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

func UnpinAllMessages(url string, chatId int64) error {
	type unpinStruct struct {
		ChatId int64 `json:"chat_id"`
	}

	payload, err := json.Marshal(unpinStruct{ChatId: chatId})
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

	if _, err := os.Stat(config.WatcherFile); errors.Is(err, os.ErrNotExist) {
		url := config.TelegramUrl + config.TelegramToken + "/sendMessage"
		messageText, err := parser.GetMenuStringHTML(time.Now(), config.Url, config.Language)
		if err != nil {
			log.Fatal(err)
		}
		message := Message{ParseMode: "HTML", ChatID: config.TelegramChatId, Text: messageText}
		messageId, err := SendMessage(url, &message)
		if err != nil {
			log.Fatal(err)
		}

		err = UnpinAllMessages(config.TelegramUrl+config.TelegramToken+"/unpinAllChatMessages", config.TelegramChatId)
		if err != nil {
			log.Fatal(err)
		}

		err = SendPinMessage(config.TelegramUrl+config.TelegramToken+"/pinChatMessage", &PinMessage{messageId, config.TelegramChatId, true})
		if err != nil {
			log.Fatal(err)
		}

		os.Create(config.WatcherFile)

	} else {
		log.Fatal("already sent menu for today")
	}
}
