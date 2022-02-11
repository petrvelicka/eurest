package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/petrvelicka/eurest/parser"
)

var BotId string
var goBot *discordgo.Session
var config Config

func main() {
	config = ParseConfig("./config.json")
	goBot, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Making our bot a user using User function .
	u, err := goBot.User("@me")
	//Handlinf error
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Storing our id from u to BotId .
	BotId = u.ID

	// Adding handler function to handle our messages using AddHandler from discordgo package. We will declare messageHandler function later.
	goBot.AddHandler(messageHandler)

	err = goBot.Open()
	//Error handling
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	//If every thing works fine we will be printing this.
	fmt.Println("Bot is running !")
	<-make(chan struct{})
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Bot musn't reply to it's own messages , to confirm it we perform this check.
	if m.Author.ID == BotId {
		return
	}
	//If we message ping to our bot in our discord it will return us pong .
	if m.Content == "/menu" {
		_, _ = s.ChannelMessageSend(m.ChannelID, parser.GetMenuString(time.Now(), config.Url))
	}
}
