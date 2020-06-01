package service

import (
	"net"
	"time"
	"errors"
	"encoding/json"
	"can/msg"
	"can/dao"
)



/*
	REGISTER int 8 = 0
	LOGIN int8 = 1
	HEART_BEAT int8 = 2
	NORMAL int8 = 3
	GROUP int8 = 4
	TEMPORARY int8 = 5
*/

func RegisterMessageParser(data []byte, key string) error {
	v := &msg.RegisterMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return errors.New("收到的消息格式错误...")
	}
	vd := v.Data
	if vd.Key != key {
		return errors.New("注册令牌错误")
	}
	return dao.RegisterNewAccount(vd.Account, vd.Password)
}


//登录
func LoginMessageParser(data []byte) (int, error) {
	v := &msg.LoginMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return 0, errors.New("收到的消息格式错误...")
	}
	vd := v.Data
	if !dao.VerifyAccountValidity(vd.Account, vd.Password) {
		return 0, errors.New("账号或密码错误")
	}
	return vd.Account, nil
}



//回复心跳包，一定时间内回复3次，一次成功视为成功，否则回复失败。后续把client剔除
func HeartBeatMessageParser(data []byte, conn net.Conn) (int, error) {
	v := &msg.HeartBeatMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return 0, errors.New("收到的消息格式错误...")
	}
	ch := make(chan error, 1)
	go func() {
		for i := 0; i < 3; i++ {
			_, err := conn.Write(append(data, '\n'))
			if err == nil {
				ch <- nil
				break
			}
		}
	}()
	select {
	case res := <- ch :
		return v.Data, res
	case <- time.After(time.Second * 5):
		return 0, errors.New("回复心跳包超时...")
	}
}


//一般普通消息
func NormalMessageParser(data []byte) (int, error) {
	v := &msg.NormalMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return 0, errors.New("收到的消息格式错误...")
	}
	/*
	todo...
	*/
	return v.Data.Receiver, nil
}


//群组消息
func GroupMessageParser(data []byte) ([]int, error) {
	v := &msg.GroupMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return []int{}, errors.New("收到的消息格式错误...")
	}
	return v.Data.Receiver, nil
}


//临时消息
func TemporaryMessageParser(data []byte) (int, error) {
	v := &msg.NormalMessage{}
	err := json.Unmarshal(data, v)
	if err != nil {
		return 0, errors.New("收到的消息格式错误...")
	}
	return v.Data.Receiver, nil
}