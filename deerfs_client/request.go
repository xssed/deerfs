package main

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/owlcache/tools"
)

//定义OwlcacheHttpclient
type OwlcacheHttpclient struct {
	sa *gorequest.SuperAgent
}

//初始化gorequest库
func NewOHC() *OwlcacheHttpclient {
	var grsa *gorequest.SuperAgent
	grsa = gorequest.New()
	return &OwlcacheHttpclient{sa: grsa}
}

//设置超时
func (ohc *OwlcacheHttpclient) Timeout(owlcache_http_request_timeout time.Duration) {
	ohc.sa.Timeout(owlcache_http_request_timeout * time.Millisecond)
}

//清理资源
func (ohc *OwlcacheHttpclient) ClearSuperAgent() {
	ohc.sa.ClearSuperAgent()
}

//发送整体的单个完整文件
func (ohc *OwlcacheHttpclient) postIntactFile() (*Response, []error) {

	ohc.sa.Post(common.JoinString(*deerfsAddress, "/deerfs_upload/upload")).Type("multipart").
		SendFile(*filePath, "", *upload_form_field)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()
	return &Response{r_res}, r_err_slices

}

//发送切割后的块文件
func (ohc *OwlcacheHttpclient) postChunksFile(cpath string) {

	post_send := make(map[string]string)
	post_send["file_sign"] = fileSign

	ohc.sa.Post(common.JoinString(*deerfsAddress, "/deerfs_upload/upload_chunks")).Type("multipart").
		SendFile(cpath, "", *upload_form_field).Send(post_send)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()

	if r_err_slices != nil || r_res.StatusCode >= 300 {
		errstr := tools.ErrorSliceJoinToString(r_err_slices)
		if errstr != "" {
			fmt.Println(tools.JoinString("postChunksFile error:", errstr))
		}
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		fmt.Println(tools.JoinString("postChunksFile error -> ioutil.ReadAll error:", ioerr.Error()))
	}
	owlres := JsonByteToOwlResponse(body)
	if owlres.Status != 200 {
		fmt.Println("      ")
		fmt.Println("PostChunksFile Error:")
		fmt.Println("Status:", owlres.Status)
		fmt.Println("Results:", owlres.Results)
		deerfs_client_exit()
	}

}

//删除服务端的块文件的临时文件
func (ohc *OwlcacheHttpclient) postClearChunksFile() {

	fmt.Println("Request temp storage from the server...")

	post_send := make(map[string]string)
	post_send["file_sign"] = fileSign

	ohc.sa.Post(common.JoinString(*deerfsAddress, "/deerfs_upload/clear_chunks")).Type("multipart").
		Send(post_send)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()

	if r_err_slices != nil || r_res.StatusCode >= 300 {
		errstr := tools.ErrorSliceJoinToString(r_err_slices)
		if errstr != "" {
			fmt.Println(tools.JoinString("postClearChunksFile error:", errstr))
		}
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		fmt.Println(tools.JoinString("postClearChunksFile error -> ioutil.ReadAll error:", ioerr.Error()))
	}
	owlres := JsonByteToOwlResponse(body)
	if owlres.Status != 200 {
		fmt.Println("      ")
		fmt.Println("postClearChunksFile Error:")
		fmt.Println("Status:", owlres.Status)
		fmt.Println("Results:", owlres.Results)
		deerfs_client_exit()
	}

}

//向服务器发送块文件合并文件请求
func (ohc *OwlcacheHttpclient) postMergeChunksFile() {

	fmt.Println("Send the merge file command to the server...")

	post_send := make(map[string]string)
	post_send["file_sign"] = fileSign
	post_send["file_name"] = fileName
	post_send["file_chunk_number"] = strconv.FormatFloat(fileChunkNumber, 'f', 0, 64)

	ohc.sa.Post(common.JoinString(*deerfsAddress, "/deerfs_upload/merge_chunks")).Type("multipart").
		Send(post_send)

	//发送请求获取数据
	r_res, _, r_err_slices := ohc.sa.EndBytes()

	if r_err_slices != nil || r_res.StatusCode >= 300 {
		errstr := tools.ErrorSliceJoinToString(r_err_slices)
		if errstr != "" {
			fmt.Println(tools.JoinString("postMergeChunksFile error:", errstr))
		}
	}

	defer r_res.Body.Close() //资源释放
	body, ioerr := ioutil.ReadAll(r_res.Body)
	if ioerr != nil {
		fmt.Println(tools.JoinString("postMergeChunksFile error -> ioutil.ReadAll error:", ioerr.Error()))
	}
	owlres := JsonByteToOwlResponse(body)
	if owlres.Status != 200 {
		fmt.Println("      ")
		fmt.Println("postMergeChunksFile Error:")
		fmt.Println("Status:", owlres.Status)
		fmt.Println("Results:", owlres.Results)
		deerfs_client_exit()
	}
	//成功
	fmt.Println("Upload success!")
	fmt.Println(tools.JoinString("Resource address:", owlres.ResponseHost, owlres.Key))

}
