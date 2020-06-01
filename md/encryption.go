package md
import (
	"crypto/md5"
	"fmt"
)

/*
	加密函数，用来加密：密码，消息内容等等
*/


//加密密码，生成长度为32的字符串。
func Encryption(originalText string) string {
	h := md5.New()
	h.Write([]byte(originalText))
	cipherText := h.Sum(nil)
	return fmt.Sprintf("%x", cipherText)
}


//用于消息加密。自定义。