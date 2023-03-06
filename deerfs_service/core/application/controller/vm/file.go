package vm

import (
	"errors"
	"fmt"

	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/application/model/owlcache_model"
	"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

//File
type File struct {
	ID         int64
	NodeId     int
	Node       Node
	FileHash   string
	FileSign   string
	FileName   string
	FileExt    string
	FileType   string
	FileMime   string
	FileSize   int64
	FileAddr   string
	Disabled   int
	Ext1       int
	Ext2       string
	CreatedBy  string
	ModifiedBy string
	PageNum    int
	PageSize   int
}

//通过文件MD5值判断文件是否存在,需要NodeID
func (f *File) HasHash() (bool, error) {
	return mysql_model.HasFileToHash(f.NodeId, f.FileHash, f.FileType, f.FileSize)
}

//通过文件Sign值判断文件是否存在
func (f *File) HasSign() (bool, error) {
	return mysql_model.HasFileToSign(f.FileSign)
}

//判断被软删除的文件是否存在,需要NodeID
func (f *File) HasDelFileHash() (bool, error) {
	return mysql_model.HasDelFileToHash(f.NodeId, f.FileHash)
}

//判断被软删除的文件是否存在
func (f *File) HasDelFileSign() (bool, error) {
	return mysql_model.HasDelFileToSign(f.FileSign)
}

//通过文件MD5值判断文件是否存在,需要NodeID,动态传值
func (f *File) HasHashFlexible(node_id int, file_hash, file_type string, file_size int64) (bool, error) {
	return mysql_model.HasFileToHash(node_id, file_hash, file_type, file_size)
}

//通过文件MD5值判断在所有节点中该文件是否存在,动态传值
func (f *File) HasHashFlexibleByAll(file_hash, file_type string, file_size int64) (bool, error) {
	return mysql_model.HasFileToHashByAll(file_hash, file_type, file_size)
}

//通过ID判断文件是否存在
func (f *File) HasID() (bool, error) {
	return mysql_model.HasFileToId(f.ID)
}

//添加文件数据
func (f *File) Add() error {

	data := map[string]interface{}{
		// 关系：File 拥有 Node
		"node_id":    f.NodeId,
		"file_hash":  f.FileHash,
		"file_sign":  f.FileSign,
		"file_name":  f.FileName,
		"file_ext":   f.FileExt,
		"file_type":  f.FileType,
		"file_mime":  f.FileMime,
		"file_size":  f.FileSize,
		"file_addr":  f.FileAddr,
		"disabled":   f.Disabled,
		"ext1":       f.Ext1,
		"ext2":       f.Ext2,
		"created_by": f.CreatedBy,
	}

	return mysql_model.AddFile(data)
}

//修改文件数据，根据文件MD5值
func (f *File) EditForFileHash(node_id int, file_hash string) error {

	data := map[string]interface{}{
		// 关系：File 拥有 Node
		"node_id":     f.NodeId,
		"file_hash":   f.FileHash,
		"file_sign":   f.FileSign,
		"file_name":   f.FileName,
		"file_ext":    f.FileExt,
		"file_type":   f.FileType,
		"file_mime":   f.FileMime,
		"file_size":   f.FileSize,
		"file_addr":   f.FileAddr,
		"disabled":    f.Disabled,
		"ext1":        f.Ext1,
		"ext2":        f.Ext2,
		"modified_by": f.ModifiedBy,
	}

	return mysql_model.EditFileForHash(node_id, file_hash, data)

}

//删除文件数据,软删除
func (f *File) DeleteForFileHash() error {
	return mysql_model.DeleteFile(f.NodeId, f.FileHash)
}

//删除文件数据,彻底删除
func (f *File) DeleteFile_Unscoped() error {
	return mysql_model.DeleteFile_Unscoped(f.NodeId, f.FileHash)
}

//更改文件可访问状态，状态(可用0/禁用1)，禁止非法文件访问,通常用在色情暴力凶杀等非法文件一键封杀
func (f *File) DisabledFile() error {
	return mysql_model.DisabledFile(f.NodeId, f.FileHash)
}

//删除文件数据,软删除
func (f *File) DeleteForFileSign() error {
	return mysql_model.DeleteFileBySign(f.NodeId, f.FileSign)
}

//删除文件数据,彻底删除
func (f *File) DeleteFileBySign_Unscoped() error {
	return mysql_model.DeleteFileBySign_Unscoped(f.FileSign)
}

//更改文件可访问状态，状态(可用0/禁用1)，禁止非法文件访问,通常用在色情暴力凶杀等非法文件一键封杀
func (f *File) DisabledFileBySign() error {
	return mysql_model.DisabledFileBySign(f.NodeId, f.FileSign)
}

//统计该节点中存在的文件数，包含被禁用的文件资源,不包含被删除的文件
func (f *File) Count() (int, error) {
	return mysql_model.GetFilesCount(f.toMap())
}

//获取所有该节点中存在的文件，包含被禁用的文件资源,不包含被删除的文件
func (f *File) GetAll() ([]*mysql_model.File, error) {
	var files []*mysql_model.File
	files, err := mysql_model.GetFiles(f.PageNum, f.PageSize, f.toMap())
	if err != nil {
		return nil, err
	}
	return files, nil
}

//获取单条数据
func (f *File) GetToHash() (*mysql_model.File, error) {
	var file *mysql_model.File
	file, err := mysql_model.GetFileToHash(f.NodeId, f.FileHash)
	if err != nil {
		return nil, err
	}
	return file, nil
}

//获取单条数据
func (f *File) GetToSign() (*mysql_model.File, error) {

	owl_file, _ := owlcache_model.GetFileToSign(f.FileSign) //先从owlcache中获取获取文件信息
	//判断文件是否存在
	if owl_file == nil || owl_file.ID < 1 {
		//不存在就从数据库获取文件信息
		my_file, _ := mysql_model.GetFileToSign(f.FileSign)
		//判断文件是否存在
		if my_file == nil || my_file.ID < 1 {
			//数据库中也不存在
			return nil, errors.New(errno.Msg[errno.FileIsNotExist])
		} else {
			//数据库中存在，则存入owlcache,弱反应文件信息是否存入成功
			save_file_err := owlcache_model.SetFileToSign(f.FileSign, my_file)
			if save_file_err != nil {
				fmt.Println("owlcache_model.SetFileToSign error:", zap.String("error", save_file_err.Error()))
				loger.Lg.Error("owlcache_model.SetFileToSign error:", zap.String("error", save_file_err.Error()))
			}
			return my_file, nil

		}

	} else {
		return owl_file, nil
	}

}

//查询条件
func (f *File) toMap() map[string]interface{} {
	m := make(map[string]interface{})
	m["deleted_on"] = 0
	m["node_id"] = global.Node_ID
	return m
}
