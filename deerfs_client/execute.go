package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/cheggaaa/pb/v3"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/owlcache/tools"
)

//剪切文件
func cut_file(chunkSize int64) {

	fmt.Println("Start cutting file......")

	//读取文件的基本信息
	fileInfo, err := os.Stat(*filePath)
	if err != nil {
		panic(err)
	}

	//计算文件MD5值，然后开展工作
	file_md5, md5_err := file_manage.GetFileMD5ByPath(*filePath)
	if md5_err != nil {
		panic(md5_err.Error())
	}

	//判断文件类型,返回文件类型字符串和MIME信息字符串
	file_type, _ := file_manage.CheckFileType(*filePath)
	//获取文件的后缀名
	fileExt = path.Ext(fileInfo.Name())

	//生成文件sign
	//sign格式为:文件MD5+加密(文件字节数+"-"+文件类型+"-"+四位随机数)
	str_filesize := strconv.FormatInt(fileInfo.Size(), 10) //int64->string
	//此处的加密只是为了url的美观性和确保文件的相对唯一性，并不是为了防止破译，
	//所以加密的key采用了包内内置，相对固定。deerfs后端接收后也用了相同的的key来解密。
	enCrypt_str := file_manage.EnCryptToString(common.JoinString(str_filesize, "-", file_type, "-", common.RandAllString(4)))
	fileSign = common.JoinString(file_md5, enCrypt_str) //计算sign

	//创建临时存储文件夹
	saveFolder = common.JoinString("./temp/", fileSign, "/")
	os.RemoveAll(saveFolder) //先删除一下本地可能存在重复的临时文件夹
	createFolder_err := file_manage.CreateFolder(saveFolder)
	if createFolder_err != nil {
		panic(createFolder_err.Error())
	}

	fileChunkNumber = math.Ceil(float64(fileInfo.Size()) / float64(chunkSize))

	cut_bar := pb.StartNew(int(fileChunkNumber)) // create and start new bar

	fi, err := os.OpenFile(*filePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return
	}
	b := make([]byte, chunkSize)
	var i int64 = 1
	for ; i <= int64(fileChunkNumber); i++ {

		cut_bar.Increment() //bar

		fi.Seek((i-1)*chunkSize, 0)
		if len(b) > int(fileInfo.Size()-(i-1)*chunkSize) {
			b = make([]byte, fileInfo.Size()-(i-1)*chunkSize)
		}
		fi.Read(b)

		f, err := os.OpenFile(saveFolder+strconv.Itoa(int(i))+".deerfs", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			deerfs_client_panic(err)
		}
		f.Write(b)
		f.Close()
	}
	fi.Close()

	// finish bar
	cut_bar.Finish()

}

//合并文件
func merge_file() {

	fmt.Println("Start merging files......")

	fii, err := os.OpenFile(common.JoinString(fileSign, fileExt), os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		deerfs_client_panic(err)
		return
	}

	merge_bar := pb.StartNew(int(fileChunkNumber)) // create and start new bar

	for i := 1; i <= int(fileChunkNumber); i++ {

		merge_bar.Increment() //bar

		f, err := os.OpenFile(saveFolder+strconv.Itoa(int(i))+".deerfs", os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
			return
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
			return
		}
		fii.Write(b)
		f.Close()
	}
	fii.Close()

	// finish bar
	merge_bar.Finish()

}

//执行发送整体的单个完整文件
func PostIntactFile() {

	ohc := NewOHC()
	ohc.Timeout(time.Duration(*http_request_timeout))

	fmt.Println(tools.JoinString("Send ", *filePath, " data to the url ", *deerfsAddress, " ......"))
	resp, errs := ohc.postIntactFile()
	if errs != nil || resp.StatusCode >= 300 {
		errstr := tools.ErrorSliceJoinToString(errs)
		if errstr != "" {
			fmt.Println(tools.JoinString("PostIntactFile error:", errstr))
		}
	}

	defer resp.Body.Close() //资源释放

	body, ioerr := ioutil.ReadAll(resp.Body)
	if ioerr != nil {
		fmt.Println(tools.JoinString("PostIntactFile error -> ioutil.ReadAll error:", ioerr.Error()))
	}

	ohc.ClearSuperAgent() //资源释放

	owlres := JsonByteToOwlResponse(body)
	if owlres.Status != 200 {
		fmt.Println("      ")
		fmt.Println("PostIntactFile Error:")
		fmt.Println("Status:", owlres.Status)
		fmt.Println("Results:", owlres.Results)
		deerfs_client_exit()
	}
	//成功
	fmt.Println("Upload success!")
	fmt.Println(tools.JoinString("Resource address:", owlres.ResponseHost, owlres.Key))

}

//执行发送切割后的块文件
func PostChunksFile() {

	//=====删除服务端的块文件的临时文件=======
	ohc_hrt := NewOHC()
	ohc_hrt.Timeout(time.Duration(*http_request_timeout))
	ohc_hrt.postClearChunksFile()
	ohc_hrt.ClearSuperAgent() //资源释放

	//=========执行发送切割后的块文件==========
	fmt.Println("Sending chunks files to the server...")
	postchunks_bar := pb.StartNew(int(fileChunkNumber)) // create and start new bar

	for i := 1; i <= int(fileChunkNumber); i++ {

		postchunks_bar.Increment() //bar

		cpath := saveFolder + strconv.Itoa(int(i)) + ".deerfs"

		ohc := NewOHC()
		ohc.Timeout(time.Duration(*http_request_timeout))

		//fmt.Println(tools.JoinString("Send ", cpath, " data to the url ", *deerfsAddress, " ......"))
		ohc.postChunksFile(cpath)

		ohc.ClearSuperAgent() //资源释放

	}

	// finish bar
	postchunks_bar.Finish()

	//向服务器发送块文件合并文件请求
	ohc_mcf := NewOHC()
	ohc_mcf.Timeout(time.Duration(*http_request_timeout))
	ohc_mcf.postMergeChunksFile()
	ohc_mcf.ClearSuperAgent() //资源释放

}
