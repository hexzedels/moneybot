package main

import (
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jedib0t/go-pretty/table"
)

func start(bot *tg.BotAPI, update tg.Update, db Databse, userId string) {
	msg := tg.NewMessage(update.Message.Chat.ID, "")
	msg.Text = db.initUserTable(userId) // TODO: Modify msg.text if table already created
	db.insertUser(userId, update)
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("ERROR While sending msg: ", err)
	}
}

func sum(bot *tg.BotAPI, update tg.Update, db Databse, userId string) {
	sum, err := db.selectSumUserTable(userId)
	errD := fmt.Sprintf(UserTableIsNotExists, userId)
	if err != nil && err.Error() == errD {
		db.initUserTable(userId)
	}

	msg := tg.NewMessage(update.Message.Chat.ID, "")
	msg.ParseMode = "HTML"
	msg.Text = fmt.Sprintf("Sum of last spendings is: %.2f", sum)
	_, err = bot.Send(msg)
	if err != nil {
		log.Println("ERROR While sending msg: ", err)
	}
}

func help(bot *tg.BotAPI, update tg.Update) {
	msg := tg.NewMessage(update.Message.Chat.ID, HelpString)
	msg.ParseMode = "HTML"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("ERROR While sending msg: ", err)
	}
}

func list(bot *tg.BotAPI, update tg.Update, db Databse, userId string) {
	data := db.selectUserTable(userId) //TODO: Add table creation
	t := table.NewWriter()
	t.AppendHeader(table.Row{"#", "Item", "Price", "Category"})
	for i, e := range data {
		t.AppendRow([]interface{}{i, e.item, e.price, e.category})
	}

	msg := tg.NewMessage(update.Message.Chat.ID, t.RenderMarkdown())
	msg.ParseMode = "Markdown"
	_, err := bot.Send(msg)
	if err != nil {
		log.Println("ERROR While sending msg: ", err)
	}
}

func add(bot *tg.BotAPI, update tg.Update, db Databse, userId string) {

}
