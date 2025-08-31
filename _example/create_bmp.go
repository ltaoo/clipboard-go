// package main

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"fmt"
// 	"os"
// )

// // BMP 文件头结构
// type BitmapFileHeader struct {
// 	Type     [2]byte
// 	Size     uint32
// 	Reserved [2]uint16
// 	Offset   uint32
// }

// // BMP 信息头结构
// type BitmapInfoHeader struct {
// 	Size            uint32
// 	Width           int32
// 	Height          int32
// 	Planes          uint16
// 	BitCount        uint16
// 	Compression     uint32
// 	SizeImage       uint32
// 	XPelsPerMeter   uint32
// 	YPelsPerMeter   uint32
// 	ColorsUsed      uint32
// 	ColorsImportant uint32
// }

// func main() {
// 	// 初始化文件头
// 	fileHeader := BitmapFileHeader{
// 		Type:     [2]byte{'B', 'M'},
// 		Size:     54 + 2*2*3, // 文件大小 = 头大小（54 字节）+ 像素数据大小
// 		Reserved: [2]uint16{0, 0},
// 		Offset:   54, // 像素数据偏移量
// 	}

// 	// 初始化信息头
// 	infoHeader := BitmapInfoHeader{
// 		Size:            40,
// 		Width:           2,
// 		Height:          2,
// 		Planes:          1,
// 		BitCount:        24,
// 		Compression:     0,
// 		SizeImage:       2 * 2 * 3,
// 		XPelsPerMeter:   0,
// 		YPelsPerMeter:   0,
// 		ColorsUsed:      0,
// 		ColorsImportant: 0,
// 	}

// 	// 像素数据，按照从下到上，从左到右的顺序排列
// 	pixelData := [][]byte{
// 		{227, 69, 22},
// 		{239, 239, 239},
// 		{227, 69, 22},
// 		{239, 239, 239},
// 	}

// 	var pixelBytes bytes.Buffer
// 	for i := len(pixelData) - 1; i >= 0; i-- {
// 		pixelBytes.Write(pixelData[i])
// 	}

// 	// 创建文件
// 	file, err := os.Create("2x2.bmp")
// 	if err != nil {
// 		fmt.Println("创建文件失败:", err)
// 		return
// 	}
// 	defer file.Close()

// 	// 写入文件头
// 	err = binary.Write(file, binary.LittleEndian, &fileHeader)
// 	if err != nil {
// 		fmt.Println("写入文件头失败:", err)
// 		return
// 	}

// 	// 写入信息头
// 	err = binary.Write(file, binary.LittleEndian, &infoHeader)
// 	if err != nil {
// 		fmt.Println("写入信息头失败:", err)
// 		return
// 	}

// 	// 写入像素数据
// 	_, err = file.Write(pixelBytes.Bytes())
// 	if err != nil {
// 		fmt.Println("写入像素数据失败:", err)
// 		return
// 	}

// 	fmt.Println("BMP 图片已创建并保存为 2x2.bmp")

// }

package main

import (
	"encoding/binary"
	"os"
)

// Pixel 表示一个像素点，包含坐标和RGB颜色
type Pixel struct {
	X, Y    int
	R, G, B uint8
}

// CreateBMP 根据像素点列表创建BMP图片
func CreateBMP(filename string, pixels []Pixel) error {
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
	bytesPerRow := width * 3
	padding := (4 - (bytesPerRow % 4)) % 4
	bytesPerRow += padding

	// 准备文件头 (14字节)
	fileHeader := make([]byte, 14)
	fileHeader[0] = 'B'
	fileHeader[1] = 'M'
	fileSize := 54 + bytesPerRow*height // 14+40 + 像素数据
	binary.LittleEndian.PutUint32(fileHeader[2:6], uint32(fileSize))
	binary.LittleEndian.PutUint32(fileHeader[10:14], 54) // 像素数据偏移量

	// 准备DIB头 (40字节)
	dibHeader := make([]byte, 40)
	binary.LittleEndian.PutUint32(dibHeader[0:4], 40)
	binary.LittleEndian.PutUint32(dibHeader[4:8], uint32(width))
	binary.LittleEndian.PutUint32(dibHeader[8:12], uint32(height))
	binary.LittleEndian.PutUint16(dibHeader[12:14], 1)
	binary.LittleEndian.PutUint16(dibHeader[14:16], 24)
	binary.LittleEndian.PutUint32(dibHeader[20:24], uint32(bytesPerRow*height))

	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// 写入文件头
	file.Write(fileHeader)
	// 写入DIB头
	file.Write(dibHeader)

	// 写入像素数据 (从下到上)
	paddingBytes := make([]byte, padding)
	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			file.Write(pixelMap[y][x])
		}
		// 写入行填充
		if padding > 0 {
			file.Write(paddingBytes)
		}
	}

	return nil
}

func main() {
	// 示例：创建2x2的BMP图片
	pixels := []Pixel{
		{X: 0, Y: 0, R: 227, G: 69, B: 22},   // 左上角
		{X: 1, Y: 0, R: 239, G: 239, B: 239}, // 右上角
		{X: 0, Y: 1, R: 227, G: 69, B: 22},   // 左下角
		{X: 1, Y: 1, R: 239, G: 239, B: 239}, // 右下角
	}

	err := CreateBMP("output.bmp", pixels)
	if err != nil {
		panic(err)
	}
}
