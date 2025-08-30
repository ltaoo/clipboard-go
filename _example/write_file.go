package main

import (
	"clipboard_t/pkg/clipboard"
	"fmt"
)

func main() {
	fmt.Println("正在向剪贴板写入文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	// files := []string{"/Users/mayfair/Documents/deploy_step4.png", "/Users/mayfair/Documents/StatsCard.tsx"}
	files := []string{"/Users/litao/Downloads/avatar.png", "/Users/litao/Downloads/face.png"}
	// files := []string{"/Users/litao/Downloads/flutterio-icon.svg"}
	err = clipboard.WriteFiles(files)
	if err != nil {
		fmt.Printf(" %v\n", err)
		return
	}
	fmt.Println("写入成功")
	// if changed != nil {
	// 	<-r
	// 	fmt.Printf("写入成功")
	// }
	// select {
	// case <-changed:
	// 	println(`"text data" is no longer available from clipboard.`)
	// }
}
