package file_manage

import (
	//"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"

	//"path"
	"path/filepath"
	"time"

	encrypt "github.com/oliverCJ/crypt"

	"github.com/h2non/filetype"
	"github.com/xssed/deerfs/deerfs_service/core/common"

	//"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
)

//判断文件类型,返回文件类型字符串和MIME信息字符串
func CheckFileType(path string) (file_type string, mime string) {
	// Open a file descriptor
	file, _ := os.Open(path)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)

	defer file.Close()

	file_match, _ := filetype.Match(head)
	if file_match == filetype.Unknown {
		return "other", ""
	}

	return file_match.Extension, file_match.MIME.Value
}

//判断文件类型,返回文件类型字符串和MIME信息字符串
func CheckFileTypeByByte(b []byte) (file_type string, mime string) {

	file_match, _ := filetype.Match(b)
	if file_match == filetype.Unknown {
		return "other", ""
	}

	return file_match.Extension, file_match.MIME.Value
}

//判断文件类型,返回文件类型字符串和MIME信息字符串
func CheckFileTypeByMultipart(file *multipart.FileHeader) (file_type string, mime string) {
	src, err := file.Open()
	if err != nil {
		return "other", ""
	}
	defer src.Close()

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	src.Read(head)

	file_match, _ := filetype.Match(head)
	if file_match == filetype.Unknown {
		return "other", ""
	}

	return file_match.Extension, file_match.MIME.Value
}

//判断文件是不是图片，返回布尔值
func IsImage(path string) bool {
	// Open a file descriptor
	file, _ := os.Open(path)

	// We only have to pass the file header = first 261 bytes
	head := make([]byte, 261)
	file.Read(head)

	defer file.Close()

	if filetype.IsImage(head) {
		return true
	} else {
		return false
	}
}

//判断文件或路径是否存在
func PathExists(path string) (bool, error) {

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err

}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

//获取文件的大小,单位字节
func FileSize(path string) int64 {
	file, err := os.Stat(path)
	if err != nil {
		return 0
	}
	return file.Size()
}

//遍历文件夹统计文件夹使用量
func ListDirSize(folder string) {

	files, errDir := ioutil.ReadDir(folder)
	if errDir != nil {
		fmt.Println(errDir.Error())
		loger.Lg.Fatal(errDir.Error()) // 记录日志
	}

	for _, file := range files {
		if file.IsDir() {
			ListDirSize(folder + "/" + file.Name())
		} else {
			// 输出绝对路径
			strAbsPath, errPath := filepath.Abs(folder + "/" + file.Name())
			if errPath != nil {
				fmt.Println(errPath)
			}
			//对已经存储文件的字节量的统计
			global.Directory_Storage_Use_Size = global.Directory_Storage_Use_Size + FileSize(strAbsPath)
			//对已经存储文件的个量的统计
			global.Directory_Storage_File_Number = global.Directory_Storage_File_Number + 1

		}
	}

}

//创建目录
func CreateFolder(folder string) error {

	//目录校验
	temp_last := folder[len(folder)-1:]
	if temp_last != "/" {
		folder = folder + "/"
	}
	//创建目录
	if folder != "" {
		return os.MkdirAll(folder, os.ModePerm)
	}
	return errors.New("Failed to create folder.")

}

//字节的单位转换,保留两位小数,将字节数转换为带单位的字符串
func FormatFileSizeToString(fileSize int64) (size string) {
	if fileSize < 1024 {
		//return strconv.FormatInt(fileSize, 10) + "B"
		return fmt.Sprintf("%.2fB", float64(fileSize)/float64(1))
	} else if fileSize < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(1024))
	} else if fileSize < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(1024*1024*1024))
	} else if fileSize < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(fileSize)/float64(1024*1024*1024*1024))
	} else {
		//if fileSize < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	}
}

//字节的单位转换,保留两位小数,将字节数转换为指定的单位(指定单位必须大写字母)
//B,KB,MB,GB,TB,EB
func FormatFileSize(fileSize int64, unit string) string {

	if unit == "B" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1))
	} else if unit == "KB" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024))
	} else if unit == "MB" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024))
	} else if unit == "GB" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024))
	} else if unit == "TB" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024))
	} else if unit == "EB" {
		return fmt.Sprintf("%.2f", float64(fileSize)/float64(1024*1024*1024*1024*1024))
	} else {
		return "0.00"
	}

}

//生成对应目录的文件存储路径字符串
func GetFilePathString() string {

	currentTime := time.Now()
	yyyymm := currentTime.Format("200601") //年月
	dd := currentTime.Format("02")         //日

	return common.JoinString(yyyymm, "/", dd)

}

//生成文件对应的缓存目录相关字符串,专属函数
func GetFileCachePathString(input_str string) (uri string, uri_md5 string, head_uri string) {

	r_uri := []rune(input_str)
	uri = string(r_uri[1:]) //去掉"/"
	uri_md5 = common.GetMd5String(uri)
	head_uri = string(r_uri[1:3])

	return uri, uri_md5, head_uri

}

//加密数据
//使用默认码表，6位加密
func EnCryptToString(src string) string {

	enc := encrypt.DefaultEncrypt6.EnCryptToString([]byte(src))
	return enc

}

//解密数据
//使用默认码表，6位加密
func DeCryptString(enc string) ([]byte, error) {

	dec, err := encrypt.DefaultEncrypt6.DeCryptString(enc)
	if err != nil {
		fmt.Println("decrypt error:", err.Error())
		return []byte(""), err
	}
	return dec, nil

}

//根据文件路径来获取文件的MD5值
func GetFileMD5ByPath(filepath string) (string, error) {

	file, err := os.Open(filepath)
	if err != nil {
		return_str1 := "Open err"
		return return_str1, err
	}
	defer file.Close()

	md5hash := md5.New()
	if _, md5_err := io.Copy(md5hash, file); md5_err != nil {
		return_str2 := "Copy err"
		return return_str2, md5_err
	}

	md5_str := hex.EncodeToString(md5hash.Sum(nil))

	return md5_str, nil
}

//根据文件资源来获取文件的MD5值
func GetFileMD5(file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return_str1 := "Open err"
		return return_str1, err
	}
	defer src.Close()

	md5hash := md5.New()
	if _, md5_err := io.Copy(md5hash, src); md5_err != nil {
		return_str2 := "Copy err"
		return return_str2, md5_err
	}

	md5_str := hex.EncodeToString(md5hash.Sum(nil))
	return md5_str, nil
}

//遍历文件夹统计文件夹获取文件夹下的文件数量
func ListDirFileNumber(folder string) int {

	var return_number int = 1

	files, errDir := ioutil.ReadDir(folder)
	if errDir != nil {
		fmt.Println(errDir.Error())
		loger.Lg.Error(errDir.Error()) // 记录日志
		return 0
	}

	for _, file := range files {
		if file.IsDir() {
			ListDirFileNumber(folder + "/" + file.Name())
		} else {
			//对已经存储文件的个量的统计
			return_number = return_number + 1
		}
	}

	return return_number

}
