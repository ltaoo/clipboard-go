package main

import (
	"context"
	"fmt"

	"clipboard_t/pkg/clipboard"
	"clipboard_t/pkg/util"
)

func main() {
	ch := clipboard.Watch(context.TODO())
	fmt.Println("开始监听粘贴板")
	for data := range ch {
		fmt.Println(data.Type)
		if data.Type == "public.file-url" {
			if files, ok := data.Data.([]string); ok {
				for _, f := range files {
					fmt.Println(f)
				}
			}
		}
		if data.Type == "public.utf8-plain-text" {
			if text, ok := data.Data.(string); ok {
				fmt.Println(text)
			}
		}
		if data.Type == "public.png" {
			if f, ok := data.Data.([]byte); ok {
				img_filepath, err := util.SaveByteAsLocalImage(f)
				if err == nil {
					fmt.Println("the image save to", img_filepath)
				}
			}
		}
	}
}
