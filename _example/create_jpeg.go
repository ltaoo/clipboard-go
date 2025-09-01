package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	// 创建一个2x2的图像
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	// 设置像素颜色
	// 注意：在image.Image坐标系中，(0,0)是左上角
	redColor := color.RGBA{227, 69, 22, 255}    // RGB(227, 69, 22)
	grayColor := color.RGBA{239, 239, 239, 255} // RGB(239, 239, 239)

	// 左上角 (0,0)
	img.Set(0, 0, redColor)
	// 右上角 (1,0)
	img.Set(1, 0, grayColor)
	// 左下角 (0,1)
	img.Set(0, 1, redColor)
	// 右下角 (1,1)
	img.Set(1, 1, grayColor)

	// 创建输出文件
	file, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 编码为JPEG格式
	err = jpeg.Encode(file, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
	}
}
