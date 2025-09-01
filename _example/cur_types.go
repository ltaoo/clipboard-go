package main

import (
	"fmt"

	"github.com/ltaoo/clipboard-go"
)

func main() {
	fmt.Println("正在读取剪贴板内容类型...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	types := clipboard.GetContentTypes(clipboard.ContentTypeParams{IsEnabled: false})
	if len(types) == 0 {
		fmt.Println("剪贴板中没有数据")
		return
	}
	fmt.Printf("粘贴板中的内容类型\n")
	for _, t := range types {
		fmt.Println(string(t))
	}
}
