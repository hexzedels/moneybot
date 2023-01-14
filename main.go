package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
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
	envs, err := godotenv.Read(".env")
	if err != nil {
		panic(fmt.Sprintf("ERROR During loading .env file: %s", err))
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
		panic(fmt.Sprintf("while checking if base exists: %s", err))
	}

	if !isExists {
		_, err := os.Create("./data.db")
		if err != nil {
			log.Printf("ERROR Base is not created: %s", err)
		}
	}

	dataBase, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Printf("base is not created: %s", err)
	}
	defer dataBase.Close()

	db := Databse{dataBase}
	db.initUsers()

	bot, err := tg.NewBotAPI(TOKEN)
	if err != nil {
		log.Printf("bot is not created %s", err)
	}

	bot.Debug = true

	log.Printf("bot started as %s", bot.Self.UserName)

	u := tg.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(fmt.Sprintf("while receiving updates: %s", err))
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		userID := strconv.Itoa(update.Message.From.ID)
		if !update.Message.IsCommand() {
			msg := tg.NewMessage(update.Message.Chat.ID, "Записал!")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ParseMode = "HTML"

			purchase, err := proccessPurchase(update)
			if err != nil {
				log.Println("while processing purchase: ", err)
				msg.Text = "<b>Неверный формат!</b>\nИспользуйте /help"
			} else {
				if (Purchase{}) == *purchase {
					log.Println("received non spending msg: ", update.Message.Text)
					msg.Text = "<b>Неверный формат!</b>\nИспользуйте /help"
				} else {
					db.insertToUserTable(userID, purchase)
				}
			}

			_, err = bot.Send(msg)
			if err != nil {
				log.Println("while sending msg: ", err)
			}
		}

		switch update.Message.Command() {
		case "start":
			go start(bot, update, db, userID)
		case "ping":
			msg := tg.NewMessage(update.Message.Chat.ID, "pong")
			_, err := bot.Send(msg)
			if err != nil {
				log.Println("while sending msg: ", err)
			}
		case "list":
			go list(bot, update, db, userID)
		case "sum":
			go sum(bot, update, db, userID)
		// TODO: Implement keyboard ask where to add spending (trigger only for multiple tables)
		case "create": // TODO: Implement search for mentioned user in database and start goroutine that asks mentioned user is he
		// willing to join initiating user's fund
		case "help":
			go help(bot, update)
		default:
		}
	}
}
