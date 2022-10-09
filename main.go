package main

import (
	"bytes"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

type Purchase struct {
	item     string
	price    float64
	category string
}

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func main() {
	// Loading .env file
	log.Println("PROCESSING .env file")
	envs, err := godotenv.Read(".env")

	if err != nil {
		log.Fatalf("ERROR During loading .env file: %s", err)
	}

	for key, value := range envs {
		if bytes.Contains([]byte(key), []byte("\xef\xbb\xbf")) {
			newkey := bytes.Trim([]byte(key), "\xef\xbb\xbf")
			delete(envs, key)
			envs[string(newkey)] = value
		}
	}
	TOKEN := envs["TOKEN"]

	// Databse init
	isExists, err := exists("./data.db")
	if err != nil {
		log.Fatal("ERROR While checking if base exists: ", err)
	}
	if !isExists {
		_, err := os.Create("./data.db")
		if err != nil {
			log.Fatal("ERROR Base is not created: ", err)
		}
	}

	db, err := sql.Open("sqlite3", "./data.db")
	defer db.Close()
	if err != nil {
		log.Fatal("ERROR Base is not created: ", err)
	}
	initUsers(db)

	bot, err := tg.NewBotAPI(TOKEN)
	if err != nil {
		log.Fatal("ERROR Bot is not created", err)

	}

	bot.Debug = true

	log.Printf("AUTHORIZED on account %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal("ERROR while receiving updates: ", err)
	}

	for update := range updates {
		userId := strconv.Itoa(update.Message.From.ID)
		if update.Message == nil {
			continue
		}
		if !update.Message.IsCommand() {
			purchase, err := proccessPurchase(update)
			if err != nil {
				log.Println("ERROR While proccessing purchase: ", err)
			}
			insertToUserTable(userId, purchase, db)
			msg := tg.NewMessage(update.Message.Chat.ID, "Записал!")
			msg.ReplyToMessageID = update.Message.MessageID
			_, err = bot.Send(msg)
			if err != nil {
				log.Println("ERROR While sending msg: ", err)
			}
		}
		switch update.Message.Command() {
		case "start":

			initUserTable(userId, db) // TODO: Modify msg.text if table already created
			if !isInUsers(userId, db) {
				insertUser(update, db)
			}
			msg := tg.NewMessage(update.Message.Chat.ID, "Таблица успешно создана!")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("ERROR While sending msg: ", err)
			}
		case "ping":
			msg := tg.NewMessage(update.Message.Chat.ID, "pong")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("ERROR While sending msg: ", err)
			}
		case "keyboard":
		case "help":
			msg := tg.NewMessage(update.Message.Chat.ID, "Use /start to create your own money tracker\n Play ping-pong with /ping")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("ERROR While sending msg: ", err)
			}
		default:

		}
	}

}
