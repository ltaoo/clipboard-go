package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clipboard_t/pkg/clipboard"
)

func readClipboard1() {
	fmt.Println("正在读取剪贴板图片...")

	// 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
	}

	// 读取剪贴板中的图片
	imgData := clipboard.Read(clipboard.FmtImage)
	if imgData == nil {
		fmt.Println("剪贴板中没有图片数据")
		os.Exit(1)
	}

	// 解码图片数据
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		fmt.Printf("解码图片失败: %v\n", err)
		os.Exit(1)
	}

	// 生成文件名（使用当前时间戳）
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("clipboard_image_%s.png", timestamp)

	// 创建输出文件
	outputFile, err := os.Create(filename)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		os.Exit(1)
	}
	defer outputFile.Close()

	// 保存图片为PNG格式
	err = png.Encode(outputFile, img)
	if err != nil {
		fmt.Printf("保存图片失败: %v\n", err)
		os.Exit(1)
	}

	// 获取文件的绝对路径
	absPath, err := filepath.Abs(filename)
	if err != nil {
		absPath = filename
	}

	fmt.Printf("图片已成功保存到: %s\n", absPath)
}

func extractFilePaths(input string) []string {
	var paths []string
	// 假设文件路径是以换行符或空格分隔的
	possiblePaths := strings.FieldsFunc(input, func(r rune) bool {
		return r == '\n' || r == ' '
	})
	// 检查字符串是否表示文件路径（可以根据实际情况进行更多检查）
	for _, p := range possiblePaths {
		if strings.Contains(p, ".") {
			paths = append(paths, p)
		}
	}
	return paths
}

func main() {
	// readClipboard1()
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
	}
	files := clipboard.Files()
	fmt.Println("files is", files)
}
