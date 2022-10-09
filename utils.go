package main

import (
	"strconv"
	"strings"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func proccessPurchase(update tg.Update) (*Purchase, error) {
	splitted := strings.Split(update.Message.Text, " ") // TODO: Implement splitting by the spaces separating only three fields
	price, err := strconv.ParseFloat(splitted[1], 32)
	p := Purchase{
		item:     splitted[0],
		price:    price,
		category: splitted[2]}

	return &p, err
}
