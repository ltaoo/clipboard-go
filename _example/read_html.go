package main

import (
	"fmt"

	"github.com/ltaoo/clipboard-go"
)

func main() {
	fmt.Println("正在读取剪贴板文本...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	html, err := clipboard.ReadHTML()
	if err != nil {
		fmt.Println("读取HTML失败", err.Error())
		return
	}
	text, err := clipboard.ReadText()
	if err != nil {
		fmt.Println("读取文本失败", err.Error())
		return
	}
	fmt.Printf("粘贴板中的文本\n")
	fmt.Println(text)
	fmt.Println("-----")
	fmt.Println(html)
	// r := parser.ParseHTMLContent(text)
	// fmt.Println(r.SourceURL)
	// fmt.Println(r.HTMLContent)
}
