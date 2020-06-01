package dao

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"fmt"
)


func VerifyAccountValidity(account int, password string) bool {
	if account < 100000 || account > 999999 {
		return false
	}
	var dbname string
	if account <= 100100 {
		dbname = "database/SpecialAccount.db"
	} else if account <= 101101 {
		dbname = "database/TemporaryAccount.db"
	} else {
		dbname = "database/GeneralAccount.db"
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer db.Close()
	var pw string
	err = db.QueryRow("SELECT password FROM accounts WHERE account=?", account).Scan(&pw)
	if err == sql.ErrNoRows {
		return false
	}
	if pw == password {
		return true
	}
	return false
}