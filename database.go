package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	defaultBalance = 0.0
	defaultDebt    = 0.0
)

func initUsers(db *sql.DB) {
	users_table := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "FirstName" TEXT,
        "LastName" TEXT,
        "UserId" TEXT,
        "Balance" FLOAT,
		"Debt" FLOAT);`

	query, err := db.Prepare(users_table)
	if err != nil {
		log.Fatal("ERROR While preparing sql query: ", err)
	}
	query.Exec()
	log.Println("TABLE Users created successfully!")
}

func isInUsers(userId string, db *sql.DB) bool {
	isInUsers_select := fmt.Sprintf(`SELECT UserId FROM users WHERE UserId='%s'`, userId)
	var userIdR string
	row := db.QueryRow(isInUsers_select)
	err := row.Scan(&userIdR)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", row.Err())
	}
	if userIdR == userId {
		return true
	} else {
		return false
	}

}

func insertUser(update tg.Update, db *sql.DB) {
	user_insert := `INSERT INTO users (FirstName, LastName, UserId, Balance, Debt) VALUES (?, ?, ?, ?, ?);`
	query, err := db.Prepare(user_insert)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", err)
	}

	_, err = query.Exec(
		update.Message.From.FirstName,
		update.Message.From.LastName,
		strconv.Itoa(update.Message.From.ID),
		defaultBalance,
		defaultDebt)

	if err != nil {
		log.Println("ERROR While executing sql query: ", err)
	}
}
func initUserTable(userId string, db *sql.DB) {
	validUserId := fmt.Sprintf("user_%s", userId)
	user_table := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"Item" STRING,
	"Price" FLOAT,
	"Category" STRING);`, validUserId) // TODO: make Category int
	log.Println(user_table)
	query, err := db.Prepare(user_table)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", err)
	}
	_, err = query.Exec()
	if err != nil {
		log.Println("ERROR While executing sql query: ", err)
	} else {
		log.Println("TABLE &s created successfully!", userId)
	}
}

func insertToUserTable(userId string, data *Purchase, db *sql.DB) {
	validUserId := fmt.Sprintf("user_%s", userId)
	purchase_insert := fmt.Sprintf(`INSERT INTO %s(Item, Price, Category) VALUES (?, ?, ?);`, validUserId)
	query, err := db.Prepare(purchase_insert)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", err)
	}
	_, err = query.Exec(data.item, data.price, data.category)
	if err != nil {
		log.Println("ERROR While executing sql query: ", err)
	}
}
