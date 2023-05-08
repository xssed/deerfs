<a href="https://github.com/xssed/deerfs/blob/master/doc/README_deerfs_client.md" target="_blank">English</a> | 中文



## 简单使用    
  deerfs使用Http协议上传和下载文件。   
  假设我们现在deerfs的服务地址为127.0.0.1:7727，它依赖的owlcache节点是127.0.0.1:7721 

### ✺上传    
```  
http://127.0.0.1:7727/deerfs_upload/upload
```  

客户端简单一般上传demo.html    
```  
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
```  
返回内容如下(deerfs作为owlcache的组件，使用了相同的响应结构)
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
请记住响应数据的"Key"和"Data"。查询文件信息和下载文件需要用到这两个信息。    



如果要使用分块上传可以使用客户端deerfs_client    
```     
deerfs_client -help
Usage of deerfs_client:
  -cut_size int
        Input cut size.Note: The cut is in KB for ease of use.
  -deerfs_address string
        Input deerfs address.
  -file_path string
        Input file path.
  -http_request_timeout uint
        Input http request timeout(Millisecond). (default 10000)
  -upload_form_field string
        Input upload form field. (default "upload")   
```      

使用示例1:分块上传
```
deerfs_client -file_path ./temp/2.mp4 -deerfs_address http://127.0.0.1:7727/  -cut_size 5120 -upload_form_field upload  -http_request_timeout 8000
```
返回:
```
Welcome to use simple deerfs client.
Current deerfs client version: 0.1
File path: ./temp/2.mp4
File name: 2.mp4
Cut file size: 5242880  byte.
The address of deerfs: http://127.0.0.1:7727
The upload form field: upload
The http request timeout: 8000  Millisecond.
Start cutting file......
33 / 33 [----------------------------------------------------------------------------------------------] 100.00% 52 p/s
Request temp storage from the server...
Sending chunks files to the server...
33 / 33 [----------------------------------------------------------------------------------------------] 100.00% 11 p/s
Send the merge file command to the server...
Upload success!
Resource address:http://127.0.0.1:7727/28923d5f4cf981a87fa77f71f78f4476MxN2O4MxO5O5N2MxMyLtbtcwN0LtQCcwTMbu
```

使用示例2:一般上传
```
deerfs_client -file_path ./temp/2023.gif -deerfs_address http://127.0.0.1:7727/
```
返回:
```
be careful! You did not input data for the 'cut_size' option, and the client will use the normal upload method.
Welcome to use simple deerfs client.
Current deerfs client version: 0.1
File path: ./temp/2023.gif
File name: 2023.gif
Cut file size: 0  byte.
The address of deerfs: http://127.0.0.1:7727
The upload form field: upload
The http request timeout: 10000  Millisecond.
Send ./temp/2023.gif data to the url http://127.0.0.1:7727 ......
Upload success!
Resource address:http://127.0.0.1:7727/283f1291f3fab643b4dfbaaedb0cad7fMyMxN3MzN3O5MwLtZnapZmLtQCcwTMbu
```  


