package global

import (
	"fmt"
	"runtime"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/owlcache/cache"
)

const (
	VERSION      string = "0.1-beta" //当前deerfs的版本
	VERSION_DATE string = "2023-03-03"
)

//当前deerfs主机的数据库中节点ID
var Node_ID int = 0

//创建一个全局的Node节点信息
var Node *mysql_model.Node

//当前deerfs主机的存储目录使用量,单位字节
var Directory_Storage_Use_Size int64 = 0

//当前deerfs主机的存储的文件数量,单位个
var Directory_Storage_File_Number int64 = 0

//创建一个全局的文件缓存表
var CacheFile_List *cache.BaseCache

//当前deerfs主机的缓存目录使用量,单位字节
var CacheFile_Storage_Use_Size int64 = 0

//设置一个全局的Token值，用来和owlcache进行数据交互
var Token string = ""

//设置一个全局的owlcache状态值，1在线,0未在线。
var Owlcache_State = 0

//设置一个全局的deerfs的Key前缀值，用来和owlcache进行数据交互
var Owlcache_deerfs_key_storage_prefix string = ""

//初始化
func Init() {

	//执行步骤信息
	fmt.Println("deerfs system global initialization...")

	//全局变量赋值
	Owlcache_deerfs_key_storage_prefix = config.OwlcacheKeyStoragePrefix()

}

//程序启动欢迎信息
func DosSayHello() {

	fmt.Println("Welcome to use deerfs. Version:" + VERSION + "\nIf you have any questions,Please contact us: xsser@xsser.cc \nProject Home:https://github.com/xssed/deerfs")
	fmt.Println(`        _                 __        `)
	fmt.Println(`     __| | ___  ___ _ __ / _|___    `)
	fmt.Println(`    / _| |/ _ \/ _ \ '__| |_/ __|   `)
	fmt.Println(`   | (_| |  __/  __/ |  |  _\__ \   `)
	fmt.Println(`    \__,_|\___|\___|_|  |_| |___/   `)
	fmt.Println(`                                    `)

}

//输出内存信息
func MemStats() string {

	var m runtime.MemStats

	// Sys 服务现在系统使用的内存
	// NumGC 垃圾回收调用次数
	// Alloc 堆空间分配的字节数
	// TotalAlloc 从服务开始运行至今分配器为分配的堆空间总和，只有增加，释放的时候不减少。

	runtime.ReadMemStats(&m)

	unit_gb := 1024 * 1024 * 1024.0
	alloc := float64(m.Alloc) / unit_gb
	totalalloc := float64(m.TotalAlloc) / unit_gb
	sys := float64(m.Sys) / unit_gb
	numgc := m.NumGC

	logstr := fmt.Sprintf("Sys = %vGB  NumGC = %v  Alloc = %vGB  TotalAlloc = %vGB  ", common.RoundedFixed(sys, 7), numgc, common.RoundedFixed(alloc, 7), common.RoundedFixed(totalalloc, 7))
	return logstr

}
