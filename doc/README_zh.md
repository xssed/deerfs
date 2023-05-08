<a href="https://github.com/xssed/deerfs" target="_blank">English</a> | 中文简介


<div align="center">

# 🦌 deerfs

![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/deer.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/deerfs.svg)](https://github.com/xssed/deerfs/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/deerfs.svg?style=popout-square)](https://github.com/xssed/deerfs/releases)

</div>

 🦌 deerfs是owlcache的一个组件扩展。使用它，您可以构建一个简单的无中心分布式文件系统。     

  主项目:<a href="https://github.com/xssed/owlcache" target="_blank"> owlcache</a>    

## deerfs_service    

   服务端文件的普通上传、分块上传、文件信息的集群查询、下载资源时的图像处理(调整尺寸、以图像中心裁剪图像、缩略图、锐化、伽玛值、亮度、饱和度、图像对比度、图像非线性对比度、图像多角度旋转、生成图像的灰度版本、反转、模糊、文字水印、图像水印)。

   <a href="https://github.com/xssed/deerfs/blob/master/doc/README_zh_deerfs_service.md" target="_blank">文档资料</a>   

## deerfs_client    

   客户端文件的普通上传、分块上传。

   <a href="https://github.com/xssed/deerfs/blob/master/doc/README_zh_deerfs_client.md" target="_blank">文档资料</a>

## deerfs需要以下服务的支持
- owlcache
- mysql(数据表文件路径deerfs_service/sql/table.sql)

## 上传与下载的权限
deerfs这边更注重功能的实现，作为一个独立服务，想接入各种不同的平台，各个平台权限的验证是多种多样的，所以需要你定制化的自己来实现。或者是网关来做权限验证这件事。

## 开发与讨论(不接商业合作)
- 联系我📪:xsser@xsser.cc
- 个人主页🛀:https://www.xsser.cc


