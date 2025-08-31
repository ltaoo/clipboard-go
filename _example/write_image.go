package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"clipboard_t/pkg/clipboard"

	"golang.org/x/image/bmp"
)

type BitmapFileHeader struct {
	// bfType     uint16
	// bfSize     uint32
	// bfReserved uint16
	// bfOffBits  uint32
	bfType     [2]byte
	bfSize     uint32
	bfReserved [2]uint16
	bfOffBits  uint32
}

type BitmapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

// Pixel 表示一个像素点，包含坐标和RGB颜色
type Pixel struct {
	X, Y    int
	R, G, B uint8
}

func create_bmp(pixels []Pixel) []byte {
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
	bytes_per_row := width * 4
	// padding := (4 - (bytesPerRow % 4)) % 4
	// bytesPerRow += padding

	// 初始化文件头
	fileHeader := BitmapFileHeader{
		bfType:    [2]byte{'B', 'M'},
		bfSize:    54 + uint32(bytes_per_row*height), // 不需要额外填充
		bfOffBits: 54,
	}

	// 初始化信息头
	infoHeader := BitmapInfoHeader{
		Size:      40,
		Width:     int32(width),
		Height:    int32(height),
		Planes:    1,
		BitCount:  32, // 32位色深
		SizeImage: uint32(bytes_per_row * height),
	}
	// 创建缓冲区
	buffer := make([]byte, 0, fileHeader.bfSize)

	// 写入文件头
	fileHeaderBytes := make([]byte, 14)
	binary.LittleEndian.PutUint16(fileHeaderBytes[0:2], binary.LittleEndian.Uint16(fileHeader.bfType[:]))
	binary.LittleEndian.PutUint32(fileHeaderBytes[2:6], fileHeader.bfSize)
	binary.LittleEndian.PutUint32(fileHeaderBytes[10:14], fileHeader.bfOffBits)
	buffer = append(buffer, fileHeaderBytes...)

	// 写入信息头
	infoHeaderBytes := make([]byte, 40)
	binary.LittleEndian.PutUint32(infoHeaderBytes[0:4], infoHeader.Size)
	binary.LittleEndian.PutUint32(infoHeaderBytes[4:8], uint32(infoHeader.Width))
	binary.LittleEndian.PutUint32(infoHeaderBytes[8:12], uint32(infoHeader.Height))
	binary.LittleEndian.PutUint16(infoHeaderBytes[12:14], infoHeader.Planes)
	binary.LittleEndian.PutUint16(infoHeaderBytes[14:16], infoHeader.BitCount)
	binary.LittleEndian.PutUint32(infoHeaderBytes[20:24], infoHeader.SizeImage)
	binary.LittleEndian.PutUint32(infoHeaderBytes[24:28], uint32(infoHeader.XPelsPerMeter))
	binary.LittleEndian.PutUint32(infoHeaderBytes[28:32], uint32(infoHeader.YPelsPerMeter))
	buffer = append(buffer, infoHeaderBytes...)

	// 写入像素数据 (从下到上)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			pixel := pixelMap[y][x]
			buffer = append(buffer, pixel[0], pixel[1], pixel[2], 0)
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
	files_name := []string{"./sample1.bmp", "./sample1.png", "./sample2.png", "./avatar.jpg", "./github-card.png"}
	var files_path []string
	for _, f := range files_name {
		ff := filepath.Join("_example", f)
		file_path, err := filepath.Abs(ff)
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
	png_file, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		fmt.Println("解码 Png 失败 :", err)
		return
	}
	var bmp_buf bytes.Buffer
	err = bmp.Encode(&bmp_buf, png_file)
	if err != nil {
		fmt.Println("转换 bmp 失败:", err)
		return
	}
	bmp_bytes := bmp_buf.Bytes()
	err = clipboard.WriteImage(bmp_bytes)
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
