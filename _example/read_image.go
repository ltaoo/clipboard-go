package main

import (
	"bytes"
	"fmt"
	"image/png"

	"github.com/ltaoo/clipboard-go"
	"github.com/ltaoo/clipboard-go/pkg/util"
)

func main() {
	fmt.Println("正在读取剪贴板图片...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	data, err := clipboard.ReadImage()
	if err != nil {
		fmt.Println("读取图片失败", err.Error())
		return
	}
	absPath, err := util.SaveByteAsLocalImage(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("粘贴板中的图片已成功保存到本地\n")
	fmt.Println(absPath)
	reader := bytes.NewReader(data)
	info, err := png.DecodeConfig(reader)
	if err != nil {
		fmt.Println("failed to decode PNG info")
		return
	}
	fmt.Printf("the image width is %v, height is %v \n", info.Width, info.Height)
}
