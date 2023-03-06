package image_handler

import (
	"bytes"
	_ "embed"
	"errors"

	//"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/unknwon/com"
)

//默认字体，方正宋体
//go:embed 1fangzhengsong.ttf
var fontBytes_song []byte

//文泉驿正黑
//go:embed 2wqy-zenhei.ttc
var fontBytes_zhenghei []byte

//方正楷体
//go:embed 3fangzhengkaiti.ttf
var fontBytes_fzkai []byte

//濑户字体
//go:embed 4setofont.ttf
var fontBytes_laihu []byte

//Lingxun英文
//go:embed 5LingxunSerif.ttf
var fontBytes_lingxun []byte

//Roboto
//go:embed 6Roboto.ttf
var fontBytes_roboto []byte

//RobotoSerif
//go:embed 7RobotoSerif.ttf
var fontBytes_robotoserif []byte

const (
	TopLeft int = iota
	TopRight
	BottomLeft
	BottomRight
	Center
)

//定义水印类
type WaterMarkText struct {
	Quality int
}

func (w *WaterMarkText) SetParam(quality int) {

	//设置图片质量,默认100
	if quality < 1 {
		w.Quality = 100
	} else {
		w.Quality = quality
	}

}

func (w *WaterMarkText) NewTextImageWaterMark(staticImg image.Image, fileType string, fontInfo []FontInfo) ([]byte, error) {

	var err error

	img := image.NewNRGBA(staticImg.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, staticImg.At(x, y))
		}
	}
	err = w.do(img, fontInfo)
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	switch fileType {
	case "png":
		err = png.Encode(&buf, img)
	default:
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: w.Quality})
	}
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil //&image.RGBA{Pix: buf.Bytes()}

}

func (w *WaterMarkText) NewTextGifWaterMark(srcFile string, fontInfo []FontInfo) ([]byte, error) {

	var err error
	imgFile, err := os.Open(srcFile)
	if err != nil {
		return nil, err
	}
	defer imgFile.Close()

	gifImg, err := gif.DecodeAll(imgFile)
	if err != nil {
		return nil, err
	}
	gifs := make([]*image.Paletted, 0)
	x0 := 0
	y0 := 0
	yuan := 0
	for k, v := range gifImg.Image {
		img := image.NewNRGBA(v.Bounds())
		if k == 0 {
			x0 = img.Bounds().Dx()
			y0 = img.Bounds().Dy()
		}
		if k == 0 && gifImg.Image[k+1].Bounds().Dx() > x0 && gifImg.Image[k+1].Bounds().Dy() > y0 {
			yuan = 1
			break
		}
		if x0 == img.Bounds().Dx() && y0 == img.Bounds().Dy() {
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					img.Set(x, y, v.At(x, y))
				}
			}
			err = w.do(img, fontInfo)
			if err != nil {
				break
			}
			p1 := image.NewPaletted(v.Bounds(), v.Palette)
			draw.Draw(p1, v.Bounds(), img, image.Point{}, draw.Src)
			gifs = append(gifs, p1)
		} else {
			gifs = append(gifs, v)
		}
	}
	if yuan == 1 {
		return nil, errors.New("gif: image block is out of bounds")
	}
	if err != nil {
		return nil, err
	}

	buf := bytes.Buffer{}
	g1 := &gif.GIF{
		Image:     gifs,
		Delay:     gifImg.Delay,
		LoopCount: gifImg.LoopCount,
	}
	err = gif.EncodeAll(&buf, g1)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil

}

func (w *WaterMarkText) do(img draw.Image, fontInfo []FontInfo) error {

	var err error
	errNum := 1
Loop:
	for _, v := range fontInfo {

		//设置字体,默认方正宋体
		v.selectFont(v.FontIndex)

		info := v.Message
		f := freetype.NewContext()
		f.SetDPI(v.DPI)
		f.SetFont(v.FontBytes)
		f.SetFontSize(v.Size)
		f.SetClip(img.Bounds())
		f.SetDst(img)
		f.SetSrc(image.NewUniform(color.RGBA{R: v.R, G: v.G, B: v.B, A: v.A}))
		first := 0
		two := 0
		switch v.Position {
		case TopLeft:
			first = v.Dx
			two = v.Dy + int(f.PointToFixed(v.Size)>>6)
		case TopRight:
			first = img.Bounds().Dx() - len(info)*4 - v.Dx
			two = v.Dy + int(f.PointToFixed(v.Size)>>6)
		case BottomLeft:
			first = v.Dx
			two = img.Bounds().Dy() - v.Dy
		case BottomRight:
			first = img.Bounds().Dx() - len(info)*4 - v.Dx
			two = img.Bounds().Dy() - v.Dy
		case Center:
			first = (img.Bounds().Dx() - len(info)*4) / 2
			two = (img.Bounds().Dy() - v.Dy) / 2
		default:
			errNum = 0
			break Loop
		}
		pt := freetype.Pt(first, two)
		_, err = f.DrawString(info, pt)
		if err != nil {
			break
		}
	}
	if errNum == 0 {
		err = errors.New("The value of the Position is wrong.")
	}
	return err
}

func (w *WaterMarkText) Save(b []byte, dstFile string) error {

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

type FontInfo struct {
	DPI       float64
	FontIndex int
	FontBytes *truetype.Font
	Size      float64 // 文字大小
	Message   string  // 文字内容
	Position  int     // 文字存放位置
	Dx        int     // 文字x轴留白距离
	Dy        int     // 文字y轴留白距离
	R         uint8   // 文字颜色值RGBA中的R值
	G         uint8   // 文字颜色值RGBA中的G值
	B         uint8   // 文字颜色值RGBA中的B值
	A         uint8   // 文字颜色值RGBA中的A值
}

//设置字体,默认方正宋体
func (fi *FontInfo) selectFont(fontIndex int) {

	fi.FontIndex = fontIndex
	switch fi.FontIndex {
	case 1:
		font, _ := freetype.ParseFont(fontBytes_song)
		fi.FontBytes = font
	case 2:
		font, _ := freetype.ParseFont(fontBytes_zhenghei)
		fi.FontBytes = font
	case 3:
		font, _ := freetype.ParseFont(fontBytes_fzkai)
		fi.FontBytes = font
	case 4:
		font, _ := freetype.ParseFont(fontBytes_laihu)
		fi.FontBytes = font
	case 5:
		font, _ := freetype.ParseFont(fontBytes_lingxun)
		fi.FontBytes = font
	case 6:
		font, _ := freetype.ParseFont(fontBytes_roboto)
		fi.FontBytes = font
	case 7:
		font, _ := freetype.ParseFont(fontBytes_robotoserif)
		fi.FontBytes = font
	default:
		font, _ := freetype.ParseFont(fontBytes_song)
		fi.FontBytes = font
	}
}

//设置文字颜色,默认白色
func (fi *FontInfo) selectColor(rgba_input string) {

	rgba_str_arr := strings.Split(rgba_input, "_")
	//fmt.Println(rgba_str_arr)//调试
	if len(rgba_str_arr) != 4 {
		fi.R = 255
		fi.G = 255
		fi.B = 255
		fi.A = 255
		return
	}
	if com.StrTo(rgba_str_arr[0]).MustUint8() > 255 {
		fi.R = 255
	} else {
		fi.R = com.StrTo(rgba_str_arr[0]).MustUint8()
	}
	if com.StrTo(rgba_str_arr[1]).MustUint8() > 255 {
		fi.G = 255
	} else {
		fi.G = com.StrTo(rgba_str_arr[1]).MustUint8()
	}
	if com.StrTo(rgba_str_arr[2]).MustUint8() > 255 {
		fi.B = 255
	} else {
		fi.B = com.StrTo(rgba_str_arr[2]).MustUint8()
	}
	if com.StrTo(rgba_str_arr[3]).MustUint8() > 255 {
		fi.A = 255
	} else {
		fi.A = com.StrTo(rgba_str_arr[3]).MustUint8()
	}

}

func (fi *FontInfo) SetParam(fontId int, rgba_input string, fontSize, dpi float64, x, y, position int, message string) {

	//设置字体,默认方正宋体
	fi.selectFont(fontId)
	//设置字体颜色
	fi.selectColor(rgba_input)
	//设置字体大小,默认17
	if fontSize < 1 {
		fi.Size = 17
	} else {
		fi.Size = fontSize
	}
	//设置DPI,默认75
	if dpi < 1 {
		fi.DPI = 75
	} else {
		fi.DPI = dpi
	}
	//设置X,Y偏移
	if x < 1 {
		fi.Dx = 0
	} else {
		fi.Dx = x
	}
	if y < 1 {
		fi.Dy = 0
	} else {
		fi.Dy = y
	}
	//设置Position
	fi.Position = position
	//设置Message
	fi.Message = message

}
