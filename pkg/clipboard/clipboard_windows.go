// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build windows

package clipboard

// Interacting with Clipboard on Windows:
// https://docs.microsoft.com/zh-cn/windows/win32/dataxchg/using-the-clipboard

import (
	"bytes"
	"context"
	"unsafe"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"reflect"
	"runtime"
	"strings"
	"syscall"
	"time"
	"unicode/utf16"

	"golang.org/x/image/bmp"
)

func initialize() error { return nil }

// readText reads the clipboard and returns the text data if presents.
// The caller is responsible for opening/closing the clipboard before
// calling this function.
func readText() (buf []byte, err error) {
	hMem, _, err := getClipboardData.Call(cFmtUnicodeText)
	if hMem == 0 {
		return nil, err
	}
	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return nil, err
	}
	defer gUnlock.Call(hMem)

	// Find NUL terminator
	n := 0
	for ptr := unsafe.Pointer(p); *(*uint16)(ptr) != 0; n++ {
		ptr = unsafe.Pointer(uintptr(ptr) + unsafe.Sizeof(*((*uint16)(unsafe.Pointer(p)))))
	}

	var s []uint16
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	h.Data = p
	h.Len = n
	h.Cap = n
	return []byte(string(utf16.Decode(s))), nil
}

// writeText writes given data to the clipboard. It is the caller's
// responsibility for opening/closing the clipboard before calling
// this function.
func writeText(buf []byte) error {
	r, _, err := emptyClipboard.Call()
	if r == 0 {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// empty text, we are done here.
	if len(buf) == 0 {
		return nil
	}

	s, err := syscall.UTF16FromString(string(buf))
	if err != nil {
		return fmt.Errorf("failed to convert given string: %w", err)
	}

	hMem, _, err := gAlloc.Call(gmemMoveable, uintptr(len(s)*int(unsafe.Sizeof(s[0]))))
	if hMem == 0 {
		return fmt.Errorf("failed to alloc global memory: %w", err)
	}

	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return fmt.Errorf("failed to lock global memory: %w", err)
	}
	defer gUnlock.Call(hMem)

	// no return value
	memMove.Call(p, uintptr(unsafe.Pointer(&s[0])),
		uintptr(len(s)*int(unsafe.Sizeof(s[0]))))

	v, _, err := setClipboardData.Call(cFmtUnicodeText, hMem)
	if v == 0 {
		gFree.Call(hMem)
		return fmt.Errorf("failed to set text to clipboard: %w", err)
	}

	return nil
}

// readImage reads the clipboard and returns PNG encoded image data
// if presents. The caller is responsible for opening/closing the
// clipboard before calling this function.
func readImage() ([]byte, error) {
	hMem, _, err := getClipboardData.Call(cFmtDIBV5)
	if hMem == 0 {
		// second chance to try FmtDIB
		return readImageDib()
	}
	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return nil, err
	}
	defer gUnlock.Call(hMem)

	// inspect header information
	info := (*bitmapV5Header)(unsafe.Pointer(p))

	// maybe deal with other formats?
	if info.BitCount != 32 {
		return nil, errUnsupported
	}

	var data []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = uintptr(p)
	sh.Cap = int(info.Size + 4*uint32(info.Width)*uint32(info.Height))
	sh.Len = int(info.Size + 4*uint32(info.Width)*uint32(info.Height))
	img := image.NewRGBA(image.Rect(0, 0, int(info.Width), int(info.Height)))
	offset := int(info.Size)
	stride := int(info.Width)
	for y := 0; y < int(info.Height); y++ {
		for x := 0; x < int(info.Width); x++ {
			idx := offset + 4*(y*stride+x)
			xhat := (x + int(info.Width)) % int(info.Width)
			yhat := int(info.Height) - 1 - y
			r := data[idx+2]
			g := data[idx+1]
			b := data[idx+0]
			a := data[idx+3]
			img.SetRGBA(xhat, yhat, color.RGBA{r, g, b, a})
		}
	}
	// always use PNG encoding.
	var buf bytes.Buffer
	png.Encode(&buf, img)
	return buf.Bytes(), nil
}

func readImageDib() ([]byte, error) {
	const (
		fileHeaderLen = 14
		infoHeaderLen = 40
		cFmtDIB       = 8
	)

	hClipDat, _, err := getClipboardData.Call(cFmtDIB)
	if err != nil {
		return nil, errors.New("not dib format data: " + err.Error())
	}
	pMemBlk, _, err := gLock.Call(hClipDat)
	if pMemBlk == 0 {
		return nil, errors.New("failed to call global lock: " + err.Error())
	}
	defer gUnlock.Call(hClipDat)

	bmpHeader := (*bitmapHeader)(unsafe.Pointer(pMemBlk))
	dataSize := bmpHeader.SizeImage + fileHeaderLen + infoHeaderLen

	if bmpHeader.SizeImage == 0 && bmpHeader.Compression == 0 {
		iSizeImage := bmpHeader.Height * ((bmpHeader.Width*uint32(bmpHeader.BitCount)/8 + 3) &^ 3)
		dataSize += iSizeImage
	}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, uint16('B')|(uint16('M')<<8))
	binary.Write(buf, binary.LittleEndian, uint32(dataSize))
	binary.Write(buf, binary.LittleEndian, uint32(0))
	const sizeof_colorbar = 0
	binary.Write(buf, binary.LittleEndian, uint32(fileHeaderLen+infoHeaderLen+sizeof_colorbar))
	j := 0
	for i := fileHeaderLen; i < int(dataSize); i++ {
		binary.Write(buf, binary.BigEndian, *(*byte)(unsafe.Pointer(pMemBlk + uintptr(j))))
		j++
	}
	return bmpToPng(buf)
}

func bmpToPng(bmpBuf *bytes.Buffer) (buf []byte, err error) {
	var f bytes.Buffer
	original_image, err := bmp.Decode(bmpBuf)
	if err != nil {
		return nil, err
	}
	err = png.Encode(&f, original_image)
	if err != nil {
		return nil, err
	}
	return f.Bytes(), nil
}

func writeImage(buf []byte) error {
	r, _, err := emptyClipboard.Call()
	if r == 0 {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// empty text, we are done here.
	if len(buf) == 0 {
		return nil
	}

	img, err := png.Decode(bytes.NewReader(buf))
	if err != nil {
		return fmt.Errorf("input bytes is not PNG encoded: %w", err)
	}

	offset := unsafe.Sizeof(bitmapV5Header{})
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	imageSize := 4 * width * height

	data := make([]byte, int(offset)+imageSize)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := int(offset) + 4*(y*width+x)
			r, g, b, a := img.At(x, height-1-y).RGBA()
			data[idx+2] = uint8(r)
			data[idx+1] = uint8(g)
			data[idx+0] = uint8(b)
			data[idx+3] = uint8(a)
		}
	}

	info := bitmapV5Header{}
	info.Size = uint32(offset)
	info.Width = int32(width)
	info.Height = int32(height)
	info.Planes = 1
	info.Compression = 0 // BI_RGB
	info.SizeImage = uint32(4 * info.Width * info.Height)
	info.RedMask = 0xff0000 // default mask
	info.GreenMask = 0xff00
	info.BlueMask = 0xff
	info.AlphaMask = 0xff000000
	info.BitCount = 32 // we only deal with 32 bpp at the moment.
	// Use calibrated RGB values as Go's image/png assumes linear color space.
	// Other options:
	// - LCS_CALIBRATED_RGB = 0x00000000
	// - LCS_sRGB = 0x73524742
	// - LCS_WINDOWS_COLOR_SPACE = 0x57696E20
	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-wmf/eb4bbd50-b3ce-4917-895c-be31f214797f
	info.CSType = 0x73524742
	// Use GL_IMAGES for GamutMappingIntent
	// Other options:
	// - LCS_GM_ABS_COLORIMETRIC = 0x00000008
	// - LCS_GM_BUSINESS = 0x00000001
	// - LCS_GM_GRAPHICS = 0x00000002
	// - LCS_GM_IMAGES = 0x00000004
	// https://docs.microsoft.com/en-us/openspecs/windows_protocols/ms-wmf/9fec0834-607d-427d-abd5-ab240fb0db38
	info.Intent = 4 // LCS_GM_IMAGES

	infob := make([]byte, int(unsafe.Sizeof(info)))
	for i, v := range *(*[unsafe.Sizeof(info)]byte)(unsafe.Pointer(&info)) {
		infob[i] = v
	}
	copy(data[:], infob[:])

	hMem, _, err := gAlloc.Call(gmemMoveable,
		uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	if hMem == 0 {
		return fmt.Errorf("failed to alloc global memory: %w", err)
	}

	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return fmt.Errorf("failed to lock global memory: %w", err)
	}
	defer gUnlock.Call(hMem)

	memMove.Call(p, uintptr(unsafe.Pointer(&data[0])),
		uintptr(len(data)*int(unsafe.Sizeof(data[0]))))

	v, _, err := setClipboardData.Call(cFmtDIBV5, hMem)
	if v == 0 {
		gFree.Call(hMem)
		return fmt.Errorf("failed to set text to clipboard: %w", err)
	}

	return nil
}

func byte_slice_to_string_slice(b []byte) ([]string, error) {
	var strs []string
	err := json.Unmarshal(b, &strs)
	return strs, err
}

type HDROPHeader struct {
	pFiles uint32
	x      int16
	y      int16
	fNC    uint32
	fWide  uint32
}

func writeFiles(buf []byte) error {
	ret, _, err := emptyClipboard.Call()
	if ret == 0 {
		return fmt.Errorf("failed to clear clipboard: %w", err)
	}

	// empty text, we are done here.
	if len(buf) == 0 {
		return nil
	}
	filePaths, err := byte_slice_to_string_slice(buf)
	if err != nil {
		return err
	}
	if len(filePaths) == 0 {
		return nil
	}

	// 计算文件路径的总长度
	var fileListSize uint32
	for _, path := range filePaths {
		var count uint32
		ret, _, err := multiByteToWideChar.Call(
			CP_UTF8,
			0,
			uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
			uintptr(int32(len(path))),
			0,
			0,
		)
		if ret == 0 {
			return fmt.Errorf("MultiByteToWideChar (to get length) for path %s failed: %w", path, err)
		}
		count = uint32(ret)
		fileListSize += count + 1
	}

	if fileListSize == 0 {
		return fmt.Errorf("No valid file paths")
	}

	// 计算总内存大小
	dropfiles := DROPFILES{
		p_files: uint32(unsafe.Sizeof(DROPFILES{})),
		pt:      POINT{x: 0, y: 0},
		f_nc:    0,
		f_wide:  1,
	}
	memSize := uintptr(unsafe.Sizeof(dropfiles)) + uintptr(fileListSize*2) + 2

	// 分配全局内存
	hMem, _, err := gAlloc.Call(0x0042, memSize)
	if hMem == 0 {
		return fmt.Errorf("GlobalAlloc failed: %w", err)
	}
	defer gFree.Call(hMem)

	// 锁定内存
	p, _, err := gLock.Call(hMem)
	if p == 0 {
		return fmt.Errorf("GlobalLock failed: %w", err)
	}
	defer gUnlock.Call(hMem)

	// 填充DROPFILES结构体
	ptr := (*DROPFILES)(unsafe.Pointer(p))
	*ptr = dropfiles

	// 填充文件路径
	dataPtr := (*[1 << 30]uint16)(unsafe.Pointer(uintptr(p) + unsafe.Sizeof(dropfiles)))
	for _, path := range filePaths {
		var count uint32
		ret, _, err := multiByteToWideChar.Call(
			CP_UTF8,
			0,
			uintptr(unsafe.Pointer(syscall.StringBytePtr(path))),
			uintptr(int32(len(path))),
			uintptr(unsafe.Pointer(&dataPtr[0])),
			uintptr(int32(fileListSize)),
		)
		if ret == 0 {
			return fmt.Errorf("MultiByteToWideChar (to write path) for path %s failed: %w", path, err)
		}
		count = uint32(ret)
		dataPtr = (*[1 << 30]uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(&dataPtr[count])) + 2))
		// *dataPtr = 0
	}
	// 添加最终的null终止符
	// *(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + 2)) = 0
	*(*uint16)(unsafe.Pointer(uintptr(unsafe.Pointer(dataPtr)) + 2)) = uint16(0)

	// 设置剪贴板数据
	ret, _, err = setClipboardData.Call(CF_HDROP, hMem)
	if ret == 0 {
		return fmt.Errorf("SetClipboardData failed: %w", err)
	}
	return nil

	// // 计算总的路径长度和HDROPHeader的大小
	// totalLength := uint32(unsafe.Sizeof(HDROPHeader{}))
	// for _, path := range filePaths {
	// 	fmt.Println("the file prepare write to clipboard", path)
	// 	totalLength += uint32(len(path)+1) * 2 // 每个字符2字节（UTF - 16），加上null终止符
	// }

	// // 分配全局内存
	// hMem, _, err := gAlloc.Call(0x0042, uintptr(totalLength))
	// if hMem == 0 {
	// 	fmt.Println("e1", err.Error())
	// 	return fmt.Errorf("GlobalAlloc failed: %w", err)
	// }
	// defer gFree.Call(hMem)

	// // 锁定内存
	// p, _, err := gLock.Call(hMem)
	// if p == 0 {
	// 	fmt.Println("e2", err.Error())
	// 	return fmt.Errorf("GlobalLock failed: %w", err)
	// }
	// defer gUnlock.Call(hMem)

	// // 填充HDROPHeader
	// hdr := (*HDROPHeader)(unsafe.Pointer(p))
	// hdr.pFiles = uint32(len(filePaths))
	// hdr.x = 0
	// hdr.y = 0
	// hdr.fNC = 0
	// hdr.fWide = 1

	// offset := uintptr(unsafe.Sizeof(HDROPHeader{}))
	// for _, path := range filePaths {
	// 	utf16Path, err := syscall.UTF16FromString(path)
	// 	if err != nil {
	// 		fmt.Println("e3", err.Error())
	// 		return fmt.Errorf("UTF16FromString failed: %w", err)
	// 	}
	// 	// 手动复制UTF - 16字符串到内存
	// 	dst := (*[1 << 30]uint16)(unsafe.Pointer(p + offset))
	// 	for i, v := range utf16Path {
	// 		dst[i] = v
	// 	}
	// 	// 添加null终止符
	// 	dst[len(utf16Path)] = 0
	// 	offset += uintptr((len(utf16Path)+1) * 2)
	// }

	// // 设置剪贴板数据
	// r, _, err = setClipboardData.Call(cFmtFilepaths, hMem)
	// fmt.Println("after setClipboardData", r, err.Error())
	// if r == 0 {
	// 	return fmt.Errorf("SetClipboardData failed: %w", err)
	// }
	// return nil

	// 填充文件路径
	// offset := uintptr(unsafe.Sizeof(HDROPHeader{}))
	// for _, path := range filePaths {
	// 	utf16Path, err := syscall.UTF16FromString(path)
	// 	if err != nil {
	// 		return fmt.Errorf("UTF16FromString failed: %w", err)
	// 	}
	// 	_, _, err = lstrcpyW.Call(p+offset, uintptr(unsafe.Pointer(&utf16Path[0])))
	// 	if err != nil {
	// 		return fmt.Errorf("lstrcpyW failed: %w", err)
	// 	}
	// 	offset += uintptr((len(utf16Path) + 1) * 2)
	// }

	// // 设置剪贴板数据
	// ret, _, err := setClipboardData.Call(cFmtFilepaths, hMem)
	// if ret == 0 {
	// 	return fmt.Errorf("SetClipboardData failed: %w", err)
	// }

	// return nil

	// infob := make([]byte, int(unsafe.Sizeof(info)))
	// for i, v := range *(*[unsafe.Sizeof(info)]byte)(unsafe.Pointer(&info)) {
	// 	infob[i] = v
	// }
	// copy(data[:], infob[:])

	// hMem, _, err := gAlloc.Call(gmemMoveable,
	// 	uintptr(len(data)*int(unsafe.Sizeof(data[0]))))
	// if hMem == 0 {
	// 	return fmt.Errorf("failed to alloc global memory: %w", err)
	// }

	// p, _, err := gLock.Call(hMem)
	// if p == 0 {
	// 	return fmt.Errorf("failed to lock global memory: %w", err)
	// }
	// defer gUnlock.Call(hMem)

	// memMove.Call(p, uintptr(unsafe.Pointer(&data[0])),
	// 	uintptr(len(data)*int(unsafe.Sizeof(data[0]))))

	// v, _, err := setClipboardData.Call(cFmtFilepaths, hMem)
	// if v == 0 {
	// 	gFree.Call(hMem)
	// 	return fmt.Errorf("failed to set files to clipboard: %w", err)
	// }

	// return nil
}

// https://stackoverflow.com/questions/77205618/when-a-file-is-on-the-windows-clipboard-how-can-i-in-python-access-its-path

func readFilepaths() ([]byte, error) {
	hMem, _, err := getClipboardData.Call(cFmtFilepaths)
	if hMem == 0 {
		fmt.Println("f1", err.Error())
		return nil, err
	}
	p, _, err := gLock.Call(hMem)
	if p == 0 {
		fmt.Println("f2", err.Error())
		return nil, err
	}
	defer gUnlock.Call(hMem)

	// 验证HDROP结构的内存布局
	type HDROPHeader struct {
		pFiles uint32
		x      int16
		y      int16
		fNC    uint32
		fWide  uint32
	}

	var count uint32
	ret, v, err := dragQueryFile.Call(p, uintptr(^uint32(0)), 0, 0, uintptr(unsafe.Sizeof(count)), uintptr(unsafe.Pointer(&count)))
	if ret == 0 {
		fmt.Println("f3", err.Error())
		return nil, fmt.Errorf("DragQueryFile (to get count) failed: %w", err)
	}

	fmt.Println("num files", ret, v, err)
	fileCount := uint32(ret)

	// // 存储文件路径
	filePaths := make([]string, fileCount)
	for i := uint32(0); i < fileCount; i++ {
		// 获取文件路径所需长度（不包含 null 终止符）
		var length uint32
		ret, _, err = dragQueryFile.Call(p, uintptr(i), 0, 0, uintptr(unsafe.Sizeof(length)), uintptr(unsafe.Pointer(&length)))
		if ret == 0 {
			return nil, fmt.Errorf("DragQueryFile (to get length) for file %d failed: %w", i, err)
		}
		length = uint32(ret)

		buffer := make([]uint16, length+1)
		ret, _, err = dragQueryFile.Call(p, uintptr(i), uintptr(unsafe.Pointer(&buffer[0])), uintptr(len(buffer)*2))
		if ret == 0 {
			return nil, fmt.Errorf("DragQueryFile (to get path) for file %d failed: %w", i, err)
		}

		filePaths = append(filePaths, syscall.UTF16ToString(buffer[:length]))
	}
	joinedPaths := strings.Join(filePaths, "\n")
	return []byte(joinedPaths), nil
}

func read(t Format) (buf []byte, err error) {
	// On Windows, OpenClipboard and CloseClipboard must be executed on
	// the same thread. Thus, lock the OS thread for further execution.
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	var format uintptr
	switch t {
	case FmtImage:
		format = cFmtDIBV5
	case FmtFilepath:
		format = cFmtFilepaths
	case FmtText:
		fallthrough
	default:
		format = cFmtUnicodeText
	}

	// check if clipboard is avaliable for the requested format
	r, _, err := isClipboardFormatAvailable.Call(format)
	fmt.Println("after check clipboard ", r, format)
	if r == 0 {
		return nil, errUnavailable
	}

	// try again until open clipboard successed
	for {
		r, _, _ = openClipboard.Call()
		if r == 0 {
			continue
		}
		break
	}
	defer closeClipboard.Call()

	switch format {
	case cFmtDIBV5:
		return readImage()
	case cFmtFilepaths:
		return readFilepaths()
	case cFmtUnicodeText:
		fallthrough
	default:
		return readText()
	}
}

// write writes the given data to clipboard and
// returns true if success or false if failed.
func write(t Format, buf []byte) (<-chan struct{}, error) {
	errch := make(chan error)
	changed := make(chan struct{}, 1)
	go func() {
		// make sure GetClipboardSequenceNumber happens with
		// OpenClipboard on the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
		for {
			r, _, _ := openClipboard.Call(0)
			if r == 0 {
				continue
			}
			break
		}

		// var param uintptr
		switch t {
		case FmtImage:
			err := writeImage(buf)
			if err != nil {
				errch <- err
				closeClipboard.Call()
				return
			}
		case FmtFilepath:
			err := writeFiles(buf)
			if err != nil {
				errch <- err
				closeClipboard.Call()
				return
			}
		case FmtText:
			fallthrough
		default:
			// param = cFmtUnicodeText
			err := writeText(buf)
			if err != nil {
				errch <- err
				closeClipboard.Call()
				return
			}
		}
		// Close the clipboard otherwise other applications cannot
		// paste the data.
		closeClipboard.Call()

		cnt, _, _ := getClipboardSequenceNumber.Call()
		errch <- nil
		for {
			time.Sleep(time.Second)
			cur, _, _ := getClipboardSequenceNumber.Call()
			if cur != cnt {
				changed <- struct{}{}
				close(changed)
				return
			}
		}
	}()
	err := <-errch
	if err != nil {
		return nil, err
	}
	return changed, nil
}

func watch(ctx context.Context, t Format) <-chan []byte {
	recv := make(chan []byte, 1)
	ready := make(chan struct{})
	go func() {
		// not sure if we are too slow or the user too fast :)
		ti := time.NewTicker(time.Second)
		cnt, _, _ := getClipboardSequenceNumber.Call()
		ready <- struct{}{}
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				cur, _, _ := getClipboardSequenceNumber.Call()
				if cnt != cur {
					b := Read(t)
					if b == nil {
						continue
					}
					recv <- b
					cnt = cur
				}
			}
		}
	}()
	<-ready
	return recv
}

const (
	cFmtBitmap      = 2 // Win+PrintScreen
	cFmtUnicodeText = 13
	cFmtFilepaths   = 15
	cFmtDIBV5       = 17
	// Screenshot taken from special shortcut is in different format (why??), see:
	// https://jpsoft.com/forums/threads/detecting-clipboard-format.5225/
	cFmtDataObject = 49161 // Shift+Win+s, returned from enumClipboardFormats
	gmemMoveable   = 0x0002

	CP_UTF8    = 65001
	CF_HDROP   = 15
	WM_DROPFILES = 0x0233
)

// BITMAPV5Header structure, see:
// https://docs.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-bitmapv5header
type bitmapV5Header struct {
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
	RedMask       uint32
	GreenMask     uint32
	BlueMask      uint32
	AlphaMask     uint32
	CSType        uint32
	Endpoints     struct {
		CiexyzRed, CiexyzGreen, CiexyzBlue struct {
			CiexyzX, CiexyzY, CiexyzZ int32 // FXPT2DOT30
		}
	}
	GammaRed    uint32
	GammaGreen  uint32
	GammaBlue   uint32
	Intent      uint32
	ProfileData uint32
	ProfileSize uint32
	Reserved    uint32
}

type bitmapHeader struct {
	Size          uint32
	Width         uint32
	Height        uint32
	PLanes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter uint32
	YPelsPerMeter uint32
	ClrUsed       uint32
	ClrImportant  uint32
}


// 定义POINT结构体
type POINT struct {
    x int32
    y int32
}

// 定义DROPFILES结构体
type DROPFILES struct {
    p_files uint32
    pt      POINT
    f_nc    int32
    f_wide  int32
}


// Calling a Windows DLL, see:
// https://github.com/golang/go/wiki/WindowsDLLs
var (
	user32 = syscall.MustLoadDLL("user32")
	// Opens the clipboard for examination and prevents other
	// applications from modifying the clipboard content.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-openclipboard
	openClipboard = user32.MustFindProc("OpenClipboard")
	// Closes the clipboard.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-closeclipboard
	closeClipboard = user32.MustFindProc("CloseClipboard")
	// Empties the clipboard and frees handles to data in the clipboard.
	// The function then assigns ownership of the clipboard to the
	// window that currently has the clipboard open.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-emptyclipboard
	emptyClipboard = user32.MustFindProc("EmptyClipboard")
	// Retrieves data from the clipboard in a specified format.
	// The clipboard must have been opened previously.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getclipboarddata
	getClipboardData = user32.MustFindProc("GetClipboardData")
	// Places data on the clipboard in a specified clipboard format.
	// The window must be the current clipboard owner, and the
	// application must have called the OpenClipboard function. (When
	// responding to the WM_RENDERFORMAT message, the clipboard owner
	// must not call OpenClipboard before calling SetClipboardData.)
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-setclipboarddata
	setClipboardData = user32.MustFindProc("SetClipboardData")
	// Determines whether the clipboard contains data in the specified format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-isclipboardformatavailable
	isClipboardFormatAvailable = user32.MustFindProc("IsClipboardFormatAvailable")
	// Clipboard data formats are stored in an ordered list. To perform
	// an enumeration of clipboard data formats, you make a series of
	// calls to the EnumClipboardFormats function. For each call, the
	// format parameter specifies an available clipboard format, and the
	// function returns the next available clipboard format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-isclipboardformatavailable
	enumClipboardFormats = user32.MustFindProc("EnumClipboardFormats")
	// Retrieves the clipboard sequence number for the current window station.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getclipboardsequencenumber
	getClipboardSequenceNumber = user32.MustFindProc("GetClipboardSequenceNumber")
	// Registers a new clipboard format. This format can then be used as
	// a valid clipboard format.
	// https://docs.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-registerclipboardformata
	registerClipboardFormatA = user32.MustFindProc("RegisterClipboardFormatA")
	// lstrcpyW                 = user32.MustFindProc("lstrcpyW")

	shell32       = syscall.NewLazyDLL("shell32")
	dragQueryFile = shell32.NewProc("DragQueryFileW")

	kernel32 = syscall.NewLazyDLL("kernel32")

	// Locks a global memory object and returns a pointer to the first
	// byte of the object's memory block.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globallock
	gLock = kernel32.NewProc("GlobalLock")
	gSize = kernel32.NewProc("GlobalSize")
	multiByteToWideChar = kernel32.NewProc("MultiByteToWideChar")
	// Decrements the lock count associated with a memory object that was
	// allocated with GMEM_MOVEABLE. This function has no effect on memory
	// objects allocated with GMEM_FIXED.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalunlock
	gUnlock = kernel32.NewProc("GlobalUnlock")
	// Allocates the specified number of bytes from the heap.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalalloc
	gAlloc = kernel32.NewProc("GlobalAlloc")
	// Frees the specified global memory object and invalidates its handle.
	// https://docs.microsoft.com/en-us/windows/win32/api/winbase/nf-winbase-globalfree
	gFree   = kernel32.NewProc("GlobalFree")
	memMove = kernel32.NewProc("RtlMoveMemory")
)
