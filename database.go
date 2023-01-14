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
	usersTable := `CREATE TABLE IF NOT EXISTS users (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "FirstName" TEXT,
        "LastName" TEXT,
        "UserId" TEXT,
        "Balance" FLOAT,
		"Debt" FLOAT);`

	query, err := r.db.Prepare(usersTable)
	if err != nil {
		panic(fmt.Sprintf("while preparing sql query: ", err))
	}

	_, err = query.Exec()
	if err != nil {
		panic(err)
	}

	log.Println("table users created successfully!")
}

func (r *Databse) isInUsers(userID string) bool {
	isInUsersSelect := fmt.Sprintf(`SELECT UserId FROM users WHERE UserId='%s'`, userID)

	var userIDR string

	row := r.db.QueryRow(isInUsersSelect)

	if err := row.Scan(&userIDR); err != nil {
		log.Println("while preparing sql query: ", row.Err())
	}

	if userIDR == userID {
		return true
	}

	return false
}

func (r *Databse) insertUser(userID string, update tg.Update) {
	if !r.isInUsers(userID) {
		userInsert := `INSERT INTO users (FirstName, LastName, UserId, Balance, Debt) VALUES (?, ?, ?, ?, ?);`

		query, err := r.db.Prepare(userInsert)
		if err != nil {
			log.Println("while preparing sql query: ", err)
		}

		_, err = query.Exec(
			update.Message.From.FirstName,
			update.Message.From.LastName,
			userID,
			defaultBalance,
			defaultDebt)

		if err != nil {
			log.Println("while executing sql query: ", err)
		}
	}
}

func (r *Databse) initUserTable(userID string) string {
	validUserID := fmt.Sprintf("user_%s", userID)
	if !r.ensureExists(validUserID) {
		userTable := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s(
			id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
			"Item" STRING,
			"Price" FLOAT,
			"Category" STRING);`, validUserID) // TODO: make Category int
		log.Println(userTable)

		query, err := r.db.Prepare(userTable)
		if err != nil {
			log.Println("while preparing sql query: ", err)
		}

		_, err = query.Exec()
		if err != nil {
			log.Println("while executing sql query: ", err)
		} else {
			log.Printf("table %s created successfully!", userID)
		}

		return SuccessUserTableInit
	}

	return EmptyString
}

func (r *Databse) insertToUserTable(userID string, data *Purchase) { // TODO: Implement timestamp saving to sort by it later
	validUserID := fmt.Sprintf("user_%s", userID)
	purchaseInsert := fmt.Sprintf(`INSERT INTO %s(Item, Price, Category) VALUES (?, ?, ?);`, validUserID)

	query, err := r.db.Prepare(purchaseInsert)
	if err != nil {
		log.Println("while preparing sql query: ", err)
	}

	_, err = query.Exec(data.item, data.price, data.category)
	if err != nil {
		log.Println("while executing sql query: ", err)
	}
}

// TODO: Think about wrapping work with rows into some over function
func (r *Databse) selectUserTable(userID string) []Purchase {
	validUserID := fmt.Sprintf("user_%s", userID)
	userSelect := fmt.Sprintf(`SELECT Item, Price, Category FROM %s`, validUserID)

	query, err := r.db.Prepare(userSelect)
	if err != nil {
		log.Println("while preparing sql query: ", err)
	}

	rows, err := query.Query()
	if err != nil {
		log.Println("while executing sql query: ", err)
	}

	defer rows.Close()

	var data []Purchase

	for rows.Next() {
		i := Purchase{}

		if err = rows.Scan(&i.item, &i.price, &i.category); err != nil {
			log.Println("while scanning sql row: ", err)
		}

		data = append(data, i)
	}

	return data
}

func (r *Databse) selectSumUserTable(userID string) (float32, error) {
	validUserID := fmt.Sprintf("user_%s", userID)
	if r.ensureExists(validUserID) {
		userSelectSum := fmt.Sprintf(`SELECT SUM(Price) FROM %s;`, validUserID)

		query, err := r.db.Prepare(userSelectSum)
		if err != nil {
			log.Println("while preparing sql query: ", err)
		}

		rows, err := query.Query()
		if err != nil {
			log.Println("while executing sql query: ", err)
		}

		defer rows.Close()

		var sum float32

		for rows.Next() {
			err = rows.Scan(&sum)
			if err != nil {
				log.Println("while scanning sql row: ", err)
			}
		}

		return sum, nil
	}

	return 0, errors.New(fmt.Sprintf(UserTableIsNotExists, userID))
}

func (r *Databse) ensureExists(tableName string) bool {
	ensureQuery := fmt.Sprintf("SELECT name FROM sqlite_master WHERE type='table' AND name='%s';", tableName)

	var tableNameR string

	row := r.db.QueryRow(ensureQuery)

	if err := row.Scan(&tableNameR); err != nil {
		log.Println("while preparing sql query: ", row.Err())
	}

	if tableNameR == tableName {
		return true
	}

	return false
}
