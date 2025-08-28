//go:build darwin && !ios

package clipboard

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
)

var (
	appkit = must(purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_GLOBAL|purego.RTLD_NOW))

	_NSPasteboardTypeString = must2(purego.Dlsym(appkit, "NSPasteboardTypeString"))
	_NSPasteboardTypePNG    = must2(purego.Dlsym(appkit, "NSPasteboardTypePNG"))
	_NSPasteboardTypeFiles  = must2(purego.Dlsym(appkit, "NSFilenamesPboardType"))

	class_NSPasteboard = objc.GetClass("NSPasteboard")
	class_NSData       = objc.GetClass("NSData")
	class_NSArray      = objc.GetClass("NSArray")
	class_NSString     = objc.GetClass("NSString")
	class_NSURL        = objc.GetClass("NSURL")

	sel_generalPasteboard        = objc.RegisterName("generalPasteboard")
	sel_length                   = objc.RegisterName("length")
	sel_getBytesLength           = objc.RegisterName("getBytes:length:")
	sel_dataForType              = objc.RegisterName("dataForType:")
	sel_propertyListForType      = objc.RegisterName("propertyListForType:")
	sel_setPropertyList_forType_ = objc.RegisterName("setPropertyListForType:")
	sel_clearContents            = objc.RegisterName("clearContents")
	sel_setDataForType           = objc.RegisterName("setData:forType:")
	sel_dataWithBytesLength      = objc.RegisterName("dataWithBytes:length:")
	sel_changeCount              = objc.RegisterName("changeCount")
	sel_count                    = objc.RegisterName("count")
	sel_UTF8String               = objc.RegisterName("UTF8String")
	sel_objectAtIndex            = objc.RegisterName("objectAtIndex:")
	sel_stringWithUTF8String     = objc.RegisterName("stringWithUTF8String:")
	sel_arrayWithObjects_count   = objc.RegisterName("arrayWithObjects:count:")
)

func must(sym uintptr, err error) uintptr {
	if err != nil {
		panic(err)
	}
	return sym
}

func must2(sym uintptr, err error) uintptr {
	if err != nil {
		panic(err)
	}
	// dlsym returns a pointer to the object so dereference like this to avoid possible misuse of 'unsafe.Pointer' warning
	return **(**uintptr)(unsafe.Pointer(&sym))
}

func initialize() error { return nil }

func read(t Format) (buf []byte, err error) {
	switch t {
	case FmtText:
		return clipboard_read_string(), nil
	case FmtImage:
		return clipboard_read_image(), nil
	case FmtFilepath:
		return clipboard_read_files(), nil
	}
	return nil, errUnavailable
}

func write(t Format, buf []byte) (<-chan struct{}, error) {
	var ok bool
	switch t {
	case FmtText:
		if len(buf) == 0 {
			ok = clipboard_write_string(nil)
		} else {
			ok = clipboard_write_string(buf)

		}
	case FmtImage:
		if len(buf) == 0 {
			ok = clipboard_write_image(nil)
		} else {
			ok = clipboard_write_image(buf)
		}
	case FmtFilepath:
		if len(buf) == 0 {
			ok = clipboard_write_files(nil)
		} else {
			ok = clipboard_write_files(buf)
		}
	default:
		return nil, errUnsupported
	}
	if !ok {
		return nil, errUnavailable
	}
	changed := make(chan struct{}, 1)
	cnt := clipboard_change_count()
	go func() {
		for {
			// not sure if we are too slow or the user too fast :)
			time.Sleep(time.Second)
			cur := clipboard_change_count()
			if cnt != cur {
				changed <- struct{}{}
				close(changed)
				return
			}
		}
	}()
	return changed, nil
}

func watch(ctx context.Context, t Format) <-chan []byte {
	recv := make(chan []byte, 1)
	ti := time.NewTicker(time.Second)
	lastCount := clipboard_change_count()
	go func() {
		for {
			select {
			case <-ctx.Done():
				close(recv)
				return
			case <-ti.C:
				this := clipboard_change_count()
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

func clipboard_read_string() []byte {
	var pasteboard = objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	var data = pasteboard.Send(sel_dataForType, _NSPasteboardTypeString)
	if data == 0 {
		return nil
	}
	var size = uint(data.Send(sel_length))
	if size == 0 {
		return nil
	}
	out := make([]byte, size)
	data.Send(sel_getBytesLength, unsafe.SliceData(out), size)
	return out
}

func clipboard_read_image() []byte {
	var pasteboard = objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	data := pasteboard.Send(sel_dataForType, _NSPasteboardTypePNG)
	if data == 0 {
		return nil
	}
	size := data.Send(sel_length)
	out := make([]byte, size)
	data.Send(sel_getBytesLength, unsafe.SliceData(out), size)
	return out
}

func readUTF8String(ptr unsafe.Pointer) string {
	if ptr == nil {
		return ""
	}
	var length int
	for ; *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(length))) != 0; length++ {
	}
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = *(*byte)(unsafe.Pointer(uintptr(ptr) + uintptr(i)))
	}
	return string(bytes)
}

func clipboard_read_files() []byte {
	var pasteboard = objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	data := pasteboard.Send(sel_propertyListForType, _NSPasteboardTypeFiles)
	if data == 0 {
		return nil
	}
	array := objc.ID(data)
	countResult := array.Send(sel_count)
	count := int(countResult)
	var strs []string
	for i := 0; i < count; i++ {
		fileObj := array.Send(sel_objectAtIndex, int(i))
		utf8Ptr := unsafe.Pointer(fileObj.Send(sel_UTF8String))
		if utf8Ptr == nil {
			continue
		}
		fileStr := readUTF8String(utf8Ptr)
		strs = append(strs, fileStr)
	}

	// 将字符串切片转换为字节切片
	var totalLen int
	for _, str := range strs {
		totalLen += len(str)
	}
	result := make([]byte, totalLen)
	offset := 0
	for _, str := range strs {
		strBytes := []byte(str)
		copy(result[offset:], strBytes)
		offset += len(strBytes)
	}
	return result
}

func clipboard_write_image(bytes []byte) bool {
	pasteboard := objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	data := objc.ID(class_NSData).Send(sel_dataWithBytesLength, unsafe.SliceData(bytes), len(bytes))
	pasteboard.Send(sel_clearContents)
	return pasteboard.Send(sel_setDataForType, data, _NSPasteboardTypePNG) != 0
}

func clipboard_write_string(bytes []byte) bool {
	pasteboard := objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	data := objc.ID(class_NSData).Send(sel_dataWithBytesLength, unsafe.SliceData(bytes), len(bytes))
	pasteboard.Send(sel_clearContents)
	return pasteboard.Send(sel_setDataForType, data, _NSPasteboardTypeString) != 0
}
func byte_slice_to_string_slice(b []byte) ([]string, error) {
	var strs []string
	err := json.Unmarshal(b, &strs)
	return strs, err
}
func constStringPtr(s string) *int8 {
	return (*int8)(unsafe.Pointer(&[]byte(s + "\x00")[0]))
}
func clipboard_write_files(bytes []byte) bool {
	files, err := byte_slice_to_string_slice(bytes)
	if err != nil {
		return false
	}
	var filePathsPtrs []unsafe.Pointer
	for _, path := range files {
		fmt.Println(path)
		nsString := objc.ID(class_NSString).Send(sel_stringWithUTF8String, unsafe.Pointer(constStringPtr(path)))
		filePathsPtrs = append(filePathsPtrs, unsafe.Pointer(nsString))
	}
	// nsArray := objc.ID(class_NSArray)
	pasteboard := objc.ID(class_NSPasteboard).Send(sel_generalPasteboard)
	// 清空粘贴板内容
	pasteboard.Send(sel_clearContents)
	// var nsArrayClass = objc.GetClass("NSArray")

	nsArray := objc.ID(class_NSArray).Send(sel_arrayWithObjects_count, unsafe.Pointer(&filePathsPtrs[0]), len(filePathsPtrs))

	// 将文件路径数组设置到粘贴板
	return pasteboard.Send(sel_setPropertyList_forType_, nsArray, _NSPasteboardTypeFiles) != 0
}

func clipboard_change_count() int {
	return int(objc.ID(class_NSPasteboard).Send(sel_generalPasteboard).Send(sel_changeCount))
}
