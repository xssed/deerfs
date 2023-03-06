package cachefile_manage

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"github.com/xssed/owlcache/cache"
	"go.uber.org/zap"
)

//缓存文件管理启动
func Cachefile_manage_start() {

	fmt.Println("Start cachefile manage...")
	loger.Lg.Info("Start cachefile manage...")

	//初始化缓存文件表
	global.CacheFile_List = cache.NewCache("cache_file")

	//创建缓存文件根目录的路径字符串
	cacheMainFolder := path.Join(config.FileStorageDirectoryPath(), "cache")
	//fmt.Println(cacheMainFolder)

	//缓存目录初始化
	cacheMainFolder_start(cacheMainFolder)

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("new watcher error:", zap.String("error", err.Error()))
		loger.Lg.Fatal("new watcher error:", zap.String("error", err.Error()))
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// fmt.Println("Watcher event:", event)
				// loger.Lg.Info("Watcher event:", zap.String("Watcher event", event.String()))

				if event.Has(fsnotify.Create) {

					//此处要判断目标是文件夹还是文件
					if file_manage.IsDir(event.Name) {
						//是文件夹
						//fmt.Println("Watcher event:", event)
						loger.Lg.Info("Watcher event:", zap.String("Watcher event", event.String()))
						//监听子目录
						listen_dir(watcher, event.Name)
					} else {
						//是文件
						//fmt.Println("Watcher event:", event)

						temp_replace_space_str := strings.ReplaceAll(event.String(), " ", "")              //过滤空格
						temp_replace_slash_str := strings.ReplaceAll(temp_replace_space_str, "\\\\", "/")  //过滤斜杠
						temp_replace_semicolon_str := strings.ReplaceAll(temp_replace_slash_str, "\"", "") //过滤分号
						temp_path_str := ""
						if temp_replace_semicolon_str[0:6] == "CREATE" {
							temp_path_str = temp_replace_semicolon_str[6:]
						}

						if temp_path_str != "" {
							cachefile_exptime, _ := time.ParseDuration(common.JoinString(config.FileStorageCacheFileExpireTime(), "s"))
							set_res := global.CacheFile_List.Set(temp_path_str, []byte(event.String()), cachefile_exptime)
							if !set_res {
								fmt.Println("Failed to add temporary record of cache file!")
							}
						}

					}

				}

				// if event.Has(fsnotify.Write) {
				// 	fmt.Println("modified file:", zap.String("modified file", event.Name))
				// 	loger.Lg.Info("modified file:", zap.String("modified file", event.Name))
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("watcher error:", zap.String("watcher error", err.Error()))
				loger.Lg.Fatal("watcher error:", zap.String("watcher error", err.Error()))
			}
		}
	}()

	// Add a path.
	//监听根目录
	listen_dir(watcher, cacheMainFolder)

	// Block main goroutine forever.
	<-make(chan struct{})
}

//监听目录
func listen_dir(watcher *fsnotify.Watcher, path string) {

	err := watcher.Add(path)
	if err != nil {
		fmt.Println("Watcher add error:", zap.String("Watcher add error", err.Error()))
		loger.Lg.Fatal("Watcher add error:", zap.String("Watcher add error", err.Error()))
	} else {
		fmt.Println("Watcher add success:", path)
		loger.Lg.Info("Watcher add success:", zap.String("path", path))
	}

}

//缓存目录初始化
func cacheMainFolder_start(cacheMainFolder string) {

	//启动时删除缓存目录下所有文件
	removeFolder_err := os.RemoveAll(cacheMainFolder)
	if removeFolder_err != nil {
		fmt.Println("Failed to create cache directory:", removeFolder_err.Error())
		loger.Lg.Fatal("Failed to create cache directory:", zap.String("error", removeFolder_err.Error()))
	}
	//建立缓存目录
	createFolder_err := file_manage.CreateFolder(cacheMainFolder)
	if createFolder_err != nil {
		fmt.Println("Failed to create cache directory:", createFolder_err.Error())
		loger.Lg.Fatal("Failed to create cache directory:", zap.String("error", createFolder_err.Error()))
	}

}
