package main

import (
	"log"
	"strconv"
	"os"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

type User struct {
	gorm.Model
	FirstName        string
	LastName         string
	TelegramID       string
	TelegramUsername string
}
func init(){
	if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
	}
}

func main() {
	TelegramToken, exist := os.LookupEnv("TELEGRAM_BOT_TOKEN")
	if(!exist){
		log.Fatal("couldnt find any bot token")
	}
	bot, err := tgbotapi.NewBotAPI(TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	db, err := gorm.Open("postgres", "host=localhost port=5432 user=amir dbname=telecoin sslmode=disable")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&User{})

	bot.Debug = true
	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		db.Create(&User{FirstName: update.Message.From.FirstName, LastName: update.Message.From.LastName,
			TelegramID: strconv.Itoa(update.Message.From.ID), TelegramUsername: update.Message.From.UserName})
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
