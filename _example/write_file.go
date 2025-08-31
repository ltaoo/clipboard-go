package main

import (
	"fmt"
	"path/filepath"

	"clipboard_t/pkg/clipboard"
)

func main() {
	fmt.Println("正在向剪贴板写入文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
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
	err = clipboard.WriteFiles(files_path)
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
