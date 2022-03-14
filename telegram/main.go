package main

import (
	"errors"
	"flag"
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

	var lang *Language
	switch language {
	case "cs":
		lang = &LanguageCzech
		break
	case "en":
		fallthrough
	default:
		lang = &LanguageEnglish
	}
	resultString := fmt.Sprintf(formatString, lang.MenuFor, menu.Date.Format(lang.DateFormat), lang.Soup, menu.Soup, lang.Meal, menu.Main[0], lang.Meal, menu.Main[1], lang.Meal, menu.Main[2], lang.Dessert, menu.Dessert, lang.EnjoyYourMeal)

	return resultString, nil
}

func main() {
	var configPath string
	var skipDayCheck bool
	flag.StringVar(&configPath, "config", "config.json", "path to config file")
	flag.BoolVar(&skipDayCheck, "skipdaycheck", false, "allow sending more than once per day")
	flag.Parse()
	log.Printf("Parsing config file %s\n", configPath)
	config := ParseConfig(configPath)

	if _, err := os.Stat(config.WatcherFile); errors.Is(err, os.ErrNotExist) || skipDayCheck {
		bot, err := tgbotapi.NewBotAPI(config.Telegram.Token)
		if err != nil {
			log.Fatal(err)
		}

		msgText, err := GetMenuStringHTML(time.Now(), config.EurestUrl, config.Language)
		if err != nil {
			log.Fatal(err)
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
