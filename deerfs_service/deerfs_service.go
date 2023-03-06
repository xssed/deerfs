package main

import (
	"runtime"

	"github.com/xssed/deerfs/deerfs_service/core/application"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/job"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
)

func main() {

	//使用多核cpu(Use multi-core cpu)
	runtime.GOMAXPROCS(runtime.NumCPU())
	//输出欢迎信息
	global.DosSayHello()
	//配置初始化
	config.Init()
	//全局变量初始化
	global.Init()
	//日志初始化
	loger.Init()
	//计划任务初始化
	job.Init()
	//开启服务器
	application.StartServer()

}
