package dao
import (
	"can/msg"
	"log"
	"encoding/json"
	"fmt"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

/*
	1.用来保存消息到数据库。包括：未读消息和历史消息
*/


//保存单个消息，保存完毕即关闭数据库的连接。消息过多时，会严重拖慢速度不宜使用
func SaveMessage(dbname string, messageText []byte, received int) {
	v := &msg.NormalMessage{}
	err := json.Unmarshal(messageText, v)
	if err != nil {
		log.Println(err)
		return
	}
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()
	stmt, err := db.Prepare("INSERT INTO message(sender, receiver, content, date, time, received) values(?,?,?,?,?,?)")
	if err != nil {
		log.Println(err)
		return
	}
	_, err = stmt.Exec(v.Data.Sender, v.Data.Receiver, v.Data.Content, v.Data.Date, v.Data.Time, received)
	if err != nil {
		log.Println(err)
		return
	}
}


//同时保存大量消息到一个数据库
func SaveMessages(dbname string, messageData [][]byte, received int) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println(err)
		return
	}
	defer db.Close()

	v := &msg.NormalMessage{}
	for _, eachmsg := range messageData {
		err := json.Unmarshal(eachmsg, v)
		if err != nil {
			log.Println(err)
			continue
		}
		stmt, err := db.Prepare("INSERT INTO message(sender, receiver, content, date, time, received) values(?,?,?,?,?,?)")
		if err != nil {
			log.Println(err)
			continue
		}
		_, err = stmt.Exec(v.Data.Sender, v.Data.Receiver, v.Data.Content, v.Data.Date, v.Data.Time, received)
		if err != nil {
			log.Println(err)
			continue
		}
	}
}


//查看dbname数据库里 所有 的消息记录。
func ReadMessages(dbname string) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM message")
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	var (
		mid int
		sender int
		receiver int
		content string
		date string
		time string
		received bool
	)
	for rows.Next() {
		err = rows.Scan(&mid, &sender, &receiver, &content, &date, &time, &received)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(mid, sender, receiver, content, date, time, received)
	}
}



//读取某个账户发送的所有消息
func ReadMessagesbySender(dbname string, account int) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM message WHERE sender=?", account)
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	var (
		mid int
		sender int
		receiver int
		content string
		date string
		time string
		received bool
	)
	for rows.Next() {
		err = rows.Scan(&mid, &sender, &receiver, &content, &date, &time, &received)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(mid, sender, receiver, content, date, time, received)
	}
}


//读取该账户收到的所有消息
func ReadMessagesbyReceiver(dbname string, account int) {
	db, err := sql.Open("sqlite3", dbname)
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM message WHERE receiver=?", account)
	if err != nil {
		log.Println("出现错误...", err)
		return
	}
	var (
		mid int
		sender int
		receiver int
		content string
		date string
		time string
		received bool
	)
	for rows.Next() {
		err = rows.Scan(&mid, &sender, &receiver, &content, &date, &time, &received)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(mid, sender, receiver, content, date, time, received)
	}
}