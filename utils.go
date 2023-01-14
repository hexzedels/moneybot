package main

import (
	"regexp"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func newPurchase(item string, price float64, category string) *Purchase {
	p := &Purchase{
		item:     item,
		price:    price,
		category: category,
	}

	return p
}

func proccessPurchase(update tg.Update) (*Purchase, error) {
	r := regexp.MustCompile(`([^\d\*\&\%\$\#\@]+) (\d+[\,\.]{0,1}\d+) ([^\d\*\&\%\$\#\@]+)`)
	splitted := r.FindStringSubmatch(update.Message.Text)
	if len(splitted) > 1 {
		price, err := strconv.ParseFloat(splitted[2], 32)
		p := newPurchase(splitted[1], price, splitted[3])

		return p, err
	}
	return &Purchase{}, nil
}
