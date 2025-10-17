# clipboard-go

## Get the content of Clipboard

### Read text

[_example/read_text.go](./_example/read_text.go)

```golang
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
	text, err := clipboard.ReadText()
	if err != nil {
		fmt.Println("读取文本失败", err.Error())
		return
	}
	fmt.Printf("粘贴板中的文本\n")
	fmt.Printf(text)
}
```

### Read html

[_example/read_html.go](./_example/read_html.go)


### Read image

[_example/read_text.go](./_example/read_text.go)

```golang
package main

import (
	"bytes"
	"fmt"
	"image/png"
	"log"

	"github.com/ltaoo/clipboard-go"
)

func main() {
	fmt.Println("正在读取剪贴板图片...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	image_bytes, err := clipboard.ReadImage()
	if err != nil {
		fmt.Println("读取文件失败", err.Error())
		return
	}
	if len(image_bytes) == 0 {
		fmt.Println("剪贴板中没有图片数据")
		return
	}
	reader := bytes.NewReader(image_bytes)
	info, err := png.DecodeConfig(reader)
	if err != nil {
		fmt.Println("failed to decode PNG info")
		return 
	}
	return info.Width, info.Height, nil
}
```

### Read files

[_example/read_file.go](./_example/read_file.go)

```golang
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
```

## Write content to Clipboard

### Write text

[_example/write_text.go](./_example/write_text.go)

### Write html

[_example/write_html.go](./_example/write_html.go)

### Write image

[_example/write_image.go](./_example/write_image.go)

### Write file

[_example/write_file.go](./_example/write_file.go)

## Watch the clipboard

```golang
package main

import (
	"context"
	"fmt"

	"github.com/ltaoo/clipboard-go"
	"github.com/ltaoo/clipboard-go/pkg/util"
)

func main() {
	ch := clipboard.Watch(context.TODO())
	fmt.Println("Start watch the clipboard...")
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
		if data.Type == "public.html" {
			if html, ok := data.Data.(string); ok {
				fmt.Println(html)
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
```

## Acknowledgments

This project was inspired by and references several excellent open-source clipboard libraries. Special thanks to:

- [clipboard-win](https://github.com/DoumanAsh/clipboard-win)
- [clipboard-rs](https://github.com/ChurchTao/clipboard-rs)
- [golang-design/clipboard](https://github.com/golang-design/clipboard)
