package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clipboard_t/pkg/clipboard"
)

func readTextFromClipboard() {
	fmt.Println("正在读取剪贴板文本...")

	// 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
	}

	// 读取剪贴板中的图片
	text := clipboard.Read(clipboard.FmtText)
	if text == nil {
		fmt.Println("剪贴板中没有文本数据")
		os.Exit(1)
	}

	fmt.Printf("文本是: %s\n", text)
}
func writeTextToClipboard() {
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
	}
	changed := clipboard.Write(clipboard.FmtText, []byte("Test content"))
	// 默认就能写入成功？
	// if changed != nil {
	// 	<-r
	// 	fmt.Printf("写入成功")
	// }
	select {
	case <-changed:
		println(`"text data" is no longer available from clipboard.`)
	}
}

func readImageFromClipboard() {
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

func readFilepathsFromClipboard() {
	fmt.Println("正在读取剪贴板文件列表...")

	// 初始化剪贴板
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
	}

	// 读取剪贴板中的图片
	files := clipboard.Read(clipboard.FmtFilepath)
	// if err != nil {
	// 	fmt.Println("读取失败", err.Error())
	// 	os.Exit(1)
	// }
	if files == nil {
		fmt.Println("剪贴板中没有文件数据")
		os.Exit(1)
	}

	fmt.Printf("粘贴板中的文件列表 %v\n", string(files))
	// for _, f := range files {
	// 	fmt.Println(string(f))
	// }
}
func writeFilesToClipboard(files []string) {
	fmt.Println("正在向剪贴板写入文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		os.Exit(1)
		return
	}

	v, err := StringSliceToByteSlice(files)
	if err != nil {
		fmt.Printf(" %v\n", err)
		return
	}
	clipboard.Write(clipboard.FmtFilepath, v)
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

func readClipboard1() {

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

func producer(c chan int) {
	for i := 0; i < 5; i++ {
		c <- i
		time.Sleep(time.Millisecond * 100)
	}
	close(c)
}

func consumer(c chan int) {
	for num := range c {
		fmt.Println("消费:", num)
	}
}
func testChan() {
	ch := make(chan int)
	go producer(ch)
	go consumer(ch)
	for {
		select {
		case num, ok := <-ch:
			if !ok {
				fmt.Println("通道已关闭，退出循环")
				return
			}
			fmt.Println("从通道获取到数据:", num)
		}
	}
}

// StringSliceToByteSlice 将 []string 转换为 []byte
func StringSliceToByteSlice(strs []string) ([]byte, error) {
	return json.Marshal(strs)
}

// ByteSliceToStringSlice 将 []byte 转换回 []string
func ByteSliceToStringSlice(b []byte) ([]string, error) {
	var strs []string
	err := json.Unmarshal(b, &strs)
	return strs, err
}

func main() {
	// arr := objc.ID(class_NSMutableArray).Send(sel_alloc).Send(sel_init)

	// files := []string{"/Users/mayfair/Documents/deploy_step4.png", "/Users/mayfair/Documents/StatsCard.tsx"}
	// files := []string{"/Users/litao/Downloads/avatar.png", "/Users/litao/Downloads/face.png"}
	files := []string{"/Users/litao/Downloads/flutterio-icon.svg"}
	writeFilesToClipboard(files)
	// readFilepathsFromClipboard()
	// bytes, err := StringSliceToByteSlice(files)
	// if err != nil {
	// 	return
	// }

	// for _, f := range files {
	// 	ss2 := (*int8)(unsafe.Pointer(&[]byte(f + "\x00")[0]))
	// 	v2 := objc.ID(class_NSString).Send(sel_stringWithUTF8String, ss2)
	// 	url2 := objc.ID(class_NSURL).Send(sel_fileURLWithPath, v2)
	// 	arr.Send(sel_addObject, url2)
	// }

	// the_file_count := uint(arr.Send(sel_count))
	// fmt.Println("the_file_count", the_file_count)

	// pasteboard := objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	// pasteboard.Send(sel_clearContents)

	// pasteboard.Send(sel_writeObjects, arr)
	// pasteboard.Send(sel_propertyListForType, _NSPasteboardTypeFiles)

}
