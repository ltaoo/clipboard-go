package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

// 将图像转换为RGBA颜色模型
func toRGBA(img image.Image) *image.RGBA {
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba
}

func jpgToPng(jpgFilePath, pngFilePath string) error {
	// 打开JPEG文件
	jpgFile, err := os.Open(jpgFilePath)
	if err != nil {
		return err
	}
	defer jpgFile.Close()

	// 解码JPEG文件
	img, err := jpeg.Decode(jpgFile)
	if err != nil {
		return err
	}

	// 将图像转换为RGBA颜色模型
	rgbaImg := toRGBA(img)

	// 创建PNG文件
	pngFile, err := os.Create(pngFilePath)
	if err != nil {
		return err
	}
	defer pngFile.Close()

	// 将图像编码为PNG格式并写入文件
	err = png.Encode(pngFile, rgbaImg)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	files_name := []string{"./avatar.jpg", "./github-card.png"}
	var files_path []string
	for _, f := range files_name {
		ff := filepath.Join("_example", f)
		file_path, err := filepath.Abs(ff)
		if err != nil {
			continue
		}
		files_path = append(files_path, file_path)
	}
	image_file_path := files_path[0]

	jpgFilePath := image_file_path
	pngFilePath := "output.png"

	err := jpgToPng(jpgFilePath, pngFilePath)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		fmt.Println("转换成功")
	}
}
