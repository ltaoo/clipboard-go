package main

import (
	"fmt"

	"github.com/ltaoo/clipboard-go"
)

func main() {
	fmt.Println("正在读取剪贴板文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	files, err := clipboard.ReadFiles()
	if err != nil {
		fmt.Println("读取文件失败", err.Error())
		return
	}
	if len(files) == 0 {
		fmt.Println("剪贴板中没有文件数据")
		return
	}
	fmt.Printf("粘贴板中的文件列表\n")
	for _, f := range files {
		fmt.Println(string(f))
	}
}
