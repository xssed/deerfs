package main

import (
	"flag"
)

const VERSION string = "0.1"

var filePath = flag.String("file_path", "", "Input file path.")
var cutSize = flag.Int64("cut_size", 0, "Input cut size.Note: The cut is in KB for ease of use.")
var deerfsAddress = flag.String("deerfs_address", "", "Input deerfs address.")                 //标准字符格式:http://127.0.0.1:7727
var upload_form_field = flag.String("upload_form_field", "upload", "Input upload form field.") //默认upload
var http_request_timeout = flag.Uint64("http_request_timeout", 10000, "Input http request timeout(Millisecond).")
var fileSign string
var fileName string
var fileExt string
var saveFolder string
var fileChunkNumber float64

func main() {

	//绑定输入数据并验证
	intput_and_verify()
	//输出必要信息
	output_info()
	//如果cutSize没有设置，按照单体文件发送单体文件
	if *cutSize == 0 {
		PostIntactFile()
	} else {
		//如果cutSize有设置，就进行切割后再上传
		//剪切文件
		cut_file(*cutSize)
		//上传剪切后的文件块
		PostChunksFile()
		//合并本地文件
		//merge_file()
		//退出
		deerfs_client_exit()
	}

}
