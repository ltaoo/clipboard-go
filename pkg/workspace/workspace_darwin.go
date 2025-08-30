//go:build darwin && !ios

package workspace

import (
	"errors"
	"unsafe"

	"github.com/ebitengine/purego"
	"github.com/ebitengine/purego/objc"
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

var (
	appkit = must(purego.Dlopen("/System/Library/Frameworks/AppKit.framework/AppKit", purego.RTLD_GLOBAL|purego.RTLD_NOW))

	_NSWorkspace          = objc.GetClass("NSWorkspace")
	_sharedWorkspace      = objc.RegisterName("sharedWorkspace")
	_frontmostApplication = objc.RegisterName("frontmostApplication")
	_localizedName        = objc.RegisterName("localizedName")

	_UTF8String = objc.RegisterName("UTF8String")

	_length = objc.RegisterName("length")
)

func cur_application() (string, error) {
	__workspace := objc.ID(_NSWorkspace).Send(_sharedWorkspace)
	if __workspace == 0 {
		return "", errors.New("初始化失败")
	}
	__app := __workspace.Send(_frontmostApplication)
	if __app == 0 {
		return "", errors.New("没有获取到当前应用")
	}
	__name := __app.Send(_localizedName)
	utf8_ptr := unsafe.Pointer(__name.Send(_UTF8String))
	if utf8_ptr == nil {
		return "", errors.New("读取失败")
	}
	return read_utf8_string(utf8_ptr), nil
}

func read_utf8_string(ptr unsafe.Pointer) string {
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
