package image_handler

import (
	"image"

	"github.com/disintegration/imaging"
)

//打开图像
func OpenImage(inFilePath string) (image.Image, error) {

	// Open a image.
	fileImg, err := imaging.Open(inFilePath)
	if err != nil {
		return nil, err
	}
	return fileImg, nil

}

//尺寸调整,Resize
//Resize resizes the image to the specified width and height using the specified resampling filter and returns the transformed image.
//If one of width or height is 0, the image aspect ratio is preserved.
//Resize使用指定的重采样过滤器将图像大小调整为指定的宽度和高度，并返回变换后的图像。
//如果宽度或高度之一为0，则图像纵横比将保留。
func Resize(fileImg image.Image, width, height int) *image.NRGBA {

	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return imaging.Resize(fileImg, width, height, imaging.Lanczos)

}

//以图片中心裁剪
//CropCenter cuts out a rectangular region with the specified size from the center of the image and returns the cropped image.
//CropCenter 从图像中心切出一个指定大小的矩形区域，并返回裁剪后的图像。
func CropCenter(fileImg image.Image, width, height int) *image.NRGBA {

	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return imaging.CropCenter(fileImg, width, height)
}

//缩略图
//Thumbnail scales the image up or down using the specified resample filter, crops it to the specified width and hight and returns the transformed image.
//缩略图使用指定的重采样过滤器放大或缩小图像，将其裁剪到指定的宽度和高度并返回转换后的图像。
func Thumbnail(fileImg image.Image, width, height int) *image.NRGBA {

	if width < 0 {
		width = 0
	}
	if height < 0 {
		height = 0
	}
	return imaging.Thumbnail(fileImg, width, height, imaging.Lanczos)
}

//锐化
//Sharpen produces a sharpened version of the image.
//Sigma parameter must be positive and indicates how much the image will be sharpened.
//锐化生成图像的锐化版本。
//Sigma参数必须为正数，并指示图像的锐化程度。
func Sharpen(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma <= 0 {
		sigma = 0.1
	}
	return imaging.Sharpen(fileImg, sigma)
}

//调整伽玛值
//AdjustGamma performs a gamma correction on the image and returns the adjusted image. Gamma parameter must be positive.
//Gamma = 1.0 gives the original image. Gamma less than 1.0 darkens the image and gamma greater than 1.0 lightens it.
//AdjustGamma对图像执行gamma校正并返回调整后的图像。Gamma参数必须为正值。
//Gamma=1.0提供原始图像。小于1.0的伽马会使图像变暗，大于1.0的伽玛会使其变亮。
func AdjustGamma(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma <= 0 {
		sigma = 1.0
	}
	return imaging.AdjustGamma(fileImg, sigma)
}

//亮度
//AdjustBrightness changes the brightness of the image using the percentage parameter and returns the adjusted image.
//The percentage must be in range (-100, 100). The percentage = 0 gives the original image.
//The percentage = -100 gives solid black image. The percentage = 100 gives solid white image.
//AdjustBrightness使用百分比参数更改图像的亮度，并返回调整后的图像。
//百分比必须在范围（-100至100）内。百分比=0表示原始图像。
//百分比=-100表示实心黑色图像。百分比=100表示纯白图像。
func AdjustBrightness(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma < -100 {
		sigma = -100
	}
	if sigma > 100 {
		sigma = 100
	}
	return imaging.AdjustBrightness(fileImg, sigma)
}

//饱和度
//AdjustSaturation changes the saturation of the image using the percentage parameter and returns the adjusted image.
//The percentage must be in the range (-100, 100).
//The percentage = 0 gives the original image.
//The percentage = 100 gives the image with the saturation value doubled for each pixel.
//The percentage = -100 gives the image with the saturation value zeroed for each pixel (grayscale).
//AdjustSaturation使用百分比参数更改图像的饱和度，并返回调整后的图像。
//百分比必须在（-100，100）范围内。
//百分比=0表示原始图像。
//百分比＝100给出了每个像素的饱和度值加倍的图像。
//百分比=-100表示每个像素（灰度）的饱和度值为零的图像。
func AdjustSaturation(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma < -100 {
		sigma = -100
	}
	if sigma > 100 {
		sigma = 100
	}
	return imaging.AdjustSaturation(fileImg, sigma)
}

//图像对比度
//AdjustContrast changes the contrast of the image using the percentage parameter and returns the adjusted image.
//The percentage must be in range (-100, 100).
//The percentage = 0 gives the original image. The percentage = -100 gives solid gray image.
//AdjustContrast使用百分比参数更改图像的对比度，并返回调整后的图像。
//百分比必须在范围（-100到100）内。
//百分比=0表示原始图像。百分比=-100表示纯灰色图像。
func AdjustContrast(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma < -100 {
		sigma = -100
	}
	if sigma > 100 {
		sigma = 100
	}
	return imaging.AdjustContrast(fileImg, sigma)
}

//图像非线性对比度
//AdjustSigmoid changes the contrast of the image using a sigmoidal function and returns the adjusted image.
//It's a non-linear contrast change useful for photo adjustments as it preserves highlight and shadow detail.
//The midpoint parameter is the midpoint of contrast that must be between 0 and 1, typically 0.5.
//The factor parameter indicates how much to increase or decrease the contrast, typically in range (-10, 10).
//If the factor parameter is positive the image contrast is increased otherwise the contrast is decreased.
//AdjustSigmoid 使用 sigmoidal 函数改变图像的对比度，并返回调整后的图像。
//这是对照片调整有用的非线性对比度变化，因为它保留了高光和阴影细节。
//midpoint参数是对比度的中点，必须在0到1之间，一般为0.5。
//factor参数表示对比度增加或减少多少，一般在(-10, 10)范围内。
//如果因子参数为正，则图像对比度增加，否则对比度降低。
func AdjustSigmoid(fileImg image.Image, midpoint, factor float64) *image.NRGBA {

	if midpoint < 0 {
		midpoint = 0
	}
	if midpoint > 1 {
		midpoint = 1
	}
	if factor < -10 {
		factor = -10
	}
	if factor > 10 {
		factor = 10
	}
	return imaging.AdjustSigmoid(fileImg, midpoint, factor)
}

//水平翻转图像（从左到右）
//FlipH flips the image horizontally (from left to right) and returns the transformed image.
//FlipH 水平翻转图像（从左到右）并返回变换后的图像。
func FlipH(fileImg image.Image) *image.NRGBA {
	return imaging.FlipH(fileImg)
}

//垂直翻转图像（从上到下）
//FlipV flips the image vertically (from top to bottom) and returns the transformed image.
//FlipV 垂直翻转图像（从上到下）并返回转换后的图像。
func FlipV(fileImg image.Image) *image.NRGBA {
	return imaging.FlipV(fileImg)
}

//图像逆时针旋转180度
//Rotate180 rotates the image 180 degrees counter-clockwise and returns the transformed image.
//Rotate180 将图像逆时针旋转 180 度并返回变换后的图像。
func Rotate180(fileImg image.Image) *image.NRGBA {
	return imaging.Rotate180(fileImg)
}

//图像逆时针旋转270度
//Rotate270 rotates the image 270 degrees counter-clockwise and returns the transformed image.
//Rotate270 将图像逆时针旋转 270 度并返回变换后的图像。
func Rotate270(fileImg image.Image) *image.NRGBA {
	return imaging.Rotate270(fileImg)
}

//图像逆时针旋转90度
//Rotate90 rotates the image 90 degrees counter-clockwise and returns the transformed image.
//Rotate90 将图像逆时针旋转 90 度并返回变换后的图像。
func Rotate90(fileImg image.Image) *image.NRGBA {
	return imaging.Rotate90(fileImg)
}

//水平翻转图像并逆时针旋转90度
//Transpose flips the image horizontally and rotates 90 degrees counter-clockwise.
//Transpose 水平翻转图像并逆时针旋转 90 度。
func Transpose(fileImg image.Image) *image.NRGBA {
	return imaging.Transpose(fileImg)
}

//垂直翻转图像，逆时针旋转90度
//Transverse flips the image vertically and rotates 90 degrees counter-clockwise.
//Transverse 垂直翻转图像，逆时针旋转90度。
func Transverse(fileImg image.Image) *image.NRGBA {
	return imaging.Transverse(fileImg)
}

//灰度
//Grayscale produces a grayscale version of the image.
//生成图像的灰度版本。
func Grayscale(fileImg image.Image) *image.NRGBA {
	return imaging.Grayscale(fileImg)
}

//反转
//Invert produces an inverted (negated) version of the image.
//Invert 生成图像的反转（否定）版本。
func Invert(fileImg image.Image) *image.NRGBA {
	return imaging.Invert(fileImg)
}

//模糊
//Blur produces a blurred version of the image using a Gaussian function.
//Sigma parameter must be positive and indicates how much the image will be blurred.
//模糊使用高斯函数生成图像的模糊版本。
//Sigma参数必须为正值，并指示图像的模糊程度。
func Blur(fileImg image.Image, sigma float64) *image.NRGBA {

	if sigma <= 0 {
		sigma = 0.1
	}
	return imaging.Blur(fileImg, sigma)
}

//保存图片
func Save(fileImg image.Image, filename string) error {
	//尺寸调整
	return imaging.Save(fileImg, filename)
}
