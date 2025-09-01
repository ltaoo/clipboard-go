package main

import (
	"encoding/binary"
	"fmt"
	"os"
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

func CreateBMPBuffer(pixels []Pixel) []byte {
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
	// 示例：创建2x2的BMP图片
	pixels := []Pixel{
		{X: 0, Y: 0, R: 227, G: 69, B: 22},   // 左上角
		{X: 1, Y: 0, R: 239, G: 239, B: 239}, // 右上角
		{X: 0, Y: 1, R: 227, G: 69, B: 22},   // 左下角
		{X: 1, Y: 1, R: 239, G: 239, B: 239}, // 右下角
	}

	buffer := CreateBMPBuffer(pixels)

	file, err := os.Create("2x2.bmp")
	if err != nil {
		fmt.Println("create file failed,", err.Error())
		return
	}

	file.Write(buffer)
}
