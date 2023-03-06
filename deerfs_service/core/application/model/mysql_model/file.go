package mysql_model

import (
	//"time"
	"errors"

	"github.com/jinzhu/gorm"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
)

//文件模型
type File struct {
	ID       int64  `gorm:"primary_key" json:"id"` // 主键，根据约定不需要
	NodeId   int    `gorm:"index" json:"node_id"`  // 索引
	Node     Node   `json:"node"`
	FileHash string `json:"file_hash"` //文件hash
	FileSign string `json:"file_sign"` //文件独有sign
	FileName string `json:"file_name"` //上传时的文件名
	FileExt  string `json:"file_ext"`  //上传时的文件扩展名
	FileType string `json:"file_type"` //文件类型(后台识别)
	FileMime string `json:"file_mime"` //文件MIME(后台识别)
	FileSize int64  `json:"file_size"` //文件大小
	FileAddr string `json:"file_addr"` //文件存储位置
	//CreateTime time.Time `json:"create_time"`
	Disabled int    `json:"disabled"` //状态(可用0/禁用1)，禁止非法文件访问
	Ext1     int    `json:"ext1"`     //备用字段1
	Ext2     string `json:"ext2"`     //备用字段2
	Model           // 自定义的 Model
}

//获取File的数据集合
func GetFiles(offset, limit int, cond map[string]interface{}) ([]*File, error) {
	var files []*File
	err := db.Where(cond).Offset(offset).Order("id desc").Limit(limit).Find(&files).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return files, nil
}

//获取一个File信息
func GetFileToHash(node_id int, file_hash string) (*File, error) {
	var file File
	err := db.Where("node_id = ? AND file_hash = ? AND deleted_on = ? ", node_id, file_hash, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &file, nil
}

//根据文件Sign获取一个File信息
func GetFileToSign(file_sign string) (*File, error) {
	var file File
	err := db.Where("file_sign = ? AND deleted_on = ? ", file_sign, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &file, nil
}

//根据ID获取一个File信息
func GetFileToId(id int) (*File, error) {
	var file File
	err := db.Where("id = ? AND deleted_on = ? ", id, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	return &file, nil
}

//判断文件是否存在(单节点)
func HasFileToHash(node_id int, file_hash, file_type string, file_size int64) (bool, error) {
	var file File
	err := db.Select("id").Where("node_id = ? AND file_hash = ? AND file_type = ? AND file_size = ?  AND deleted_on = ? ", node_id, file_hash, file_type, file_size, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//根据Sign判断文件是否存在
func HasFileToSign(file_sign string) (bool, error) {
	var file File
	err := db.Select("id").Where("file_sign = ?  AND deleted_on = ? ", file_sign, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//判断文件是否存在(所有节点)
func HasFileToHashByAll(file_hash, file_type string, file_size int64) (bool, error) {
	var file File
	err := db.Select("id").Where(" file_hash = ? AND file_type = ? AND file_size = ? AND deleted_on = ? ", file_hash, file_type, file_size, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//判断被软删除的文件是否存在
func HasDelFileToHash(node_id int, file_hash string) (bool, error) {
	var file File
	err := db.Select("id").Where("node_id = ? AND file_hash = ?  AND deleted_on != ? ", node_id, file_hash, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//判断被软删除的文件是否存在
func HasDelFileToSign(file_sign string) (bool, error) {
	var file File
	err := db.Select("id").Where(" file_hash = ?  AND deleted_on != ? ", file_sign, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//根据文件ID判断文件是否存在
func HasFileToId(id int64) (bool, error) {
	var file File
	err := db.Select("id").Where("id = ? AND deleted_on = ? ", id, 0).First(&file).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return false, err
	}
	// id 为正数才表示存在
	if file.ID > 0 {
		return true, nil
	}
	return false, nil
}

//统计文件数据条数
func GetFilesCount(cond map[string]interface{}) (int, error) {
	var count int
	if err := db.Model(&File{}).Where(cond).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

//添加一个文件
func AddFile(data map[string]interface{}) error {

	// 根据 data 参数构造 file 结构体
	file := File{
		// 关系：File 拥有 Node
		NodeId:   data["node_id"].(int),
		FileHash: data["file_hash"].(string),
		FileSign: data["file_sign"].(string),
		FileName: data["file_name"].(string),
		FileExt:  data["file_ext"].(string),
		FileType: data["file_type"].(string),
		FileMime: data["file_mime"].(string),
		FileSize: data["file_size"].(int64),
		FileAddr: data["file_addr"].(string),
		//CreateTime: data["create_time"].(time.Time),
		Disabled: data["disabled"].(int),
		Ext1:     data["ext1"].(int),
		Ext2:     data["ext2"].(string),
		Model: Model{
			CreatedBy: data["created_by"].(string),
		},
	}
	// 插入记录
	if err := db.Create(&file).Error; err != nil {
		return err
	}
	return nil
}

//修改节点数据,根据ID来修改
func EditFileForId(id int, data map[string]interface{}) error {
	if err := db.Model(&File{}).Where("id = ?  AND deleted_on = ?", id, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//修改一个文件,根据Hash修改
func EditFileForHash(node_id int, file_hash string, data map[string]interface{}) error {
	if err := db.Model(&File{}).Where("node_id = ? AND file_hash = ?  AND deleted_on = ?", node_id, file_hash, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//修改一个文件,根据Sign修改
func EditFileForSign(file_sign string, data map[string]interface{}) error {
	if err := db.Model(&File{}).Where(" file_sign = ?  AND deleted_on = ?", file_sign, 0).Updates(data).Error; err != nil {
		return err
	}
	return nil
}

//更改文件可访问状态，状态(可用0/禁用1)，禁止非法文件访问,通常用在色情暴力凶杀等非法文件一键封杀
func DisabledFile(node_id int, file_hash string) error {

	//获取一个File信息
	file, err := GetFileToHash(node_id, file_hash)
	if err != nil {
		return err
	}
	//判断文件存在
	if file.ID > 0 {

		if file.Disabled == 1 {
			if err := db.Model(&File{}).Where("node_id = ? AND file_hash = ?  AND deleted_on = ?", node_id, file_hash, 0).Update("disabled", 0).Error; err != nil {
				return err
			}
		} else {
			if err := db.Model(&File{}).Where("node_id = ? AND file_hash = ?  AND deleted_on = ?", node_id, file_hash, 0).Update("disabled", 1).Error; err != nil {
				return err
			}
		}
		return nil

	}
	return errors.New(errno.Msg[errno.FileIsNotExist])

}

//根据Sign更改文件可访问状态，状态(可用0/禁用1)，禁止非法文件访问,通常用在色情暴力凶杀等非法文件一键封杀
func DisabledFileBySign(node_id int, file_sign string) error {

	//获取一个File信息
	file, err := GetFileToSign(file_sign)
	if err != nil {
		return err
	}
	//判断文件存在
	if file.ID > 0 {

		if file.Disabled == 1 {
			if err := db.Model(&File{}).Where("node_id = ? AND file_sign = ?  AND deleted_on = ?", node_id, file_sign, 0).Update("disabled", 0).Error; err != nil {
				return err
			}
		} else {
			if err := db.Model(&File{}).Where("node_id = ? AND file_sign = ?  AND deleted_on = ?", node_id, file_sign, 0).Update("disabled", 1).Error; err != nil {
				return err
			}
		}
		return nil

	}
	return errors.New(errno.Msg[errno.FileIsNotExist])

}

//删除文件数据
func DeleteFile(node_id int, file_hash string) error {
	if err := db.Where("node_id = ? AND file_hash = ?", node_id, file_hash).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}

//永久删除节点数据
func DeleteFile_Unscoped(node_id int, file_hash string) error {
	// Unscoped 返回所有记录，包含软删除的记录
	if err := db.Unscoped().Where("node_id = ? AND file_hash = ? AND deleted_on != ?", node_id, file_hash, 0).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}

//删除文件数据
func DeleteFileBySign(node_id int, file_sign string) error {
	if err := db.Where("node_id = ? AND  file_sign = ?", node_id, file_sign).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}

//永久删除节点数据
func DeleteFileBySign_Unscoped(file_sign string) error {
	// Unscoped 返回所有记录，包含软删除的记录
	if err := db.Unscoped().Where(" file_sign = ? AND deleted_on != ?", file_sign, 0).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}

//永久删除所有过期数据
func DeleteFiles() error {
	// Unscoped 返回所有记录，包含软删除的记录
	if err := db.Unscoped().Where("deleted_on != ?", 0).Delete(&File{}).Error; err != nil {
		return err
	}
	return nil
}
