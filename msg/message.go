package msg
import (
	"net"
	"encoding/json"
	"errors"
	"can/tf"
)

/*
	定义各种消息结构体。
	一个消息基本类。
	衍生出：
		0) 注册消息
		1) 登录消息		
		2) 心跳包
		3) 普通个人消息
		4) 群组消息
		5) 临时消息(不保存/或者短期保存，阅后即焚，过期即删)
		6) 在线消息(必须双方都在线，一方发送，另一方才能收到)
		7) 离线消息(对方如果离线，消息会保存)
*/

const (
	REGISTER int = 0
	LOGIN int = 1
	HEART_BEAT int = 2
	NORMAL int = 3
	GROUP int = 4
	TEMPORARY int = 5
)


//消息基本类
type BaseMessage struct {
	Type int `json:"type"`
	Data interface{} `json:"data"`
}



//注册消息中的数据格式
type RegisterMessageData struct {
	Account int `json:"account"`
	Password string `json:"password"`
	Key string `json:"key"`
}


//注册消息
type RegisterMessage struct {
	Type int `json:"type"`
	Data RegisterMessageData `json:"data"`
}


//登录消息中的数据格式
type LoginMessageData struct {
	Account int `json:"account"`
	Password string `json:"password"`
}


//登录消息
type LoginMessage struct {
	Type int `json:"type"`
	Data LoginMessageData `json:"data"`
}


//心跳包, Data部分放账号。当然必须先登录，否则心跳包不起作用
type HeartBeatMessage struct {
	Type int `json:"type"`
	Data int `json:"data"`
}


//点对点消息中的数据格式
type NormalMessageData struct {
	Sender int `json:"sender"`
	Receiver int `json:"receiver"`
	Content string `json:"content"`
	Date string `json:"date"` //2020.12.12
	Time string `json:"time"` //20:30:05
	Received bool `json:"received"`// 0:false 1:yes
}


//点对点消息包
type NormalMessage struct {
	Type int `json:"type"`
	Data NormalMessageData `json:"data"`
}



//群组消息中的数据格式
type GroupMessageData struct {
	Sender int `json:"sender"`
	Receiver []int `json:"receiver"`
	Content string `json:"content"`
	Date string `json:"date"` //2020.12.12
	Time string `json:"time"` //20:30:05
	Received []bool `json:"received"`// 0:false 1:yes
}


//群组消息包
type GroupMessage struct {
	Type int `json:"type"`
	Data GroupMessageData `json:"data"`
}






/* 
	-----------------------------------------------------------------------
*/

//生成一个注册账号消息
func NewRegisterMessage(account int, password string, key string) *RegisterMessage {
	data := RegisterMessageData{
		Account: account,
		Password: password,
		Key: key,
	}
	return &RegisterMessage{
		Type: REGISTER,
		Data: data,
	}
}


//生成一个登录消息
func NewLoginMessage(sender int, password string) *LoginMessage {
	data := LoginMessageData{
		Account: sender,
		Password: password,
	}
	return &LoginMessage{
		Type: LOGIN,
		Data: data,
	}
}


//生成一个普通消息：即点对点，单人消息，私聊消息
func NewNormalMessage(sender int, receiver int, content string) *NormalMessage {
	t := tf.FormatTime2()
	date := t[0]
	time := t[1]
	data := NormalMessageData{
		Sender:sender,
		Receiver:receiver,
		Content:content,
		Date:date,
		Time:time,
		Received:false,
	}
	return &NormalMessage{
		Type: NORMAL,
		Data: data,
	}
}


//生成一个群组消息
func NewGroupMessage() *GroupMessage {
	data := GroupMessageData{}
	return &GroupMessage{
		Type: GROUP,
		Data: data,
	}
}


//生成一个心跳包: data应该为account
func NewHeartBeatMessage(data int) *HeartBeatMessage{
	return &HeartBeatMessage{
		Type: HEART_BEAT,
		Data: data,
	}
}



//服务器的连接地址
//const address = "49.235.179.226:10000"
const address = "127.0.0.1:10000"
//将消息发送给服务器或客户端
func SendMessage(msg interface{}, conn net.Conn, close bool) error {
	if close {
		defer conn.Close()
	}
	transport_data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	transport_data = append(transport_data, '\t')	
	_, err = conn.Write(transport_data)
	if err != nil {
		return errors.New("message发送失败...")
	}
	return nil
}


func SendData(data []byte, conn net.Conn, close bool) error {
	if close {
		defer conn.Close()
	}
	_, err := conn.Write(data)
	if err != nil {
		return errors.New("data发送失败...")
	}
	return nil
}