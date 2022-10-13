package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	defaultBalance = 0.0
	defaultDebt    = 0.0
)

type Databse struct { // TODO: Update custom interface over this struct
	db *sql.DB
}

func (r *Databse) initUsers() {
	users_table := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "FirstName" TEXT,
        "LastName" TEXT,
        "UserId" TEXT,
        "Balance" FLOAT,
		"Debt" FLOAT);`

	query, err := r.db.Prepare(users_table)
	if err != nil {
		log.Fatal("ERROR While preparing sql query: ", err)
	}
	query.Exec()
	log.Println("TABLE Users created successfully!")
}

func (r *Databse) isInUsers(userId string) bool {
	isInUsers_select := fmt.Sprintf(`SELECT UserId FROM users WHERE UserId='%s'`, userId)
	var userIdR string
	row := r.db.QueryRow(isInUsers_select)
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

func (r *Databse) insertUser(userId string, update tg.Update) {
	if !r.isInUsers(userId) {
		user_insert := `INSERT INTO users (FirstName, LastName, UserId, Balance, Debt) VALUES (?, ?, ?, ?, ?);`
		query, err := r.db.Prepare(user_insert)
		if err != nil {
			log.Println("ERROR While preparing sql query: ", err)
		}

		_, err = query.Exec(
			update.Message.From.FirstName,
			update.Message.From.LastName,
			userId,
			defaultBalance,
			defaultDebt)

		if err != nil {
			log.Println("ERROR While executing sql query: ", err)
		}
	}
}
func (r *Databse) initUserTable(userId string) string {
	validUserId := fmt.Sprintf("user_%s", userId)
	if !r.ensureExists(validUserId) {
		user_table := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"Item" STRING,
			"Price" FLOAT,
			"Category" STRING);`, validUserId) // TODO: make Category int
		log.Println(user_table)
		query, err := r.db.Prepare(user_table)
		if err != nil {
			log.Println("ERROR While preparing sql query: ", err)
		}
		_, err = query.Exec()
		if err != nil {
			log.Println("ERROR While executing sql query: ", err)
		} else {
			log.Println("TABLE &s created successfully!", userId)
		}
		return SUCCESS_USER_TABLE_INIT
	} else {
		return EMPTY
	}
}

func (r *Databse) insertToUserTable(userId string, data *Purchase) { // TODO: Implement timestamp saving to sort by it later
	validUserId := fmt.Sprintf("user_%s", userId)
	purchase_insert := fmt.Sprintf(`INSERT INTO %s(Item, Price, Category) VALUES (?, ?, ?);`, validUserId)
	query, err := r.db.Prepare(purchase_insert)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", err)
	}
	_, err = query.Exec(data.item, data.price, data.category)
	if err != nil {
		log.Println("ERROR While executing sql query: ", err)
	}
}

// TODO: Think about wrapping work with rows into some over function
func (r *Databse) selectUserTable(userId string) []Purchase {
	validUserId := fmt.Sprintf("user_%s", userId)
	user_select := fmt.Sprintf(`SELECT Item, Price, Category FROM %s`, validUserId)
	query, err := r.db.Prepare(user_select)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", err)
	}
	rows, err := query.Query()
	if err != nil {
		log.Println("ERROR While executing sql query: ", err)
	}
	defer rows.Close()
	var data []Purchase

	for rows.Next() {
		i := Purchase{}
		err = rows.Scan(&i.item, &i.price, &i.category)
		if err != nil {
			log.Println("ERROR While scanning sql row: ", err)
		}
		data = append(data, i)
	}
	return data
}

func (r *Databse) selectSumUserTable(userId string) (float32, error) {
	validUserId := fmt.Sprintf("user_%s", userId)
	if r.ensureExists(validUserId) {
		user_select_sum := fmt.Sprintf(`SELECT SUM(Price) FROM %s;`, validUserId)
		query, err := r.db.Prepare(user_select_sum)
		if err != nil {
			log.Println("ERROR While preparing sql query: ", err)
		}
		rows, err := query.Query()
		if err != nil {
			log.Println("ERROR While executing sql query: ", err)
		}
		defer rows.Close()

		var sum float32
		for rows.Next() {
			err = rows.Scan(&sum)
			if err != nil {
				log.Println("ERROR While scanning sql row: ", err)
			}
		}
		return sum, nil
	} else {
		err := errors.New(fmt.Sprintf(USER_TABLE_IS_NOT_EXISTS, userId))
		return 0, err
	}
}

func (r *Databse) ensureExists(tableName string) bool {
	ensureQuery := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", tableName)
	var tableNameR string
	row := r.db.QueryRow(ensureQuery)
	err := row.Scan(&tableNameR)
	if err != nil {
		log.Println("ERROR While preparing sql query: ", row.Err())
	}
	if tableNameR == tableName {
		return true
	} else {
		return false
	}
}
