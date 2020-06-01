package tf
import (
	"time"
	"strings"
)

/*
	用于时间的格式化
*/


//将日期和时间生成一个字符串
func FormatTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}


//将日期和时间拆分
func FormatTime2() []string {
	return strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")
}