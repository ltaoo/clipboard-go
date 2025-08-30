package main

import (
	"clipboard_t/pkg/clipboard"
	"fmt"
)

func main() {
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	err = clipboard.WriteText("Test content")
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	fmt.Println("写入成功")
	// 默认就能写入成功？
	// if changed != nil {
	// 	<-r
	// 	fmt.Printf("写入成功")
	// }
	// select {
	// case <-changed:
	// 	println(`"text data" is no longer available from clipboard.`)
	// }
}
