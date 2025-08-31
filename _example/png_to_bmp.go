package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/png"
	"os"
)

// BMPHeader 定义 BMP 文件头
type BMPHeader struct {
	Type      uint16
	Size      uint32
	Reserved1 uint16
	Reserved2 uint16
	OffBits   uint32
}

// BMPInfoHeader 定义 BMP 信息头
type BMPInfoHeader struct {
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

func pngToBmp(pngFilePath, bmpFilePath string) error {
	// 打开 PNG 文件
	pngFile, err := os.Open(pngFilePath)
	if err != nil {
		return err
	}
	defer pngFile.Close()

	// 解码 PNG 图像
	pngImg, err := png.Decode(pngFile)
	if err != nil {
		return err
	}

	// 获取图像边界
	bounds := pngImg.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	// 创建 BMP 文件
	bmpFile, err := os.Create(bmpFilePath)
	if err != nil {
		return err
	}
	defer bmpFile.Close()

	// 计算图像数据大小
	padding := (4 - (width*3)%4) % 4
	imageSize := uint32((width*3 + padding) * height)

	// 写入 BMP 文件头
	bmpHeader := BMPHeader{
		Type:      0x4D42, // 'BM'
		Size:      14 + 40 + imageSize,
		Reserved1: 0,
		Reserved2: 0,
		OffBits:   14 + 40,
	}
	err = binary.Write(bmpFile, binary.LittleEndian, bmpHeader)
	if err != nil {
		return err
	}

	// 写入 BMP 信息头
	bmpInfoHeader := BMPInfoHeader{
		Size:          40,
		Width:         int32(width),
		Height:        int32(height),
		Planes:        1,
		BitCount:      24,
		Compression:   0,
		SizeImage:     imageSize,
		XPelsPerMeter: 0,
		YPelsPerMeter: 0,
		ClrUsed:       0,
		ClrImportant:  0,
	}
	err = binary.Write(bmpFile, binary.LittleEndian, bmpInfoHeader)
	if err != nil {
		return err
	}

	// 写入图像数据
	for y := height - 1; y >= 0; y-- {
		lineBuffer := bytes.NewBuffer(nil)
		for x := 0; x < width; x++ {
			r, g, b, _ := pngImg.At(x, y).RGBA()
			buffer := []byte{byte(b >> 8), byte(g >> 8), byte(r >> 8)}
			lineBuffer.Write(buffer)
		}
		// 填充到 4 字节对齐
		for i := 0; i < padding; i++ {
			lineBuffer.WriteByte(0)
		}
		_, err = bmpFile.Write(lineBuffer.Bytes())
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	err := pngToBmp("./_example/sample1.png", "output.bmp")
	if err != nil {
		fmt.Println("转换失败:", err)
	} else {
		fmt.Println("转换成功")
	}
}
