package config

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-ini/ini"
	"github.com/gookit/filter"
)

var conf *ini.File

//配置初始化
func Init() {

	//执行步骤信息
	fmt.Println("deerfs system config initialization...")

	var err error
	conf, err = ini.Load("deerfs_service.conf")
	if err != nil {
		log.Fatalf(err.Error())
	} else {
		fmt.Println("deerfs success loaded the configuration file.")
	}

}

//通用获取数据基本方法
func Get(str string) string {

	strArr := strings.Split(str, ".")

	if len(strArr) == 2 {
		return conf.Section(strArr[0]).Key(strArr[1]).String()
	}

	return conf.Section("").Key(strArr[0]).String()
}

//给当前主机的deerfs节点起一个在集群中的唯一名字，例如deerfs1，deerfs2，deerfs3。。。。。。
//deerfs启动后会将名字与mysql服务器进行数据更新，一旦名字注册后请谨慎修改名字。
func SystemTagName() string {

	name := Get("system_tag.name")
	return name

}

//该节点的资源定位地址
func SystemTagUriAddress() string {

	uri_address := Get("system_tag.uri_address")
	return uri_address

}

//是否允许该节点在集群中冗余存储相同文件。0关闭，1开启。默认为关闭。开启后，该节点存储集群中已经存在的文件时会重复存储(按需配置)。
func SystemTagAllowDuplicates() int {

	allow_duplicates, _ := filter.Int(Get("system_tag.allow_duplicates"))
	return allow_duplicates

}

//定时任务-自动输出内存信息。单位分钟。默认为5。
func SystemTagTaskMemoryInfoToLog() int {

	task_memory_info_to_log, _ := filter.Int(Get("system_tag.task_memory_info_to_log"))
	return task_memory_info_to_log

}

//deerfs的IP地址。
func HttpAddr() string {

	http_addr := Get("http.http_addr")
	if http_addr == "" {
		return "0.0.0.0"
	}

	return http_addr
}

//deerfs的端口。
func HttpPort() string {

	http_port := Get("http.http_port")
	if http_port == "" {
		return "7727"
	}

	return http_port
}

func HttpMode() string {

	http_mode := Get("http.http_mode")
	if http_mode == "" {
		return "debug"
	}

	return http_mode
}

//HTTP读超时,单位秒.
func HttpReadTimeout() int {

	http_read_timeout, _ := filter.Int(Get("http.http_read_timeout"))
	if http_read_timeout == 0 {
		return 60
	}

	return http_read_timeout
}

//HTTP写超时,单位秒。
func HttpWriteTimeout() int {

	http_write_timeout, _ := filter.Int(Get("http.http_write_timeout"))
	if http_write_timeout == 0 {
		return 60
	}

	return http_write_timeout
}

//owlcache(HTTP服务)的地址。
func OwlcacheAddr() string {

	owlcache_addr := Get("owlcache.owlcache_addr")
	if owlcache_addr == "" {
		return "http://127.0.0.1:7721"
	}

	return owlcache_addr
}

//owlcache(HTTP服务)的密码。
func OwlcachePassword() string {

	owlcache_password := Get("owlcache.owlcache_password")

	return owlcache_password
}

//间隔多久向owlcache发送一次ping命令。单位是秒。默认一秒。
func OwlcachePingInterval() int {

	owlcache_ping_interval, _ := filter.Int(Get("owlcache.owlcache_ping_interval"))
	if owlcache_ping_interval == 0 {
		return 1
	}

	return owlcache_ping_interval
}

//向owlcache(HTTP服务)请求超时的时间。
func OwlcacheHttpRequestTimeout() time.Duration {

	owlcache_http_request_timeout, _ := filter.Int64(Get("http.owlcache_http_request_timeout"))
	if owlcache_http_request_timeout == 0 {
		return time.Duration(4000)
	}

	return time.Duration(owlcache_http_request_timeout)
}

//向owlcache(HTTP服务)中存储文件信息的Key前缀字符串。默认是“deerfs::”。
func OwlcacheKeyStoragePrefix() string {

	owlcache_key_storage_prefix := Get("http.owlcache_key_storage_prefix")
	if owlcache_key_storage_prefix == "" {
		return "deerfs::"
	}
	return owlcache_key_storage_prefix
}

//向owlcache(HTTP服务)中存储文件信息的数据过期时间。单位秒。默认值为0永不过期。
func OwlcacheKeyStorageExpire() string {

	owlcache_key_storage_expire := Get("http.owlcache_key_storage_expire")
	if owlcache_key_storage_expire == "" {
		return "0"
	}
	return owlcache_key_storage_expire
}

//文件存储的目录。
func FileStorageDirectoryPath() string {

	directory_path := Get("file_storage.directory_path")
	if directory_path == "" {
		return "./deerfs_data/"
	}

	return directory_path
}

//节点目录最大存储容量。默认大小1TB。单位是字节。
func FileStorageDirectoryStorageMaxSize() int64 {

	directory_storage_max_size, _ := filter.Int64(Get("file_storage.directory_storage_max_size"))
	if directory_storage_max_size == 0 {
		return 1099511627776
	}

	return directory_storage_max_size
}

//单个文件的最大存储大小。默认大小5M。单位是字节。
func FileStorageFileStorageMaxSize() int64 {

	file_storage_max_size, _ := filter.Int64(Get("file_storage.file_storage_max_size"))
	if file_storage_max_size == 0 {
		return 5242880
	}

	return file_storage_max_size
}

//缓存文件的过期时间。默认为3600(一小时)。单位是秒。
func FileStorageCacheFileExpireTime() string {

	cache_file_expire_time := Get("file_storage.cache_file_expire_time")
	if cache_file_expire_time == "" {
		return "3600"
	}

	return cache_file_expire_time
}

//定时任务-自动清理过期的缓存文件。单位分钟。默认为1。
func FileStorageClearExpireCacheFileData() int {

	task_clear_expire_cache_file_data, _ := filter.Int(Get("file_storage.task_clear_expire_cache_file_data"))
	return task_clear_expire_cache_file_data

}

//表单提交字段。
func UploadFormField() string {

	form_field := Get("upload.form_field")
	if form_field == "" {
		return "upload"
	}

	return form_field
}

//上传区块的字段名。
func UploadFormChunksField() string {

	form_chunks_field := Get("upload.form_chunks_field")
	if form_chunks_field == "" {
		return "upload"
	}

	return form_chunks_field
}

//level options:debug,info,warn,error,dpanic,panic,fatal
func LogLevel() string {

	log_level := Get("log.log_level")
	if log_level == "" {
		return "debug"
	}

	return log_level
}

//日志文件的位置
func LogFilename() string {

	log_filename := Get("log.log_filename")
	if log_filename == "" {
		return "./deerfs_log/deerfs.log"
	}

	return log_filename
}

//在进行切割之前，日志文件的最大大小（以MB为单位）
func LogMaxsize() int {

	log_maxsize, _ := filter.Int(Get("log.log_maxsize"))
	if log_maxsize == 0 {
		return 200
	}

	return log_maxsize
}

//保留旧文件的最大天数
func LogMaxage() int {

	log_maxage, _ := filter.Int(Get("log.log_maxage"))
	if log_maxage == 0 {
		return 7
	}

	return log_maxage
}

//保留旧文件的最大个数
func LogMaxbackups() int {

	log_maxbackups, _ := filter.Int(Get("log.log_maxbackups"))
	if log_maxbackups == 0 {
		return 10
	}

	return log_maxbackups
}

//mysql的地址。
func MysqlHost() string {

	host := Get("mysql.host")
	if host == "" {
		return "127.0.0.1"
	}

	return host
}

//mysql的端口。
func MysqlPort() string {

	port := Get("mysql.port")
	if port == "" {
		return "3306"
	}

	return port
}

//连接mysql的账号。
func MysqlUser() string {

	user := Get("mysql.user")
	return user
}

//连接mysql的密码。
func MysqlPassword() string {

	password := Get("mysql.password")
	return password
}

//mysql中deerfs的数据库名称。
func MysqlDatabase() string {

	database := Get("mysql.database")
	if database == "" {
		return "deerfs"
	}

	return database
}

//mysql中deerfs的数据库对应编码。
func MysqlCharset() string {

	charset := Get("mysql.charset")
	if charset == "" {
		return "utf8"
	}

	return charset
}

//查询结果是否自动解析为时间。
func MysqlParsetime() string {

	parsetime := Get("mysql.parsetime")
	if parsetime == "" {
		return "True"
	}

	return parsetime
}

//MySQL的时区设置。
func MysqlLoc() string {

	loc := Get("mysql.loc")
	if loc == "" {
		return "Local"
	}

	return loc
}

//最大链接数。
func MysqlMaxidleconns() int {

	maxidleconns, _ := filter.Int(Get("log.maxidleconns"))
	if maxidleconns == 0 {
		return 10
	}

	return maxidleconns
}

//最大打开链接。
func MysqlMaxopenconns() int {

	maxopenconns, _ := filter.Int(Get("log.maxopenconns"))
	if maxopenconns == 0 {
		return 100
	}

	return maxopenconns
}

//设置表前缀。
func MysqlTableprefix() string {

	tableprefix := Get("mysql.tableprefix")
	if tableprefix == "" {
		return "deerfs_"
	}

	return tableprefix
}
