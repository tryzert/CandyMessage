package dao
import (
	"database/sql"
	"log"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"can/tf"
)


/*
	用于账号注册服务。因为管理员账户和临时账户已经全部注册完毕，对外开放的仅仅是普通账户
*/


const (
	SPECIAL_ACCOUNT = 1
	TEMPORARY_ACCOUNT = 2
	GENERAL_ACCOUNT = 3
)


//注册成功后会向普通账户数据库中插入一个账户
func insertAccount(dbname string, account int, password string) error {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO accounts(account, password, registertime, isavailable) values(?,?,?,?)")
	if err != nil {
		log.Println(err)
		return err
	}
	time := tf.FormatTime()
	_, err = stmt.Exec(account, password, time, 1)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}



//查看数据库中所有存在的账户(包括可用和不可用的账户， 不包括删除过后的账户)
func ReadAccounts(accountType int) {
	var dbname string
	switch accountType {
	case SPECIAL_ACCOUNT:
		dbname = "database/SpecialAccount.db"
	case TEMPORARY_ACCOUNT:
		dbname = "database/TemporaryAccount.db"
	case GENERAL_ACCOUNT:
		dbname = "database/GeneralAccount.db"
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println("什么结果都没有查到...", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM accounts")
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	var (
		aid int
		account int
		password string
		registertime string
		isavailable bool
	)
	for rows.Next() {
		err = rows.Scan(&aid, &account, &password, &registertime, &isavailable)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(aid, account, password, registertime, isavailable)
	}
}



//倒序查看所有账户，不包括已经删除的账户
func ReadReversedAccounts(accountType int) {
	var dbname string
	switch accountType {
	case SPECIAL_ACCOUNT:
		dbname = "database/SpecialAccount.db"
	case TEMPORARY_ACCOUNT:
		dbname = "database/TemporaryAccount.db"
	case GENERAL_ACCOUNT:
		dbname = "database/GeneralAccount.db"
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println("什么结果都没有查到...", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM accounts order by aid desc")
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	var (
		aid int
		account int
		password string
		registertime string
		isavailable bool
	)
	for rows.Next() {
		err = rows.Scan(&aid, &account, &password, &registertime, &isavailable)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(aid, account, password, registertime, isavailable)
	}
}


//更新账户的状态：正常使用和停用。1:正常 2：停用
func UpdateAccount(account, isavailable int) {
	if account <= 100100 {
		return
	}
	var dbname string
	if account <= 101101 {
		dbname = "database/TemporaryAccount.db"
	} else {
		dbname = "database/GeneralAccount.db"
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	//更新数据
	stmt, err := db.Prepare("update accounts set isavailable=? where account=?")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(isavailable, account)
	if err != nil {
		log.Println(err)
	}
}


//删除一个账户。
func deleteAccount(account int) {
	if account <= 101101 {
		return
	}
	db, err := sql.Open("sqlite3", "database/GeneralAccount.db")
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	//删除数据
	stmt, err := db.Prepare("delete from accounts where account=?")
	if err != nil {
		log.Println(err)
		return
	}
	// res, err := stmt.Exec(3)
	_, err = stmt.Exec(account)
	if err != nil {
		log.Println(err)
		return
	}
	// affect, err := res.RowsAffected()
	// checkErr(err)
	// fmt.Println(affect)
}


//账户是否不在于数据库中，如果不存在，则登录失败，或者可以注册，等等。用于后续
func isAccountNotExist(accountType int, account int) bool {
	var dbname string
	switch accountType {
	case SPECIAL_ACCOUNT:
		dbname = "database/SpecialAccount.db"
	case TEMPORARY_ACCOUNT:
		dbname = "database/TemporaryAccount.db"
	case GENERAL_ACCOUNT:
		dbname = "database/GeneralAccount.db"
	default:
		return false
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return false
	}
	var acc int
	err = db.QueryRow("SELECT account FROM accounts WHERE account=?", account).Scan(&acc)
	if err == sql.ErrNoRows {
		return true
	}
	return false
}


//无差别注册一个账户，包括管理员账户，临时账户，普通账户。不开放
func register(accountType int, account int, password string) error {
	if account <= 101101 {
		return errors.New("账户已经存在...")
	}
	var dbname string
	switch accountType {
	case SPECIAL_ACCOUNT:
		dbname = "database/SpecialAccount.db"
	case TEMPORARY_ACCOUNT:
		dbname = "database/TemporaryAccount.db"
	case GENERAL_ACCOUNT:
		dbname = "database/GeneralAccount.db"
	default:
		dbname = "database/GeneralAccount.db"
	}
	if isAccountNotExist(accountType, account) {
		return insertAccount(dbname, account, password)	
	}
	return errors.New("账户已经存在...")
}


//密码为加密过后的密码
//注册一个普通账户
func RegisterNewAccount(account int, password string) error {
	return register(GENERAL_ACCOUNT, account, password)
}


//生成一个可用的、未被注册的账户
func NewAccountNumber() (int, error) {
	db, err := sql.Open("sqlite3", "database/GeneralAccount.db")
	// db, err := sql.Open("sqlite3", "database/TemporaryAccount.db")
	if err != nil {
		log.Println(err)
		return 0, errors.New("未知错误...")
	}
	var acc int
	// err = db.QueryRow("SELECT account FROM accounts WHERE account=?", account).Scan(&acc)
	err = db.QueryRow("SELECT account FROM accounts order by aid desc").Scan(&acc)
	// e, err := db.Exec("SELECT MAX(ACCOUNT) FROM accounts").Scan(&acc)
	// fmt.Println(e, err)
	if err == sql.ErrNoRows {
		return 101102, nil
	}
	return acc + 1, nil
}