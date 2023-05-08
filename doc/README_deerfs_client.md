English | <a href="https://github.com/xssed/deerfs/blob/master/doc/README_zh_deerfs_client.md" target="_blank">中文</a>



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



To use chunks upload, run deerfs_client.    
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


Example 1: chunks upload(Note that the unit of cut_size is KB)
```
deerfs_client -file_path ./temp/2.mp4 -deerfs_address http://127.0.0.1:7727/  -cut_size 5120 -upload_form_field upload  -http_request_timeout 8000
```
Return:
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

Example 2: General upload
```
deerfs_client -file_path ./temp/2023.gif -deerfs_address http://127.0.0.1:7727/
```
Return:
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


