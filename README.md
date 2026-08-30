

English | <a href="https://github.com/xssed/deerfs/blob/master/doc/README_zh.md" target="_blank">中文简介</a>

<div align="center">

# 🦌 deerfs

![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/deer.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/deerfs.svg)](https://github.com/xssed/deerfs/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/deerfs.svg?style=popout-square)](https://github.com/xssed/deerfs/releases)

</div>

 🦌 deerfs is a component extension of owlcache. with it you can build a simple non-centralized distributed file system.     

  Main project:<a href="https://github.com/xssed/owlcache" target="_blank"> owlcache</a>     

## deerfs_service    

   General upload of server files, chunks upload, cluster query of file information, image processing when downloading resources (Resize, CropCenter, Thumbnail, Sharpen, Gamma, Brightness, Saturation, Contrast, Sigmoid, Rotate, Transverse, Transpose, Grayscale, invert, blur, Text watermark, Image watermark).

   <a href="https://github.com/xssed/deerfs/blob/master/doc/README_deerfs_service.md" target="_blank">Document</a>   

## deerfs_client    

   General upload and chunks upload of client files. 

   <a href="https://github.com/xssed/deerfs/blob/master/doc/README_deerfs_client.md" target="_blank">Document</a>

## deerfs requires support for the following services
- owlcache
- mysql(Data table file path:deerfs_service/sql/table.sql)

## Upload and download permissions
Deerfs focuses more on functionality. As a standalone service, it needs to access a variety of different platforms. Permission verification is different on each platform, so you need to modify the source code to implement it yourself. Or perform permission verification on the gateway.    

## Development and discussion(not involved in business cooperation)
- Email📪:xsser@xsser.cc
- Homepage🛀:https://www.xsser.cc
