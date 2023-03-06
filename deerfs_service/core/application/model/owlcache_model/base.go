package owlcache_model

import (
	"fmt"
	"os"

	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
)

//初始化连接owlcache
func Init() {

	// 记录日志
	loger.Lg.Info("Database owlcache connection initialization......")
	fmt.Println("Database owlcache connection initialization......")

	//发送ping命令检验owlcache是否正常运行
	if Ping() == "" {
		//fmt.Println("Ping owlcache error!")
		os.Exit(0)
	}
	//如果正常运行获取一个token,失败就退出
	if GetToken() == "" {
		fmt.Println("Failed to get token. Please check <owlcache_password> in the config file.Or owlcache http service exception.")
		os.Exit(0)
	}

}
