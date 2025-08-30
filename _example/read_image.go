package main

import (
	"fmt"

	"clipboard_t/pkg/clipboard"
	"clipboard_t/pkg/util"
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
		fmt.Println("剪贴板中没有图片数据", err.Error())
		return
	}
	absPath, err := util.SaveByteAsLocalImage(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("粘贴板中的图片已成功保存到本地\n")
	fmt.Printf(absPath)
}
