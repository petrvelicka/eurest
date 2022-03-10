package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/petrvelicka/eurest/parser"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func GetMenuStringHTML(day time.Time, url string, language string) (string, error) {
	menu, err := parser.GetDay(day, parser.ParseCSV(url))
	if err != nil {
		return "", err
	}

	formatString := "%s %s:\n\n<b>%s</b>: %s\n\n<b>%s 1</b>: %s\n\n<b>%s 2</b>: %s\n\n<b>%s 3</b>: %s\n\n<b>%s</b>: %s\n\n%s"
	resultString := ""
	switch language {
	case "cs":
		resultString = fmt.Sprintf(formatString, "Menu pro", menu.Date.Format("02.01.2006"), "Polévka", menu.Soup, "Jídlo", menu.Main[0], "Jídlo", menu.Main[1], "Jídlo", menu.Main[2], "Dezert", menu.Dessert, "Dobrou chuť!")
		break
	case "en":
		fallthrough
	default:
		resultString = fmt.Sprintf(formatString, "Menu for", menu.Date.Format("2006-01-02"), "Soup", menu.Soup, "Meal", menu.Main[0], "Meal", menu.Main[1], "Meal", menu.Main[2], "Dessert", menu.Dessert, "Enjoy your meal!")
	}

	return resultString, nil
}

func main() {
	config_path := "config.json"

	if len(os.Args) == 2 {
		config_path = os.Args[1]
	}

	config := ParseConfig(config_path)

	if _, err := os.Stat(config.WatcherFile); errors.Is(err, os.ErrNotExist) {

		bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
		if err != nil {
			log.Panic(err)
		}

		msgText, err := GetMenuStringHTML(time.Now(), config.EurestUrl, config.Language)
		if err != nil {
			log.Panic(err)
		}

		msg := tgbotapi.NewMessage(config.Telegram.ChatId, msgText)
		msg.ParseMode = "HTML"
		response, err := bot.Send(msg)
		if err != nil {
			log.Panic(err)
		}
		messageId := response.MessageID
		unPinAll := tgbotapi.UnpinAllChatMessagesConfig{
			ChatID: config.Telegram.ChatId,
		}
		_, err = bot.Request(unPinAll)
		if err != nil {
			log.Panic(err)
		}

		pinMessage := tgbotapi.PinChatMessageConfig{
			ChatID:              config.Telegram.ChatId,
			MessageID:           messageId,
			DisableNotification: true,
		}
		_, err = bot.Request(pinMessage)
		if err != nil {
			log.Panic(err)
		}

		os.Create(config.WatcherFile)

	} else {
		log.Fatal("already sent menu for today")
	}
}
