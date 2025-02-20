package article_service

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"

	"github.com/3Eeeecho/go-gin-example/pkg/file"
	"github.com/3Eeeecho/go-gin-example/pkg/qrcode"
	"github.com/fogleman/gg"
)

type ArticlePoster struct {
	PosterName string
	*Article
	Qr *qrcode.QrCode
}

func NewArticlePoster(postName string, article *Article, qr *qrcode.QrCode) *ArticlePoster {
	return &ArticlePoster{
		PosterName: postName,
		Article:    article,
		Qr:         qr,
	}
}

func GetPosterFlag() string {
	return "poster"
}

func (a *ArticlePoster) CheckMergedImage(path string) bool {
	return !file.CheckFileNotExist(path + a.PosterName)
}

func (a *ArticlePoster) OpenMergedImage(path string) (*os.File, error) {
	f, err := file.MustOpen(a.PosterName, path)
	if err != nil {
		return nil, err
	}

	return f, nil
}

type ArticlePosterBg struct {
	Name string
	*ArticlePoster
	*Rect
	*Pt
}

type Rect struct {
	Name string
	X0   int
	Y0   int
	X1   int
	Y1   int
}

type Pt struct {
	X int
	Y int
}

func NewArticlePosterBg(name string, ap *ArticlePoster, rect *Rect, pt *Pt) *ArticlePosterBg {
	return &ArticlePosterBg{
		Name:          name,
		ArticlePoster: ap,
		Rect:          rect,
		Pt:            pt,
	}
}

func (a *ArticlePosterBg) Generate() (string, string, error) {
	fullPath := qrcode.GetQrCodeFullPath()
	fileName, path, err := a.Qr.Encode(fullPath)
	if err != nil {
		return "", "", err
	}

	if !a.CheckMergedImage(path) {
		mergedF, err := a.OpenMergedImage(path)
		if err != nil {
			return "", "", err
		}
		defer mergedF.Close()

		bgF, err := file.MustOpen(a.Name, path)
		if err != nil {
			return "", "", err
		}
		defer bgF.Close()

		qrF, err := file.MustOpen(fileName, path)
		if err != nil {
			return "", "", err
		}
		defer qrF.Close()

		bgImage, err := jpeg.Decode(bgF)
		if err != nil {
			return "", "", err
		}
		qrImage, err := jpeg.Decode(qrF)
		if err != nil {
			return "", "", err
		}

		// 获取背景图片的宽度和高度
		imgWidth := bgImage.Bounds().Dx()
		imgHeight := bgImage.Bounds().Dy()

		// 计算标题的 X 坐标，确保水平居中
		X0 := float64(imgWidth) / 2
		// 将标题放置在图片的四分之一高度处
		Y0 := float64(imgHeight) / 4
		// 计算副标题的 Y 坐标，将其放置在标题下方 80px
		Y1 := Y0 + 80

		jpg := image.NewRGBA(image.Rect(a.Rect.X0, a.Rect.Y0, a.Rect.X1, a.Rect.Y1))
		draw.Draw(jpg, jpg.Bounds(), bgImage, bgImage.Bounds().Min, draw.Over)
		draw.Draw(jpg, jpg.Bounds(), qrImage, qrImage.Bounds().Min.Sub(image.Pt(a.Pt.X, a.Pt.Y)), draw.Over)

		err = a.DrawPoster(&DrawText{
			JPG:    jpg,
			Merged: mergedF,

			Title: "Golang Gin GitHub",
			X0:    int(X0),
			Y0:    int(Y0),
			Size0: 42,

			SubTitle: "---Eecho",
			X1:       int(X0),
			Y1:       int(Y1),
			Size1:    36,
		}, "msyhbd.ttc")

		if err != nil {
			return "", "", err
		}
	}

	return fileName, path, nil
}

type DrawText struct {
	JPG      image.Image // 要绘制文本的图片
	Merged   *os.File    // 合并后的目标文件
	Title    string      // 标题文本
	X0, Y0   int         // 标题文本位置
	Size0    float64     // 标题字体大小
	SubTitle string      // 副标题文本
	X1, Y1   int         // 副标题文本位置
	Size1    float64     // 副标题字体大小
}

func (a *ArticlePosterBg) DrawPoster(dt *DrawText, fontPath string) error {
	dc := gg.NewContextForRGBA(dt.JPG.(*image.RGBA))
	dc.SetRGB(0, 0, 0)

	err := dc.LoadFontFace(fontPath, dt.Size0)
	if err != nil {
		return fmt.Errorf("加载标题字体失败: %v", err)
	}
	dc.DrawStringAnchored(dt.Title, float64(dt.X0), float64(dt.Y0), 0.5, 0.5)

	err = dc.LoadFontFace(fontPath, dt.Size1)
	if err != nil {
		return fmt.Errorf("加载副标题字体失败: %v", err)
	}
	dc.DrawStringAnchored(dt.SubTitle, float64(dt.X1), float64(dt.Y1), 0.5, 0.5)

	if err = jpeg.Encode(dt.Merged, dc.Image(), nil); err != nil {
		return fmt.Errorf("保存图片失败: %v", err)
	}

	return nil
}
