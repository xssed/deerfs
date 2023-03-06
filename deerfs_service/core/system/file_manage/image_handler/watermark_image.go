package image_handler

import (
	"bytes"
	"errors"
	"image"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"strings"

	//"github.com/h2non/filetype"
	"github.com/xssed/deerfs/deerfs_service/core/system/file_manage"
)

// 水印的位置
const (
	TopLeftImg     Pos = iota //0 顶部左边
	TopRightImg               //1 顶部右边
	BottomLeftImg             //2 底部左边
	BottomRightImg            //3 底部右边
	CenterImg                 //4 中间
)

var (
	// ErrUnsupportedWatermarkType 不支持的水印类型
	ErrUnsupportedWatermarkType = errors.New("Unsupported watermark type.")

	// ErrWatermarkTooLarge 当水印位置距离右下角的范围小于水印图片时，返回错误。
	ErrWatermarkTooLarge = errors.New("Watermark size is too large.")

	//ErrInvalidPos  pos值无效
	ErrInvalidPos = errors.New("Invalid 'pos' value")
)

// Pos 表示水印的位置
type Pos int

// WatermarkImage 用于给图片添加水印功能
//
// 目前支持 gif、jpeg 和 png 三种图片格式。
// 若是 gif 图片，则只取图片的第一帧；png 支持透明背景。
type WatermarkImage struct {
	image   image.Image // 水印图片
	gifImg  *gif.GIF    // 如果是 GIF 图片，image 保存第一帧的图片， gifImg 保存全部内容
	padding int         // 水印留的边白
	pos     Pos         // 水印的位置
}

// NewFromFile 从文件声明一个 WatermarkImage 对象
//
// path 为水印文件的路径；
// padding 为水印在目标图像上的留白大小；
// pos 水印的位置的代表数字。
func NewFromFile(path string, padding int, pos int) (*WatermarkImage, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file_type, _ := file_manage.CheckFileType(path)

	//根据水印数值来设置Pos
	if pos > 4 || pos < 0 {
		//防止偏移错误,默认底部右侧
		return New(f, file_type, padding, BottomRightImg)
	} else {
		if pos == 0 {
			//顶部左边
			return New(f, file_type, padding, TopLeftImg)
		} else if pos == 1 {
			//顶部右边
			return New(f, file_type, padding, TopRightImg)
		} else if pos == 2 {
			//底部左边
			return New(f, file_type, padding, BottomLeftImg)
		} else if pos == 4 {
			//中间
			return New(f, file_type, padding, CenterImg)
		} else {
			//默认底部右侧
			return New(f, file_type, padding, BottomRightImg)
		}

	}

}

// New 声明 WatermarkImage 对象
//
// r 为水印图片内容；
// ext 为水印图片的扩展名，会根据扩展名判断图片类型；
// padding 为水印在目标图像上的留白大小；
// pos 图片位置；
func New(r io.Reader, ext string, padding int, pos Pos) (w *WatermarkImage, err error) {

	if pos < TopLeftImg || pos > CenterImg {
		return nil, ErrInvalidPos
	}

	var img image.Image
	var gifImg *gif.GIF
	switch ext {
	case "jpg":
		img, err = jpeg.Decode(r)
	case "png":
		img, err = png.Decode(r)
	case "gif":
		gifImg, err = gif.DecodeAll(r)
		img = gifImg.Image[0]
	default:
		return nil, ErrUnsupportedWatermarkType
	}
	if err != nil {
		return nil, err
	}

	return &WatermarkImage{
		image:   img,
		gifImg:  gifImg,
		padding: padding,
		pos:     pos,
	}, nil
}

// NewForImage 声明 WatermarkImage 对象
//
// wmImg 为水印图片内容；
// padding 为水印在目标图像上的留白大小；
// pos 图片位置；
func NewForImage(wmImg image.Image, padding int, pos Pos) (w *WatermarkImage, err error) {

	if pos < TopLeftImg || pos > CenterImg {
		return nil, ErrInvalidPos
	}

	var img image.Image
	var gifImg *gif.GIF

	img = wmImg

	return &WatermarkImage{
		image:   img,
		gifImg:  gifImg,
		padding: padding,
		pos:     pos,
	}, nil
}

// NewForGif 声明 WatermarkImage 对象
//
// r 为水印图片内容；
// padding 为水印在目标图像上的留白大小；
// pos 图片位置；
func NewForGif(r io.Reader, padding int, pos Pos) (w *WatermarkImage, err error) {

	if pos < TopLeftImg || pos > CenterImg {
		return nil, ErrInvalidPos
	}

	var img image.Image
	var gifImg *gif.GIF
	gifImg, err = gif.DecodeAll(r)
	img = gifImg.Image[0]

	if err != nil {
		return nil, err
	}

	return &WatermarkImage{
		image:   img,
		gifImg:  gifImg,
		padding: padding,
		pos:     pos,
	}, nil

}

// MarkFile 给指定的文件打上水印
func (w *WatermarkImage) MarkFile(path string) ([]byte, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	file_type, _ := file_manage.CheckFileType(path)

	return w.Mark(f, file_type)
}

func (w *WatermarkImage) Save(b []byte, dstFile string) error {

	newFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer newFile.Close()
	_, err = newFile.Write(b) //将字符串转换成字节切片
	if err != nil {
		return err
	}
	return nil

}

// Mark 将水印写入 src 中，由 ext 确定当前图片的类型。
func (w *WatermarkImage) Mark(src io.ReadWriteSeeker, ext string) ([]byte, error) {

	var srcImg image.Image
	var err error

	ext = strings.ToLower(ext)
	switch ext {
	case "gif":
		return w.markGIF(src) // GIF 另外单独处理
	case "jpg":
		srcImg, err = jpeg.Decode(src)
	case "png":
		srcImg, err = png.Decode(src)
	default:
		return nil, ErrUnsupportedWatermarkType
	}
	if err != nil {
		return nil, err
	}

	bound := srcImg.Bounds()
	point, pos_err := w.getPoint(bound.Dx(), bound.Dy())
	if pos_err != nil {
		return nil, pos_err
	}

	if err = w.checkTooLarge(point, bound); err != nil {
		return nil, err
	}

	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)

	if _, err = src.Seek(0, 0); err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}

	switch ext {
	case "jpg":
		err := jpeg.Encode(&buf, dstImg, nil)
		return buf.Bytes(), err
	case "png":
		err := png.Encode(&buf, dstImg)
		return buf.Bytes(), err
	default:
		return nil, ErrUnsupportedWatermarkType
	}
}

// MarkForImage 将水印写入 image.Image 中
func (w *WatermarkImage) MarkForImage(bgImg image.Image, ext string) ([]byte, error) {

	var srcImg image.Image
	var err error

	ext = strings.ToLower(ext)
	if ext == "jpg" || ext == "png" {
		srcImg = bgImg
	} else {
		return nil, ErrUnsupportedWatermarkType
	}

	bound := srcImg.Bounds()
	point, pos_err := w.getPoint(bound.Dx(), bound.Dy())
	if pos_err != nil {
		return nil, pos_err
	}

	if err = w.checkTooLarge(point, bound); err != nil {
		return nil, err
	}

	dstImg := image.NewNRGBA64(srcImg.Bounds())
	draw.Draw(dstImg, dstImg.Bounds(), srcImg, image.Point{}, draw.Src)
	draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)

	buf := bytes.Buffer{}

	switch ext {
	case "jpg":
		err := jpeg.Encode(&buf, dstImg, nil)
		return buf.Bytes(), err
	case "png":
		err := png.Encode(&buf, dstImg)
		return buf.Bytes(), err
	default:
		return nil, ErrUnsupportedWatermarkType
	}

}

func (w *WatermarkImage) markGIF(src io.ReadWriteSeeker) ([]byte, error) {

	srcGIF, err := gif.DecodeAll(src)
	if err != nil {
		return nil, err
	}
	bound := srcGIF.Image[0].Bounds()
	point, pos_err := w.getPoint(bound.Dx(), bound.Dy())
	if pos_err != nil {
		return nil, pos_err
	}

	if err = w.checkTooLarge(point, bound); err != nil {
		return nil, err
	}

	if w.gifImg == nil {
		for index, img := range srcGIF.Image {
			dstImg := image.NewPaletted(img.Bounds(), img.Palette)
			draw.Draw(dstImg, dstImg.Bounds(), img, image.Point{}, draw.Src)
			draw.Draw(dstImg, dstImg.Bounds(), w.image, point, draw.Over)
			srcGIF.Image[index] = dstImg
		}
	} else { // 水印也是 GIF
		windex := 0
		wmax := len(w.gifImg.Image)
		for index, img := range srcGIF.Image {
			dstImg := image.NewPaletted(img.Bounds(), img.Palette)
			draw.Draw(dstImg, dstImg.Bounds(), img, image.Point{}, draw.Src)

			// 获取对应帧数的水印图片
			if windex >= wmax {
				windex = 0
			}
			draw.Draw(dstImg, dstImg.Bounds(), w.gifImg.Image[windex], point, draw.Over)
			windex++

			srcGIF.Image[index] = dstImg
		}
	}

	if _, err = src.Seek(0, 0); err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	err = gif.EncodeAll(&buf, srcGIF)
	return buf.Bytes(), err
}

func (w *WatermarkImage) checkTooLarge(start image.Point, dst image.Rectangle) error {
	// 允许的最大高宽
	width := dst.Dx() - start.X - w.padding
	height := dst.Dy() - start.Y - w.padding

	if width < w.image.Bounds().Dx() || height < w.image.Bounds().Dy() {
		return ErrWatermarkTooLarge
	}
	return nil
}

func (w *WatermarkImage) getPoint(width, height int) (image.Point, error) {

	var point image.Point

	switch w.pos {
	case TopLeftImg:
		point = image.Point{X: -w.padding, Y: -w.padding}
	case TopRightImg:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -w.padding,
		}
	case BottomLeftImg:
		point = image.Point{
			X: -w.padding,
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case BottomRightImg:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()),
			Y: -(height - w.padding - w.image.Bounds().Dy()),
		}
	case CenterImg:
		point = image.Point{
			X: -(width - w.padding - w.image.Bounds().Dx()) / 2,
			Y: -(height - w.padding - w.image.Bounds().Dy()) / 2,
		}
	default:
		return image.Point{
			X: 0,
			Y: 0,
		}, ErrInvalidPos
	}

	return point, nil

}
