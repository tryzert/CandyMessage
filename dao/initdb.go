package dao
import (
	"database/sql"
	"log"
	_ "github.com/mattn/go-sqlite3"
)


/*
	用来生成并初始化数据库
	目前包括5个数据库：
		3个账户数据库(账户共6位10进制数字：范围: 100000 - 999999)：
			1) 管理员账户数据库。已经注册完毕，且不对外开放。账户范围: 100000 - 100100
			2) 临时账户数据库。这部分账户已经初始化，对外开放，可以免费使用。主要用于测试。范围: 100101 - 101101
				密码都为: 123456
			3) 普通账户数据库。开放注册。范围：101102 - 999999
		2个消息数据库：
*/

//账户数据库的初始化和表格的生成
/*
	说明：isavailable 1:账户状态正常 0：停用，封禁等等
	账户数据库表格式为:
			aid	account password registertime   		isavailable
	比如：	0	100101	1a2b3c4d 2020-05-30 18:23:35	1	
*/
func createAccountTable(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	sql_table := `
        CREATE TABLE IF NOT EXISTS accounts(
            aid INTEGER PRIMARY KEY AUTOINCREMENT,
            account INTEGER NULL,
            password CHAR(32) NULL,
            registertime CHAR(19) NULL,
            isavailable INTEGER NULL
        );
    `
	_, err = db.Exec(sql_table)
	if err != nil {
		log.Println(err)
		return
	}
}


//消息数据库的初始化和表格生成
func createMessageTable(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	sql_table := `
        CREATE TABLE IF NOT EXISTS message(
            mid INTEGER PRIMARY KEY AUTOINCREMENT,
            sender INTEGER NULL,
            receiver INTEGER NULL,
            content TEXT NULL,
            date CHAR(10) NULL,
            time CHAR(8) NULL,
            received INTEGER NULL
        );
    `
	_, err = db.Exec(sql_table)
	if err != nil {
		log.Println(err)
		return
	}
}


//统一初始化
func InitDataBase() {
	createAccountTable("database/SpecialAccount.db")
	createAccountTable("database/TemporaryAccount.db")
	createAccountTable("database/GeneralAccount.db")
	createMessageTable("database/UnreadMessage.db")
	createMessageTable("database/HistoryMessage.db")
}