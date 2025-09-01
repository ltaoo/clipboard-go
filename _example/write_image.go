package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ltaoo/clipboard-go"
)

func main() {
	fmt.Println("正在向剪贴板写入文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	files_name := []string{"./_example/sample1.bmp", "./_example/sample1.png", "./_example/sample2.png", "./_example/avatar.jpg", "./_example/github-card.png"}
	var files_path []string
	for _, f := range files_name {
		file_path, err := filepath.Abs(f)
		if err != nil {
			continue
		}
		files_path = append(files_path, file_path)
	}
	image_file_path := files_path[4]
	data, err := os.ReadFile(image_file_path)
	if err != nil {
		fmt.Println("打开文件失败:", err)
		return
	}
	err = clipboard.WriteImage(data)
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
