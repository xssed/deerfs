package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

//绑定输入数据并验证
func intput_and_verify() {

	// 把用户传递的命令行参数解析为对应变量的值
	flag.Parse()

	temp_cutSize := *cutSize
	*cutSize = temp_cutSize * 1024 //默认单位是KB

	//===========必填参数部分============
	//filePath
	if len(*filePath) < 1 {
		fmt.Println("The 'file_path' option is required.")
		os.Exit(0)
	}
	fileName = path.Base(*filePath)

	//deerfsAddress
	if len(*deerfsAddress) < 11 {
		fmt.Println("The 'deerfs_address' option is required.")
		os.Exit(0)
	}
	//对deerfsAddress进行数据补全或者删除
	if strings.HasPrefix(*deerfsAddress, "//") {
		temp_deerfsAddress := *deerfsAddress
		*deerfsAddress = temp_deerfsAddress[2:]
	}
	if !strings.HasPrefix(*deerfsAddress, "http://") && !strings.HasPrefix(*deerfsAddress, "https://") {
		*deerfsAddress = "http://" + *deerfsAddress
	}
	if strings.HasSuffix(*deerfsAddress, "/") {
		temp_deerfsAddress := *deerfsAddress
		*deerfsAddress = temp_deerfsAddress[:len(temp_deerfsAddress)-1]
	}
	if strings.HasSuffix(*deerfsAddress, "//") {
		temp_deerfsAddress := *deerfsAddress
		*deerfsAddress = temp_deerfsAddress[:len(temp_deerfsAddress)-2]
	}
	u, parse_err := url.Parse(*deerfsAddress)
	if parse_err != nil {
		fmt.Println("The value of 'deerfsAddress' is invalid.", parse_err.Error())
		os.Exit(0)
	}
	*deerfsAddress = u.String()
	//===========选填参数部分============
	//cutSize
	if *cutSize == 0 {
		fmt.Println("be careful! You did not input data for the 'cut_size' option, and the client will use the normal upload method.")
	}

}

//输出必要信息
func output_info() {
	// 输出命令行参数
	fmt.Println("Welcome to use simple deerfs client.")
	fmt.Println("Current deerfs client version:", VERSION)
	fmt.Println("File path:", *filePath)
	fmt.Println("File name:", fileName)
	fmt.Println("Cut file size:", *cutSize, " byte.")
	fmt.Println("The address of deerfs:", *deerfsAddress)
	fmt.Println("The upload form field:", *upload_form_field)
	fmt.Println("The http request timeout:", *http_request_timeout, " Millisecond.")
}

//友好化panic
func deerfs_client_panic(v interface{}) {
	//先删除本地临时文件夹
	os.RemoveAll(saveFolder)
	//报异常
	panic(v)
}

//友好化exit
func deerfs_client_exit() {
	//先删除本地临时文件夹
	os.RemoveAll(saveFolder)
	//退出
	os.Exit(0)
}
