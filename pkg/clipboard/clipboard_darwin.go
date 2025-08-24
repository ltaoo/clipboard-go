// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin && !ios

package clipboard

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Cocoa
#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

unsigned int clipboard_read_string(void **out);
unsigned int clipboard_read_image(void **out);
unsigned int clipboard_get_files(void **out);
int clipboard_write_string(const void *bytes, NSInteger n);
int clipboard_write_image(const void *bytes, NSInteger n);
NSInteger clipboard_change_count();
*/
import "C"
import (
	"context"
	"errors"
	"time"
	"unsafe"
)

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	var (
		data unsafe.Pointer
		n    C.uint
	)
	switch t {
	case FmtText:
		n = C.clipboard_read_string(&data)
	case FmtImage:
		n = C.clipboard_read_image(&data)
	}
	if data == nil {
		return nil, errUnavailable
	}
	defer C.free(unsafe.Pointer(data))
	if n == 0 {
		return nil, nil
	}
	return C.GoBytes(data, C.int(n)), nil
}

// write writes the given data to clipboard and
// returns true if success or false if failed.
func write(t Format, buf []byte) (<-chan struct{}, error) {
	var ok C.int
	switch t {
	case FmtText:
		if len(buf) == 0 {
			ok = C.clipboard_write_string(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_string(unsafe.Pointer(&buf[0]),
				C.NSInteger(len(buf)))
		}
	case FmtImage:
		if len(buf) == 0 {
			ok = C.clipboard_write_image(unsafe.Pointer(nil), 0)
		} else {
			ok = C.clipboard_write_image(unsafe.Pointer(&buf[0]),
				C.NSInteger(len(buf)))
		}
	default:
		return nil, errUnsupported
	}
	if ok != 0 {
		return nil, errUnavailable
	}

	// use unbuffered data to prevent goroutine leak
	changed := make(chan struct{}, 1)
	cnt := C.long(C.clipboard_change_count())
	go func() {
		for {
			// not sure if we are too slow or the user too fast :)
			time.Sleep(time.Second)
			cur := C.long(C.clipboard_change_count())
			if cnt != cur {
				changed <- struct{}{}
				close(changed)
				return
			}
		}
	}()
	return changed, nil
}

var (
	errNoFiles       = errors.New("no files in clipboard")
	errInvalidData   = errors.New("invalid data type in clipboard")
	errMemoryAlloc   = errors.New("memory allocation failed")
	errElementType   = errors.New("clipboard contains non-string elements")
	errStringConvert = errors.New("failed to convert string to C format")
	errUnknown       = errors.New("unknown error")
)

func get_files() ([]string, error) {
	// var (
	// 	data unsafe.Pointer
	// 	n    C.uint
	// )
	// n = C.clipboard_get_files(&data)
	// if data == nil {
	// 	return nil, errUnavailable
	// }
	// defer C.free(unsafe.Pointer(data))
	// if n == 0 {
	// 	return nil, nil
	// }
	// return C.GoBytes(data, C.int(n)), nil
	// var cFiles **C.char // 对应 C 的 char**（指向字符串数组的指针）
	var outPtr unsafe.Pointer
	var errCode C.uint
	var (
		data unsafe.Pointer
	)

	// 调用 C 函数获取文件路径数组（错误码保存在 errCode）
	// errCode = C.clipboard_get_files(&data)
	errCode = C.clipboard_get_files(&outPtr)

	// 处理 C 函数返回的错误码
	switch errCode {
	case 0: // 成功
		// 无需处理，继续后续逻辑
	case 1:
		return nil, errInvalidData
	case 2:
		return nil, errNoFiles
	case 3, 6:
		return nil, errMemoryAlloc
	case 4:
		return nil, errElementType
	case 5:
		return nil, errStringConvert
	default:
		return nil, errUnknown
	}
	defer C.free(unsafe.Pointer(data))

	cFiles := (**C.char)(outPtr)
	// 遍历 C 的字符串数组（以 NULL 结尾）
	var files []string
	for i := 0; ; i++ {
		// 通过指针运算获取第 i 个字符串的指针
		cStr := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cFiles)) + uintptr(i)*unsafe.Sizeof(cFiles)))
		if cStr == nil { // 遇到 NULL 结尾标记，停止遍历
			break
		}
		// 将 C 字符串转换为 Go 字符串
		files = append(files, C.GoString(cStr))
	}

	// 手动释放 C 分配的内存（关键！避免内存泄漏）
	// 1. 先释放每个字符串的内存
	for i := 0; ; i++ {
		cStr := *(**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cFiles)) + uintptr(i)*unsafe.Sizeof(cFiles)))
		if cStr == nil {
			break
		}
		C.free(unsafe.Pointer(cStr)) // 释放单个字符串
	}
	// 2. 最后释放整个数组的内存
	C.free(unsafe.Pointer(cFiles))

	return files, nil
}

func watch(ctx context.Context, t Format) <-chan []byte {
	recv := make(chan []byte, 1)
	// not sure if we are too slow or the user too fast :)
	ti := time.NewTicker(time.Second)
	lastCount := C.long(C.clipboard_change_count())
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				this := C.long(C.clipboard_change_count())
				if lastCount != this {
					b := Read(t)
					if b == nil {
						continue
					}
					recv <- b
					lastCount = this
				}
			}
		}
	}()
	return recv
}
