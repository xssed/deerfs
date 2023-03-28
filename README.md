English | <a href="https://github.com/xssed/deerfs/blob/master/doc/README_zh.md" target="_blank">中文简介</a>

<div align="center">

# 🦌 deerfs

![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/deer.jpg?raw=true)

[![License](https://img.shields.io/github/license/xssed/deerfs.svg)](https://github.com/xssed/deerfs/blob/master/LICENSE)
[![release](https://img.shields.io/github/release/xssed/deerfs.svg?style=popout-square)](https://github.com/xssed/deerfs/releases)

</div>

 🦌 deerfs is a component extension of owlcache. with it you can build a simple non-centralized distributed file system.     

  Main project:<a href="https://github.com/xssed/owlcache" target="_blank"> owlcache</a>     

## Simple use    
  deerfs uses Http to upload and download files.   
  Suppose we now have a deerfs service address of 127.0.0.1:7727, and the owlcache node it relies on is 127.0.0.1:7721. 

### ✺Upload    
```  
http://127.0.0.1:7727/deerfs_upload/upload
```  

Simple client demo    
```  
<html>
<head>
<meta charset="utf-8">
<title></title>
</head>
<body>
<form action="http://127.0.0.1:7727/deerfs_upload/upload" method="post" enctype="multipart/form-data">
    <label for="file">file：</label>
    <input type="file" name="upload" id="file">
    <!-- name value in the config file "[upload] form_field" --><br>
    <input type="submit" name="submit" value="submit">
</form>
</body>
</html>
```  
The returned content is as follows: (Deerfs as a component of owlcache, uses the same response struct)
```  
{
    "Cmd": "",
    "Status": 200,
    "Results": "Success",
    "Key": "8b82fd400c864077026d3421d833be95N2MyMyMwN1N2LtcwZkZmLtQCcwTMbu",
    "Data": "http:/127.0.0.1:7727/8b82fd400c864077026d3421d833be95N2MyMyMwN1N2LtcwZkZmLtQCcwTMbu",
    "ResponseHost": "http://127.0.0.1:7727/",
    "KeyCreateTime": "0001-01-01T00:00:00Z"
}
```  
Please remember the "Key" and "Data" of the response data. Querying file information and downloading files require both information.       


### ✺Query(owlcache)     

💡💡💡💡💡<b>Note: The query is performed in owlcache.</b>     

After successful uploading, the file information will be stored in the middle owlcache (theoretically, it is permanent storage, because owlcache supports disk storage of data, but it does not exclude data loss caused by special circumstances such as sudden power failure of the host. Even if the file is lost, the information of the file will be rewritten to owlcache when downloading the file).      

<b>The key format of the query is "deerfs::"+the key of the file.</b>    

* Single node query
```  
http://127.0.0.1:7721/data/?cmd=get&key=deerfs::6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
```  

* Cluster node query
```  
http://127.0.0.1:7721/group_data/?cmd=get&key=deerfs::6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
```  


### ✺Download      

Normal Download      
```  
http://127.0.0.1:7727/6b49865b0d6a3fc51ae372c9545bbc36N0N3N0N0MzN3LtaqcwZnLtZmZnREcz
```  

## Image processing when downloading resources 

Temporary files generated during image processing will be cached locally, and will not be generated repeatedly within the time set in the configuration file, saving system resources, and performing automatic management.    

To enable image processing, you need to add the "action=imageView" parameter to the Get request.       
```  
http://127.0.0.1:7727/File_Key?action=imageView
```  

### ✺Modify Image     

Image processing types: png, jpg.        
The functions are executed in list order with decreasing priority.    

| function             | Parameter, meaning, value range                     | examples |
| ---------------- | -------------------------------------- | ---- |
| Resize | w、Adjusted width、 0-∞<br>H、Height after adjustment、0-∞ | &w=300&h=200 |
| CropCenter(Crop an image at the image center)             | crop_c_w、Adjusted width、 0-∞<br>crop_c_h、Height after adjustment、0~∞                                   | &crop_c_w=300&<br>crop_c_h=200 |
| Thumbnail            | thumbnail_w、Adjusted width、 0-∞<br>thumbnail_h、Height after adjustment、0-∞                                   | &thumbnail_w=300&<br>thumbnail_h=200 |
| Sharpen             | sharpen、Sharpen processing、0.1-∞(Not recommended to be too large)                                   | &sharpen=20 |
| Gamma             | gamma、Gamma processing、0.1-∞(Gamma=1.0 provides the original image. Gamma less than 1.0 darkens the image, and gamma greater than 1.0 brightens it)                                   | &gamma=0.1 |
| Brightness            | brightness、brightness processing、-100~100(0 represents the original image)                                   | &brightness=-20 |
| Saturation             | saturation、saturation processing、-100-100(0 represents the original image)                                   | &saturation=-20 |
| Contrast            | contrast、contrast processing、-100-100(0 represents the original image)                                   | &contrast=10 |
| Sigmoid           | sigmoid_midpoint 、Midpoint of contrast、0~1(Typically 0.5)<br>sigmoid_factor、How much contrast increases or decreases、-10-10(If the parameter is positive, the image contrast increases, otherwise the contrast decreases)                                   | &sigmoid_midpoint=0.5&<br>sigmoid_factor=10 |
| FlipH(Flip the image horizontally (from left to right))             | flip_h、FlipH processing、not have                                   | &flip_h |
| FlipV(Flip the image vertically (from top to bottom))             | flip_v、FlipV processing、not have                                   | &flip_v |
| Rotate180(Image rotates 180 degrees counterclockwise)             | rotate180、Rotate180 processing、not have                                   | &rotate180 |
| Rotate270(Image rotates 270 degrees counterclockwise)             | rotate270、Rotate270 processing、not have                                   | &rotate270 |
| Rotate90(Image rotates 90 degrees counterclockwise)             | rotate90、Rotate90 processing、not have                                   | &rotate90 |
| Transpose(Flip the image horizontally and rotate it 90 degrees counterclockwise)             | transpose、Transpose processing、not have                                   | &transpose |
| Transverse(Flip the image vertically and rotate it 90 degrees counterclockwise)             | transverse、Transverse processing、not have                                   | &transverse |
| Grayscale(Generate a grayscale version of the image)             | grayscale、Grayscale processing、not have                                   | &grayscale |
| invert            | invert、invert processing、not have                                   | &invert |
| blur            | blur、blur processing、propose 1-20                                   | &blur=10 |

### ✺Add a watermark to an image     

In addition to the gif format, other formats can be modified before adding a watermark.       
When adding a text watermark, you cannot add another image watermark. You can only choose between two options.     
The functions are executed in list order with decreasing priority.      

#### ⚪Text Watermark    

Processing types: png, jpg.    
To enable text watermarks, you need to add the "watermark_text" parameter to the Get request.   
```  
http://127.0.0.1:7727/File_Key?action=imageView&watermark_text
```  

| function             | Parameter, meaning, value range                     | examples |
| ---------------- | -------------------------------------- | ---- |
| text | Base64 encoding of text watermarks to be added | &text=5L2g5aW9 |
| q             | Image quality, 1-100 (default 100)                                   | &q=100 |
| font             | Set font ID, with 7 built-in fonts.<br/>💡If infringement occurs, please contact me for replacement.<br/> 1，Default font, 方正宋体<br/> 2，文泉驿正黑<br/>  3，方正楷体<br/>  4，濑户字体<br/>  5，Lingxun<br/>  6，Roboto<br/>  7，RobotoSerif<br/>                           | &font=2 |
| rgba             | Set the value of font rgba color.<br/>The four values of rgba are separated by "_".<br/>Default White "255_255_255_0"                                   | &rgba=60_179_113_100 |
| size             | Set font size, default 17                                   | &size=20 |
| dpi             | Set image DPI, default 75                                   | &dpi=75 |
| pos             | Set the relative position of the text watermark, with five values.<br/>	The default value is the lower right corner.<br/>0，TopLeft <br/>	1，TopRight<br/>	2，BottomLeft<br/>	3，BottomRight<br/>	4，Center                                 | &pos=3 |
| x             | Set the offset value x of the text watermark to adjust the watermark position.                                  | &x=35 |
| y             | Set the offset value x of the text watermark to adjust the watermark position.                                  | &y=10 |




#### ⚪Image Watermark   
Processing types: png, jpg, gif.    
To enable image watermarks, you need to add the "watermark_image" parameter to the Get request.   
```  
http://127.0.0.1:7727/File_Key?action=imageView&watermark_image
```  

| function             | Meaning or value range                     | examples |
| ---------------- | -------------------------------------- | ---- |
| wmi_id | The key of the watermark image, please upload the watermark image to deerfs first. | &wmi_id=783025f1813323b8530c419b<br>68bb0b3bN0MyMzN2LtaqcwZnLtQCcwTMbu |
| pad             | How many pixels are offset                                  | &pad=20 |
| pos             | Set the relative position of the image watermark, with five values.<br/>	The default value is the lower right corner.<br/>0，TopLeft <br/>	1，TopRight<br/>	2，BottomLeft<br/>	3，BottomRight<br/>	4，Center                                  | &pos=0 |

---      

Image samples tested   
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu.jpg?raw=true)

Image processing demo1   
```  
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&thumbnail_w=300&thumbnail_h=200&sharpen=20&brightness=20&contrast=10
```  

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo1.jpg?raw=true)

Image processing demo2   
```  
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&crop_c_w=300&crop_c_h=200&sigmoid_midpoint=0.5&sigmoid_factor=10&rotate270&saturation=20
```  

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo2.jpg?raw=true)

Image processing demo3  
```  
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&w=300&h=230&blur=9.5
```  

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo3.jpg?raw=true)

Image adding text watermark demo4   
```  
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&font=4&watermark_text&w=300&h=230&text=5L2g5aW9&q=100&rgba=34_139_34_100&size=20&pos=3&x=30&y=10&dpi=75
```  

输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo4.jpg?raw=true)



Sample static and dynamic graphs for testing  
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/liangchaowei.gif?raw=true)  
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/caicai.jpg?raw=true)
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/fox.gif?raw=true)

Image added static graph watermark demo5  
```  
http://127.0.0.1:7727/733f67e13a8d770f19a9be203a19bdf2MyN2N0N0O4LtaqcwZnLtQCcwTMbu?action=imageView&watermark_image&wmi_id=87a132e181225bb608b25dadea08ddfaMxO5N3N1MzLtaqcwZnLtYjMzVXRE&crop_c_w=400&crop_c_h=320
```  
输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/baichunlu_demo5.jpg?raw=true)

Dynamic Image added Dynamic Image watermark demo6 
```
http://127.0.0.1:7727/f6a15b0e95baee0f5a081ea87fa9b3d2MxMxMwMwN2MzMxLtZnapZmLtQCcwTMbu?action=imageView&watermark_image&wmi_id=545264737d51c5e83dc803fcfb30bddcMzN1O4N3N3LtZnapZmLtZmZnREcz&pos=2
```  
输出    
![Image text](https://github.com/xssed/deerfs/blob/master/doc/assets/liangchaowei_demo6.gif?raw=true)



## deerfs requires support for the following services
- owlcache
- mysql(Data table file path:deerfs_service/sql/table.sql)

## Upload and download permissions
Deerfs focuses more on functionality. As a standalone service, it needs to access a variety of different platforms. Permission verification is different on each platform, so you need to modify the source code to implement it yourself. Or perform permission verification on the gateway.    

## Development and discussion(not involved in business cooperation)
- Email📪:xsser@xsser.cc
- Homepage🛀:https://www.xsser.cc



