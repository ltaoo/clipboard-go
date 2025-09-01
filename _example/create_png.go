package main

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

func main() {
	// 创建一个2x2的图像
	img := image.NewRGBA(image.Rect(0, 0, 2, 2))

	// 设置像素颜色
	// 左上角(0,0)
	img.Set(0, 0, color.RGBA{R: 227, G: 69, B: 22, A: 255})
	// 右上角(1,0)
	img.Set(1, 0, color.RGBA{R: 239, G: 239, B: 239, A: 255})
	// 左下角(0,1)
	img.Set(0, 1, color.RGBA{R: 227, G: 69, B: 22, A: 255})
	// 右下角(1,1)
	img.Set(1, 1, color.RGBA{R: 239, G: 239, B: 239, A: 255})

	// 创建输出文件
	file, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// 编码为PNG格式并写入文件
	err = png.Encode(file, img)
	if err != nil {
		panic(err)
	}
}
