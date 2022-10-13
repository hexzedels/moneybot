package main

import (
	"regexp"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const ZERO = 0.0

func proccessPurchase(update tg.Update) (*Purchase, error) {
	r, err := regexp.Compile(`([^\d\*\&\%\$\#\@]+) (\d+[\,\.]{0,1}\d+) ([^\d]+)`)
	splitted := r.FindStringSubmatch(update.Message.Text)
	p := Purchase{}
	if len(splitted) > 1 {
		price, err := strconv.ParseFloat(splitted[2], 32)
		p.item = splitted[1]
		p.price = price
		p.category = splitted[3]
		return &p, err
	}
	return &p, err
}
