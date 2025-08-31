package main

import (
	"encoding/binary"
	"fmt"

	"clipboard_t/pkg/clipboard"
)

// Pixel 表示一个像素点，包含坐标和RGB颜色
type Pixel struct {
	X, Y    int
	R, G, B uint8
}

func create_bmp(pixels []Pixel) []byte {
	// 确定图像尺寸
	var width, height int
	for _, p := range pixels {
		if p.X+1 > width {
			width = p.X + 1
		}
		if p.Y+1 > height {
			height = p.Y + 1
		}
	}

	// 创建像素数据映射
	pixelMap := make([][][]uint8, height)
	for i := range pixelMap {
		pixelMap[i] = make([][]uint8, width)
		for j := range pixelMap[i] {
			pixelMap[i][j] = []uint8{0, 0, 0} // 默认黑色
		}
	}

	// 填充像素数据
	for _, p := range pixels {
		if p.X >= 0 && p.X < width && p.Y >= 0 && p.Y < height {
			pixelMap[p.Y][p.X] = []uint8{p.B, p.G, p.R} // BMP存储顺序是BGR
		}
	}

	// 计算行字节数 (每行必须是4的倍数)
	bytes_per_row := width * 3
	padding := (4 - (bytes_per_row % 4)) % 4
	bytes_per_row += padding

	// 准备文件头 (14字节)
	file_header := make([]byte, 14)
	file_header[0] = 'B'
	file_header[1] = 'M'
	file_size := 14 + 40 + bytes_per_row*height // 14+40 + 像素数据
	binary.LittleEndian.PutUint32(file_header[2:6], uint32(file_size))
	binary.LittleEndian.PutUint32(file_header[10:14], 54) // 像素数据偏移量

	// 准备DIB头 (40字节)
	info_header := make([]byte, 40)
	binary.LittleEndian.PutUint32(info_header[0:4], 40)
	binary.LittleEndian.PutUint32(info_header[4:8], uint32(width))
	binary.LittleEndian.PutUint32(info_header[8:12], uint32(height))
	binary.LittleEndian.PutUint16(info_header[12:14], 1)
	binary.LittleEndian.PutUint16(info_header[14:16], 24)
	binary.LittleEndian.PutUint32(info_header[20:24], uint32(bytes_per_row*height))

	// 创建缓冲区
	buffer := make([]byte, 0, file_size)
	buffer = append(buffer, file_header...)
	buffer = append(buffer, info_header...)

	// 添加像素数据 (从下到上)
	padding_bytes := make([]byte, padding)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			buffer = append(buffer, pixelMap[y][x]...)
		}
		// 添加行填充
		if padding > 0 {
			buffer = append(buffer, padding_bytes...)
		}
	}

	return buffer
}

func main() {
	fmt.Println("正在向剪贴板写入文件列表...")
	err := clipboard.Init()
	if err != nil {
		fmt.Printf("初始化剪贴板失败: %v\n", err)
		return
	}
	// files_name := []string{"./avatar.jpg", "./github-card.png"}
	// files_name := []string{"./sample1.bmp", "./sample1.png", "./sample2.png"}
	// var files_path []string
	// for _, f := range files_name {
	// 	ff := filepath.Join("_example", f)
	// 	file_path, err := filepath.Abs(ff)
	// 	if err != nil {
	// 		continue
	// 	}
	// 	files_path = append(files_path, file_path)
	// }
	// image_file_path := files_path[0]
	// data, err := os.ReadFile(image_file_path)
	// if err != nil {
	// 	fmt.Println("打开文件失败:", err)
	// 	return
	// }
	// fmt.Println("the original buffer", len(data))
	// fmt.Println(data)
	// // png_file, err := png.Decode(bytes.NewReader(data))
	// // if err != nil {
	// // 	fmt.Println("解码 Png 失败 :", err)
	// // 	return
	// // }
	// // var bmp_buf bytes.Buffer

	// // // fmt.Println(png_file.Bounds())

	// img, err := bmp.Decode(bytes.NewReader(data))
	// if err != nil {
	// 	fmt.Println("decode bmp failed, because", err.Error())
	// 	return
	// }
	// fmt.Println("the bmp img", img.Bounds())
	// // // 将image.Image编码为BMP格式写入缓冲区
	// // err = bmp.Encode(&bmp_buf, png_file)
	// // fmt.Println("the bmp bytes", len(bmp_buf.Bytes()))
	// // fmt.Println(bmp_buf.Bytes())
	// // if err != nil {
	// // 	fmt.Println("转换 bmp 失败:", err)
	// // 	return
	// // }

	pixels := []Pixel{
		{X: 0, Y: 0, R: 227, G: 69, B: 22},   // 左上角
		{X: 1, Y: 0, R: 239, G: 239, B: 239}, // 右上角
		{X: 0, Y: 1, R: 227, G: 69, B: 22},   // 左下角
		{X: 1, Y: 1, R: 239, G: 239, B: 239}, // 右下角
	}
	data := create_bmp(pixels)
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
