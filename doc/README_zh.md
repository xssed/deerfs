<a href="https://github.com/xssed/deerfs" target="_blank">English</a> | 中文简介


<div align="center">

# 🦌 deerfs

![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/deer.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/deerfs.svg)](https://github.com/xssed/deerfs/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/deerfs.svg?style=popout-square)](https://github.com/xssed/deerfs/releases)

</div>

 🦌 deerfs是owlcache的一个组件扩展。使用它，您可以构建一个简单的无中心分布式文件系统。     

  主项目:<a href="https://github.com/xssed/owlcache" target="_blank"> owlcache</a>    

## 简单使用    
  deerfs使用Http协议上传和下载文件。   
  假设我们现在deerfs的服务地址为127.0.0.1:7727，它依赖的owlcache节点是127.0.0.1:7721 

### ✺上传    
~~~shell
http://127.0.0.1:7727/deerfs_upload/upload
~~~

客户端简单demo    
~~~shell
<html>
<head>
<meta charset="utf-8">
<title></title>
</head>
<body>
<form action="http://127.0.0.1:7727/deerfs_upload/upload" method="post" enctype="multipart/form-data">
    <label for="file">文件：</label>
    <input type="file" name="upload" id="file">
    <!-- name值在配置文件中"[upload]的form_field" --><br>
    <input type="submit" name="submit" value="提交">
</form>
</body>
</html>
~~~
返回内容如下(deerfs作为owlcache的组件，使用了相同的响应结构)
~~~shell
{
    "Cmd": "",
    "Status": 200,
    "Results": "Success",
    "Key": "8b82fd400c864077026d3421d833be95N2MyMyMwN1N2LtcwZkZmLtQCcwTMbu",
    "Data": "http:/127.0.0.1:7727/8b82fd400c864077026d3421d833be95N2MyMyMwN1N2LtcwZkZmLtQCcwTMbu",
    "ResponseHost": "http://127.0.0.1:7727/",
    "KeyCreateTime": "0001-01-01T00:00:00Z"
}
~~~
请记住响应数据的"Key"和"Data"。查询文件信息和下载文件需要用到这两个信息。    


### ✺查询(owlcache)     

💡💡💡💡💡<b>注意查询是在owlcache中进行</b>     

上传成功后文件信息会存储在中owlcache之中(理论上是永久存储，因为owlcache支持数据的磁盘存储，但是不排除主机的忽然断电等特殊情况导致的数据丢失，即使丢失，在下载这个文件时这个文件的信息也会被重新写入owlcache)    

<b>查询的Key格式"deerfs::"+File的Key</b>    

* 单节点查询
~~~shell
http://127.0.0.1:7721/data/?cmd=get&key=deerfs::6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
~~~

* 集群节点查询
~~~shell
http://127.0.0.1:7721/group_data/?cmd=get&key=deerfs::6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
~~~


### ✺下载      

普通下载      
~~~shell
http://127.0.0.1:7727/6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
~~~

## 下载资源时的图像处理 

图像处理时生成的临时文件会被缓存到本地，配置文件中设置的时间内不会重复生成，节省系统资源，并进行自动管理。    
开启图像处理需要在Get请求时添加“action=imageView”参数。   
~~~shell
http://127.0.0.1:7727/File的Key?action=imageView
~~~   

### ✺修改图像     

图像处理类型png、jpg。    
函数按列表顺序执行，优先级依次递减。

| 功能             | 参数、含义、值范围                     | 示例 |
| ---------------- | -------------------------------------- | ---- |
| Resize(调整尺寸) | w、调整后宽度、 0-∞<br>H、调整后高度、0-∞ | &w=300&h=200 |
| CropCenter(以图像中心裁剪图像)             | crop_c_w、调整后宽度、 0-∞<br>crop_c_h、调整后高度、0~∞                                   | &crop_c_w=300&crop_c_h=200 |
| Thumbnail(缩略图)             | thumbnail_w、调整后宽度、 0-∞<br>thumbnail_h、调整后高度、0-∞                                   | &thumbnail_w=300&thumbnail_h=200 |
| Sharpen(锐化)             | sharpen、进行sharpen处理、0.1-∞(不建议过大)                                   | &sharpen=20 |
| Gamma(伽玛值)             | gamma、进行gamma处理、0.1-∞(Gamma=1.0提供原始图像。小于1.0的伽马会使图像变暗，大于1.0的伽玛会使其变亮)                                   | &gamma=0.1 |
| Brightness(亮度)             | brightness、进行brightness处理、-100~100(0表示原始图像)                                   | &brightness=-20 |
| Saturation(饱和度)             | saturation、进行saturation处理、-100-100(0表示原始图像)                                   | &saturation=-20 |
| Contrast(图像对比度)             | contrast、进行contrast处理、-100-100(0表示原始图像)                                   | &contrast=10 |
| Sigmoid(图像非线性对比度,对照片调整有用的非线性对比度变化，因为它保留了高光和阴影细节)             | sigmoid_midpoint 、对比度的中点、0~1(一般为0.5)<br>sigmoid_factor、对比度增加或减少多少、-10-10(参数为正，则图像对比度增加，否则对比度降低)                                   | &sigmoid_midpoint=0.5&sigmoid_factor=10 |
| FlipH(水平翻转图像（从左到右）)             | flip_h、进行FlipH处理、无                                   | &flip_h |
| FlipV(垂直翻转图像（从上到下）)             | flip_v、进行FlipV处理、无                                   | &flip_v |
| Rotate180(图像逆时针旋转180度)             | rotate180、进行Rotate180处理、无                                   | &rotate180 |
| Rotate270(图像逆时针旋转270度)             | rotate270、进行Rotate270处理、无                                   | &rotate270 |
| Rotate90(图像逆时针旋转90度)             | rotate90、进行Rotate90处理、无                                   | &rotate90 |
| Transpose(水平翻转图像并逆时针旋转90度)             | transpose、进行Transpose处理、无                                   | &transpose |
| Transverse(垂直翻转图像，逆时针旋转90度)             | transverse、进行Transpose处理、无                                   | &transverse |
| Grayscale(生成图像的灰度版本)             | grayscale、进行Grayscale处理、无                                   | &grayscale |
| invert(反转)             | invert、进行invert处理、无                                   | &invert |
| blur(模糊)             | blur、进行blur处理、建议1-20                                   | &blur=10 |

### ✺为图像添加水印     

除了gif格式，其它格式是可以先进行图像修改后添加水印。     
添加文字水印的同时不能再添加图片水印只能二选一。    
函数按列表顺序执行，优先级依次递减。     

#### ⚪文字水印    

处理类型png、jpg。    
开启文字水印需要在Get请求时添加“watermark_text”参数。   
~~~shell
http://127.0.0.1:7727/File的Key?action=imageView&watermark_text
~~~   

| 参数             | 含义或值范围                     | 示例 |
| ---------------- | -------------------------------------- | ---- |
| text | 需要添加的文字水印的Base64编码 | &text=5L2g5aW9 |
| q             | 图片质量，1-100(默认100)                                   | &q=100 |
| font             | 设置字体ID,内置了7种字体。<br/>💡如果侵权请联系我更换。<br/> 1，默认字体，方正宋体<br/> 2，文泉驿正黑<br/>  3，方正楷体<br/>  4，濑户字体<br/>  5，Lingxun<br/>  6，Roboto<br/>  7，RobotoSerif<br/>                           | &font=2 |
| rgba             | 设置字体rgba颜色的值<br/>rgba的四个值，之间用“_”分割<br/>默认白色"255_255_255_0"                                   | &rgba=60_179_113_100 |
| size             | 设置字体大小,默认17                                   | &size=20 |
| dpi             | 设置图像DPI,默认75                                   | &dpi=75 |
| pos             | 设置文字水印的相对位置，五个值。<br/>	默认值为右下角。<br/>0，TopLeft <br/>	1，TopRight<br/>	2，BottomLeft<br/>	3，BottomRight<br/>	4，Center                                 | &pos=3 |
| x             | 设置文字水印的偏移值x,用来调整水印位置                                  | &x=35 |
| y             | 设置文字水印的偏移值x,用来调整水印位置                                  | &y=10 |




#### ⚪图片水印   
处理类型png、jpg、gif。
开启图片水印需要在Get请求时添加“watermark_image”参数。   
~~~shell
http://127.0.0.1:7727/File的Key?action=imageView&watermark_image
~~~  

| 参数             | 含义或值范围                     | 示例 |
| ---------------- | -------------------------------------- | ---- |
| wmi_id | 水印图片的Key，请先将水印图片上传至deerfs之中 | &wmi_id=783025f1813323b8530c419b68bb0b3bN0MyMzN2LtaqcwZnLtQCcwTMbu |
| pad             | 偏移多少个像素                                  | &pad=20 |
| pos             | 设置图片水印的相对位置，五个值。<br/>	默认值为右下角。<br/>0，TopLeft <br/>	1，TopRight<br/>	2，BottomLeft<br/>	3，BottomRight<br/>	4，Center                                  | &pos=0 |

====================分割线====================      

测试的图像样本    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu.jpg?raw=true)

图像处理demo1    
~~~shell
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&thumbnail_w=300&thumbnail_h=200&sharpen=20&brightness=20&contrast=10
~~~ 

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo1.jpg?raw=true)

图像处理demo2    
~~~shell
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&crop_c_w=300&crop_c_h=200&sigmoid_midpoint=0.5&sigmoid_factor=10&rotate270&saturation=20
~~~ 

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo2.jpg?raw=true)

图像处理demo3   
~~~shell
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&w=300&h=230&blur=9.5
~~~ 

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo3.jpg?raw=true)

图像添加文字水印demo4   
~~~shell
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&font=4&watermark_text&w=300&h=230&text=5L2g5aW9&q=100&rgba=34_139_34_100&size=20&pos=3&x=30&y=10&dpi=75
~~~ 

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo4.jpg?raw=true)



用于测试的静态图和动态图样本  
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/liangchaowei.gif?raw=true)  
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/caicai.jpg?raw=true)
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/fox.gif?raw=true)

图像添加静态图水印demo5   
~~~shell
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&watermark_image&wmi_id=87a132e181225bb608b25dadea08ddfaMxO5N3N1MzLtaqcwZnLtYjMzVXRE&crop_c_w=400&crop_c_h=320
~~~     
输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo5.jpg?raw=true)

动态图像添加动态图水印demo6  
```
http://127.0.0.1:7727/f6a15b0e95baee0f5a081ea87fa9b3d2MxMxMwMwN2MzMxLtZnapZmLtQCcwTMbu?action=imageView&watermark_image&wmi_id=545264737d51c5e83dc803fcfb30bddcMzN1O4N3N3LtZnapZmLtZmZnREcz&pos=2
```  
输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/liangchaowei_demo6.gif?raw=true)



## deerfs需要以下服务的支持
- owlcache
- mysql(数据表文件路径deerfs_service/sql/table.sql)


## 开发与讨论(不接商业合作)
- 联系我📪:xsser@xsser.cc
- 个人主页🛀:https://www.xsser.cc


