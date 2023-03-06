package job

import (
	"fmt"
	"os"
	"time"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/owlcache_model"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"github.com/xssed/owlcache/cache"
	"go.uber.org/zap"
)

func Init() {

	//执行步骤信息
	fmt.Println("deerfs system job initialization...")
	//定期清理DB中过期的数据
	ClearExpireCacheFileData()
	//定时自动输出Owl的内存信息
	MemoryInfoToLog()
	//定时Ping命令
	Ping()
}

//清理过期的缓存文件数据
func ClearExpireCacheFileData() {

	task_clear_expire_cache_file_data := config.FileStorageClearExpireCacheFileData()
	ticker := time.NewTicker(time.Minute * time.Duration(task_clear_expire_cache_file_data))
	go func() {
		for _ = range ticker.C {
			//清理数据库过期数据
			//遍历集合
			global.CacheFile_List.KvStoreItems.Range(func(k, v interface{}) bool {
				//检查是否过期
				if v.(*cache.KvStore).IsExpired() {
					//过期的缓存文件进行删除
					//先判断是不是正常文件，是就删掉他
					if file_manage.IsFile(v.(*cache.KvStore).Key) {
						remove_err := os.Remove(v.(*cache.KvStore).Key) //删除文件
						if remove_err == nil {
							if !global.CacheFile_List.Delete(v.(*cache.KvStore).Key) {
								loger.Lg.Info("ClearExpireCacheFile:", zap.String("info", v.(*cache.KvStore).Key)) //日志记录
							} else {
								loger.Lg.Info("ClearExpireCacheFileError:", zap.String("info", v.(*cache.KvStore).Key)) //日志记录
							}
						}
					}
				}
				return true
			})

		}
	}()

}

//统计内存使用情况
func MemoryInfoToLog() {

	task_memoryinfotolog := config.SystemTagTaskMemoryInfoToLog()
	ticker := time.NewTicker(time.Minute * time.Duration(task_memoryinfotolog))
	go func() {
		for _ = range ticker.C {
			fmt.Println("Memory_info:", global.MemStats())
			loger.Lg.Info("Memory_info:", zap.String("info", global.MemStats())) //日志记录
		}
	}()

}

//定时Ping命令
func Ping() {

	owlcache_ping_interval := config.OwlcachePingInterval()
	ticker := time.NewTicker(time.Second * time.Duration(owlcache_ping_interval))
	go func() {
		for _ = range ticker.C {
			//发送ping命令检验owlcache是否正常运行
			if owlcache_model.Ping() == "" {
				//ping命令异常
				global.Owlcache_State = 0
			} else {
				global.Owlcache_State = 1
			}

		}
	}()

}
