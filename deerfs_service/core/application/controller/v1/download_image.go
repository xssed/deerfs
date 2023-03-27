package v1

import (
	"fmt"
	"path"

	//"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"

	//"github.com/xssed/deerfs/deerfs_service/core/application/controller"
	"github.com/xssed/deerfs/deerfs_service/core/application/controller/vm"
	"github.com/xssed/deerfs/deerfs_service/core/common"
	"github.com/xssed/deerfs/deerfs_service/core/system/config"

	//"github.com/xssed/deerfs/deerfs_service/core/system/errno"
	"github.com/xssed/deerfs/deerfs_service/core/application/model/mysql_model"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage/image_handler"
	"github.com/xssed/deerfs/deerfs_service/core/system/global"
	"github.com/xssed/deerfs/deerfs_service/core/system/loger"
	"go.uber.org/zap"
)

//图片处理主引导控制器
func ImageMainHandler(c *gin.Context, file *mysql_model.File) {

	//判断是否是图片
	if file_manage.IsImage(file.FileAddr) {

		//获取文件类型
		f_type, _ := file_manage.CheckFileType(file.FileAddr)
		//能进入图片处理的都是对图片源文件进行修改的,根据条件查找缓存
		_, uri_md5, head_uri := file_manage.GetFileCachePathString(c.Request.RequestURI)
		//fmt.Println(uri_md5)  //调试
		//fmt.Println(head_uri) //调试

		//创建缓存文件路径字符串
		cacheFileFolder := path.Join(config.FileStorageDirectoryPath(), "cache", head_uri)
		cacheFilePath := path.Join(config.FileStorageDirectoryPath(), "cache", head_uri, common.JoinString(uri_md5, ".", f_type))
		//判断缓存文件是否存在
		exist, _ := file_manage.PathExists(cacheFilePath)
		//fmt.Println("exist:", exist) //调试
		//存在则直接将文件数据返回
		if exist {
			c.File(cacheFilePath)
			return
		} else {
			//没有存在数据则创建新的缓存
			createFolder_err := file_manage.CreateFolder(cacheFileFolder) //创建存储文件夹
			if createFolder_err != nil {
				fmt.Println("Failed to create storage directory:", createFolder_err.Error())
				loger.Lg.Error("Failed to create storage directory:", zap.String("error", createFolder_err.Error()))
			}

			//根据图像类型去选择不同的图像操作器
			switch f_type {
			case "gif":
				//Gif图片处理控制器
				ImageGifHandler(c, file, cacheFilePath)
			default:
				//图片处理控制器
				ImageHandler(c, file, f_type, cacheFilePath)
			}

		}

	}

	//无有效指令,还是直接输出原文件
	c.File(path.Join(file.FileAddr))
	return
}

//图片处理控制器
func ImageHandler(c *gin.Context, file *mysql_model.File, imgType, cacheFilePath string) {

	//打开图像
	img, err := image_handler.OpenImage(file.FileAddr)
	if err != nil {
		fmt.Println("OpenImage error:", zap.String("error", err.Error()))
		loger.Lg.Error("OpenImage error:", zap.String("error", err.Error()))
		//打开图像失败,直接输出原文件
		c.File(path.Join(file.FileAddr))
		return
	}

	//更改图像尺寸,优先级最高
	_, w_exist := c.GetQuery("w") //宽度
	_, h_exist := c.GetQuery("h") //高度
	if w_exist || h_exist {
		w := com.StrTo(c.Query("w")).MustInt()
		h := com.StrTo(c.Query("h")).MustInt()
		img = image_handler.Resize(img, w, h)
	}

	//根据图像中心来裁剪图像
	_, crop_c_w_exist := c.GetQuery("crop_c_w") //裁剪宽度
	_, crop_c_h_exist := c.GetQuery("crop_c_h") //裁剪高度
	if crop_c_w_exist && crop_c_h_exist {
		crop_c_w := com.StrTo(c.Query("crop_c_w")).MustInt()
		crop_c_h := com.StrTo(c.Query("crop_c_h")).MustInt()
		img = image_handler.CropCenter(img, crop_c_w, crop_c_h)
	}

	//缩略图
	_, thumbnail_w_exist := c.GetQuery("thumbnail_w") //缩略图宽度
	_, thumbnail_h_exist := c.GetQuery("thumbnail_h") //缩略图高度
	if thumbnail_w_exist && thumbnail_h_exist {
		thumbnail_w := com.StrTo(c.Query("thumbnail_w")).MustInt()
		thumbnail_h := com.StrTo(c.Query("thumbnail_h")).MustInt()
		img = image_handler.Thumbnail(img, thumbnail_w, thumbnail_h)
	}

	//锐化
	_, sharpen_exist := c.GetQuery("sharpen")
	if sharpen_exist {
		sharpen := com.StrTo(c.Query("sharpen")).MustFloat64()
		img = image_handler.Sharpen(img, sharpen)
	}

	//调整伽玛值
	_, gamma_exist := c.GetQuery("gamma")
	if gamma_exist {
		gamma := com.StrTo(c.Query("gamma")).MustFloat64()
		img = image_handler.AdjustGamma(img, gamma)
	}

	//调整亮度
	_, brightness_exist := c.GetQuery("brightness")
	if brightness_exist {
		brightness := com.StrTo(c.Query("brightness")).MustFloat64()
		img = image_handler.AdjustBrightness(img, brightness)
	}

	//调整饱和度
	_, saturation_exist := c.GetQuery("saturation")
	if saturation_exist {
		saturation := com.StrTo(c.Query("saturation")).MustFloat64()
		img = image_handler.AdjustSaturation(img, saturation)
	}

	//调整图像对比度
	_, contrast_exist := c.GetQuery("contrast")
	if contrast_exist {
		contrast := com.StrTo(c.Query("contrast")).MustFloat64()
		img = image_handler.AdjustContrast(img, contrast)
	}

	//调整图像非线性对比度
	_, sigmoid_exist := c.GetQuery("sigmoid_midpoint")
	if sigmoid_exist {
		midpoint := com.StrTo(c.Query("sigmoid_midpoint")).MustFloat64()
		factor := com.StrTo(c.Query("sigmoid_factor")).MustFloat64()
		img = image_handler.AdjustSigmoid(img, midpoint, factor)
	}

	//水平翻转图像（从左到右）
	_, flip_h_exist := c.GetQuery("flip_h")
	if flip_h_exist {
		img = image_handler.FlipH(img)
	}

	//垂直翻转图像（从上到下）
	_, flip_v_exist := c.GetQuery("flip_v")
	if flip_v_exist {
		img = image_handler.FlipV(img)
	}

	//图像逆时针旋转180度
	_, rotate180_exist := c.GetQuery("rotate180")
	if rotate180_exist {
		img = image_handler.Rotate180(img)
	}

	//图像逆时针旋转270度
	_, rotate270_exist := c.GetQuery("rotate270")
	if rotate270_exist {
		img = image_handler.Rotate270(img)
	}

	//图像逆时针旋转90度
	_, rotate90_exist := c.GetQuery("rotate90")
	if rotate90_exist {
		img = image_handler.Rotate90(img)
	}

	//水平翻转图像并逆时针旋转90度
	_, transpose_exist := c.GetQuery("transpose")
	if transpose_exist {
		img = image_handler.Transpose(img)
	}

	//垂直翻转图像，逆时针旋转90度
	_, transverse_exist := c.GetQuery("transverse")
	if transverse_exist {
		img = image_handler.Transverse(img)
	}

	//灰度
	_, grayscale_exist := c.GetQuery("grayscale")
	if grayscale_exist {
		img = image_handler.Grayscale(img)
	}

	//反转
	_, invert_exist := c.GetQuery("invert")
	if invert_exist {
		img = image_handler.Invert(img)
	}

	//模糊
	_, blur_exist := c.GetQuery("blur")
	if blur_exist {
		blur := com.StrTo(c.Query("blur")).MustFloat64()
		img = image_handler.Blur(img, blur)
	}

	//====水印部分开始====
	_, wmt_exist := c.GetQuery("watermark_text")  //开启文字水印
	_, wmi_exist := c.GetQuery("watermark_image") //开启图片水印

	//给图片打文字水印
	_, t_exist := c.GetQuery("text")
	wmt := image_handler.WaterMarkText{}
	var wm_img []byte

	if wmt_exist && t_exist && wmi_exist == false {

		r_text := c.Query("text") //将内容进行base64解密
		text, de64_err := common.Base64_Decode(r_text)

		//文字解码成功
		if de64_err == nil {

			fi := image_handler.FontInfo{}

			quality := com.StrTo(c.Query("q")).MustInt()
			wmt.SetParam(quality)

			fontId := com.StrTo(c.Query("font")).MustInt()
			rgba_input := c.Query("rgba") //默认白色，"255_255_255_0"
			fontSize := com.StrTo(c.Query("size")).MustFloat64()
			dpi := com.StrTo(c.Query("dpi")).MustFloat64()
			//===处理默认水印位置开始===
			//pos
			_, pos_exist := c.GetQuery("pos")
			var position int
			if pos_exist {
				position = com.StrTo(c.Query("pos")).MustInt()
				if position < 0 || position > 4 {
					position = 3 //设定默认值为右下角
				}
			} else {
				position = 3
			}
			//x
			_, x_exist := c.GetQuery("x")
			var x int
			if x_exist {
				x = com.StrTo(c.Query("x")).MustInt()
				if x <= 0 {
					x = 37
				}
			} else {
				x = 37
			}
			//y
			_, y_exist := c.GetQuery("y")
			var y int
			if y_exist {
				y = com.StrTo(c.Query("y")).MustInt()
				if y <= 0 {
					y = 10
				}
			} else {
				y = 10
			}
			//===处理默认水印位置结束===
			message := string(text)
			fi.SetParam(fontId, rgba_input, fontSize, dpi, x, y, position, message)

			fi_list := []image_handler.FontInfo{fi}

			wm_text_img, err := wmt.NewTextImageWaterMark(img, imgType, fi_list)
			wm_img = wm_text_img
			if err != nil {
				fmt.Println("Error creating text watermark:", zap.String("error", de64_err.Error()))
				loger.Lg.Error("Error creating text watermark:", zap.String("error", err.Error()))
				//添加水印失败,直接输出原文件
				c.File(path.Join(file.FileAddr))
				return
			}

		} else {
			fmt.Println("Text Watermark Base64_Decode error:", zap.String("error", de64_err.Error()))
			loger.Lg.Error("Text Watermark Base64_Decode error:", zap.String("error", de64_err.Error()))
			//水印文字Base64解码失败,输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

	}

	//图片水印部分
	_, wmi_id_exist := c.GetQuery("wmi_id")

	if wmi_exist && wmi_id_exist && wmt_exist == false {

		wmi_id := c.Query("wmi_id") //获取水印图片的标识

		// 构造结构体查询水印图片
		vmFile := vm.File{FileSign: wmi_id}
		//查询水印图片文件数据
		wmifile, err := vmFile.GetToSign()
		if err != nil {
			fmt.Println("get watermark file database info error:", zap.String("error", err.Error()))
			loger.Lg.Error("get watermark file database info error:", zap.String("error", err.Error()))
			//水印图片文件查询失败,直接输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}
		//判断水印图片文件是否存在
		if wmifile == nil || wmifile.ID < 1 {
			//水印图片文件不存在,直接输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

		//进行水印处理的参数接收
		padding := com.StrTo(c.Query("pad")).MustInt() //偏移多少个像素
		//===处理默认水印位置开始===
		//pos
		_, pos_exist := c.GetQuery("pos")
		var position int
		if pos_exist {
			position := com.StrTo(c.Query("pos")).MustInt() //水印位置
			if position < 0 || position > 4 {
				position = 3 //设定默认值为右下角
			}
		} else {
			position = 3
		}
		//===处理默认水印位置结束===

		//读取水印图片，进行修改操作
		w, err := image_handler.NewFromFile(wmifile.FileAddr, padding, position)
		if err != nil {
			//操作水印原图失败
			fmt.Println("Image Watermark action error:", zap.String("error", err.Error()))
			loger.Lg.Error("Image Watermark action error:", zap.String("error", err.Error()))
			//输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

		//读取原图，设定水印
		temp_img, err := w.MarkForImage(img, imgType)
		wm_img = temp_img
		if err != nil {
			fmt.Println("Original Image Watermark action error:", zap.String("error", err.Error()))
			loger.Lg.Error("Original Image Watermark action error:", zap.String("error", err.Error()))
			//失败,输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

	}
	//====水印部分结束====

	//保存文件
	var save_err error
	if wm_img != nil {
		save_err = wmt.Save(wm_img, cacheFilePath) //统一保存
	} else {
		save_err = image_handler.Save(img, cacheFilePath)
	}
	if save_err != nil {
		fmt.Println("Failed to save image:", zap.String("error", save_err.Error()))
		loger.Lg.Error("Failed to save image:", zap.String("error", save_err.Error()))
		//保存失败,直接输出原文件
		c.File(path.Join(file.FileAddr))
		return
	} else {
		//缓存文件保存成功后累计当前缓存文件目录使用量
		global.CacheFile_Storage_Use_Size = global.CacheFile_Storage_Use_Size + file_manage.FileSize(cacheFilePath)
	}

	//输出缓存文件
	c.File(cacheFilePath)
	return

}

//GIF图片处理控制器
func ImageGifHandler(c *gin.Context, file *mysql_model.File, cacheFilePath string) {

	_, wt_exist := c.GetQuery("watermark_text")
	_, wmi_exist := c.GetQuery("watermark_image")

	//给图片打文字水印
	_, t_exist := c.GetQuery("text")

	wm := image_handler.WaterMarkText{}

	if wt_exist && t_exist { //&& wmi_exist == false

		r_text := c.Query("text") //将内容进行base64解密
		text, de64_err := common.Base64_Decode(r_text)

		//文字解码成功
		if de64_err == nil {

			fi := image_handler.FontInfo{}

			quality := com.StrTo(c.Query("q")).MustInt()
			wm.SetParam(quality)

			fontId := com.StrTo(c.Query("font")).MustInt()
			rgba_input := c.Query("rgba") //"255_255_255_0"
			fontSize := com.StrTo(c.Query("size")).MustFloat64()
			dpi := com.StrTo(c.Query("dpi")).MustFloat64()
			x := com.StrTo(c.Query("x")).MustInt()
			y := com.StrTo(c.Query("y")).MustInt()
			position := com.StrTo(c.Query("posi")).MustInt()
			message := string(text)
			fi.SetParam(fontId, rgba_input, fontSize, dpi, x, y, position, message)

			fi_list := []image_handler.FontInfo{fi}

			new_file, err := wm.NewTextGifWaterMark(file.FileAddr, fi_list)
			if err != nil {
				fmt.Println("Error creating Gif text watermark:", zap.String("error", de64_err.Error()))
				loger.Lg.Error("Error creating Gif text watermark:", zap.String("error", err.Error()))
				//添加水印失败,直接输出原文件
				c.File(path.Join(file.FileAddr))
				return
			}

			err = wm.Save(new_file, cacheFilePath)
			if err != nil {
				loger.Lg.Error("Failed to save gif image:", zap.String("error", err.Error()))
				//保存失败,直接输出原文件
				c.File(path.Join(file.FileAddr))
				return
			}

		} else {
			fmt.Println("Text Watermark Base64_Decode error:", zap.String("error", de64_err.Error()))
			loger.Lg.Error("Text Watermark Base64_Decode error:", zap.String("error", de64_err.Error()))
			//水印文字Base64解码失败,输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

	}

	//给图片打图片水印
	_, wmi_id_exist := c.GetQuery("wmi_id")

	if wmi_exist && wmi_id_exist { //&& wt_exist == false

		wmi_id := c.Query("wmi_id")

		// 构造结构体查询水印图片
		vmFile := vm.File{FileSign: wmi_id}
		//查询水印图片文件数据
		wmifile, err := vmFile.GetToSign()
		if err != nil {
			loger.Lg.Error("get watermark file database info error:", zap.String("error", err.Error()))
			//水印图片文件查询失败,直接输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}
		//判断水印图片文件是否存在
		if wmifile == nil || wmifile.ID < 1 {
			//水印图片文件不存在,直接输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

		//进行水印处理的参数接收
		padding := com.StrTo(c.Query("pad")).MustInt() //偏移多少个像素
		//水印位置
		var position int
		_, pos_exist := c.GetQuery("pos")
		if pos_exist {
			position = com.StrTo(c.Query("pos")).MustInt()
		} else {
			position = 3
		}

		//读取水印图片，进行修改操作
		w, err := image_handler.NewFromFile(wmifile.FileAddr, padding, position)
		if err != nil {
			//操作水印原图失败
			fmt.Println("Original Image Watermark action error:", zap.String("error", err.Error()))
			loger.Lg.Error("Original Image Watermark action error:", zap.String("error", err.Error()))
			//输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}

		//读取原图，设定水印信息
		img, err := w.MarkFile(file.FileAddr)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("Image Watermark action error:", zap.String("error", err.Error()))
			loger.Lg.Error("Image Watermark action error:", zap.String("error", err.Error()))
			//失败,输出原文件
			c.File(path.Join(file.FileAddr))
			return
		}
		//水印添加成功后保存水印图
		err = w.Save(img, cacheFilePath)
		if err != nil {
			loger.Lg.Error("Failed to save Gif Watermark:", zap.String("error", err.Error()))
			//保存失败,直接输出原文件
			c.File(path.Join(file.FileAddr))
			return
		} else {
			//缓存文件保存成功后累计当前缓存文件目录使用量
			global.CacheFile_Storage_Use_Size = global.CacheFile_Storage_Use_Size + file_manage.FileSize(cacheFilePath)
		}

	}

	//输出缓存文件
	c.File(cacheFilePath)
	return

}
