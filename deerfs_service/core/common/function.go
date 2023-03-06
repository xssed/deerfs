package common

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"math"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
)

//生成32位md5字串
func GetMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

//判断字符串是否为数字
func IsNumber(s string) bool {
	if s == "" {
		return false
	}
	//去除首尾空格
	s = strings.TrimSpace(s)
	for i := 0; i < len(s); i++ {
		// 存在 e 或 E, 判断是否为科学计数法
		if s[i] == 'e' || s[i] == 'E' {
			return IsSciNum(s[:i], s[i+1:])
		}
	}
	// 否则判断是否为整数或小数
	return IsInt(s) || IsDec(s)
}

//是否为科学计数法
func IsSciNum(num1, num2 string) bool {
	//e 前后字符串长度为0 是错误的
	if len(num1) == 0 || len(num2) == 0 {
		return false
	}
	// e 后面必须是整数，前面可以是整数或小数  4  +
	return (IsInt(num1) || IsDec(num1)) && IsInt(num2)
}

//判断是否为小数
func IsDec(s string) bool {
	// eg: 11.15, -0.15, +10.15, 3., .15,
	// err: +. 0..
	match1, _ := regexp.MatchString(`^[\+-]?\d*\.\d+$`, s)
	match2, _ := regexp.MatchString(`^[\+-]?\d+\.\d*$`, s)
	return match1 || match2
}

//判断是否为整数
func IsInt(s string) bool {
	match, _ := regexp.MatchString(`^[\+-]?\d+$`, s)
	return match
}

//高效拼接字符串
func JoinString(args ...string) string {
	var args_buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		args_buffer.WriteString(args[i])
	}
	return args_buffer.String()
}

//Int64转字符串
func Int64ToString(args int64) string {
	return strconv.FormatInt(args, 10)
}

//生成UUID
func GetUUIDString() string {
	uuid := uuid.Must(uuid.NewV4())
	return uuid.String()
}

/*RandAllString  生成随机字符串([a~zA~Z0~9])
  lenNum 长度
*/
func RandAllString(lenNum int) string {

	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	str := strings.Builder{}
	length := len(chars)
	for i := 0; i < lenNum; i++ {
		l := chars[rand.Intn(length)]
		str.WriteString(l)
	}
	return str.String()
}

/*RandNumString  生成随机数字字符串([0~9])
  lenNum 长度
*/
func RandNumString(lenNum int) string {

	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	str := strings.Builder{}
	length := 10
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[52+rand.Intn(length)])
	}
	return str.String()
}

/*RandString  生成随机字符串(a~zA~Z])
  lenNum 长度
*/
func RandString(lenNum int) string {

	var chars = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}

	str := strings.Builder{}
	length := 52
	for i := 0; i < lenNum; i++ {
		str.WriteString(chars[rand.Intn(length)])
	}
	return str.String()
}

//将数据进行Base64编码
func Base64_Encode(src []byte) []byte {

	enc_str := base64.StdEncoding.EncodeToString(src)
	return []byte(enc_str)

}

//将数据进行Base64解码
func Base64_Decode(enc_str string) ([]byte, error) {

	return base64.StdEncoding.DecodeString(enc_str)

}

//对浮点数四舍五入-保留小数点后n位
func RoundedFixed(val float64, n int) float64 {
	change := math.Pow(10, float64(n))
	fv := 0.0000000001 + val //对浮点数产生.xxx999999999 计算不准进行处理
	return math.Floor(fv*change+.5) / change
}

//高效将错误类型切片拼接字符串
func ErrorSliceJoinToString(args []error) string {
	var args_buffer bytes.Buffer
	for i := 0; i < len(args); i++ {
		args_buffer.WriteString(args[i].Error())
	}
	return args_buffer.String()
}
