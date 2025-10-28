// Package clipboard provides a simple, cross‑platform API for interacting
// with the system clipboard. It supports reading and writing plain text,
// HTML, images (PNG-encoded), and file paths, as well as watching for
// clipboard changes.
//
// Typical usage:
//
//	// Initialize once at program startup
//	if err := clipboard.Init(); err != nil {
//		panic(err)
//	}
//
//	// Write and read plain text
//	if err := clipboard.WriteText("hello"); err != nil {
//		panic(err)
//	}
//	text, _ := clipboard.ReadText()
//	_ = text
//
//	// Watch clipboard updates
//	ctx, cancel := context.WithCancel(context.Background())
//	defer cancel()
//	for change := range clipboard.Watch(ctx) {
//		_ = change // handle ClipboardContent
//	}
//
// Platform notes:
//   - On some operating systems (e.g., Darwin), concurrent reads are not safe;
//     this package serializes access internally.
//   - If initialization fails, subsequent calls may panic.
package clipboard

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var (
	err_unavailable = errors.New("clipboard unavailable")
	err_unsupported = errors.New("unsupported format")
)

// Format represents the format of clipboard data.
type Format int

// All sorts of supported clipboard data.
const (
	// FmtText indicates plain text clipboard format
	FmtText Format = iota
	// FmtImage indicates image/png clipboard format
	FmtImage
	FmtFilepath
)

type ClipboardContent struct {
	Type       string // public.utf8-plain-text纯文本 public.file-url文件 public.png图片 public.html富文本
	Data       interface{}
	BackupData interface{}
	Error      error
}

var (
	// Due to the limitation on operating systems (such as darwin),
	// concurrent read can even cause panic, use a global lock to
	// guarantee one read at a time.
	lock      = sync.Mutex{}
	initOnce  sync.Once
	initError error
)

// Init initializes the clipboard package. It returns an error
// if the clipboard is not available to use. This may happen if the
// target system lacks required dependency, such as libx11-dev in X11
// environment. For example,
//
//	err := clipboard.Init()
//	if err != nil {
//		panic(err)
//	}
//
// If Init returns an error, any subsequent Read/Write/Watch call
// may result in an unrecoverable panic.
func Init() error {
	initOnce.Do(func() {
		initError = initialize()
	})
	return initError
}

// Read returns a chunk of bytes of the clipboard data if it presents
// in the desired format t presents. Otherwise, it returns nil.
func Read(t Format) ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()
	buf, err := read(t)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

// ReadText returns the current clipboard contents as UTF-8 text.
// It returns an empty string and a non-nil error on failure.
func ReadText() (string, error) {
	lock.Lock()
	defer lock.Unlock()
	t, err := read_text()
	if err != nil {
		return "", nil
	}
	return t, nil
}

// ReadHTML returns the current clipboard contents as HTML.
// If available, the HTML markup string is returned; otherwise an error.
func ReadHTML() (string, error) {
	lock.Lock()
	defer lock.Unlock()
	t, err := read_html()
	if err != nil {
		return "", nil
	}
	return t, nil
}

// ReadImage returns the current clipboard image encoded as PNG bytes.
// It returns an error if no image data is available.
func ReadImage() ([]byte, error) {
	lock.Lock()
	defer lock.Unlock()
	return read_image()
}

// ReadFiles returns file paths currently in the clipboard (if any).
// Paths are absolute when provided by the OS.
func ReadFiles() ([]string, error) {
	lock.Lock()
	defer lock.Unlock()
	return read_files()
}

// Write writes a given buffer to the clipboard in a specified format.
// Write returned a receive-only channel can receive an empty struct
// as a signal, which indicates the clipboard has been overwritten from
// this write.
// If format t indicates an image, then the given buf assumes
// the image data is PNG encoded.
func Write(t Format, buf []byte) (<-chan struct{}, error) {
	lock.Lock()
	defer lock.Unlock()
	changed, err := write(t, buf)
	if err != nil {
		return nil, err
	}
	return changed, nil
}

// WriteText writes a UTF-8 string to the clipboard as plain text.
func WriteText(text string) error {
	lock.Lock()
	defer lock.Unlock()
	return write_text(text)
}

// WriteHTML writes HTML markup to the clipboard.
// Optionally provide a plain-text fallback as the first variadic argument.
func WriteHTML(html string, args ...string) error {
	lock.Lock()
	defer lock.Unlock()
	var text *string
	if len(args) > 0 {
		text = &args[0]
	}
	return write_html(html, text)
}

// WriteImage writes an image to the clipboard. Data must be PNG-encoded.
func WriteImage(data []byte) error {
	lock.Lock()
	defer lock.Unlock()
	return write_image(data)
}

// WriteFiles writes a list of file paths to the clipboard.
func WriteFiles(files []string) error {
	lock.Lock()
	defer lock.Unlock()
	return write_files(files)
}

// Watch returns a receive-only channel that received the clipboard data
// whenever any change of clipboard data in the desired format happens.
//
// The returned channel will be closed if the given context is canceled.
// Watch emits clipboard changes as `ClipboardContent` values until the
// provided context is canceled. Only one watcher should be active per process.
func Watch(ctx context.Context) <-chan ClipboardContent {
	return watch(ctx)
}

// ContentTypeParams specifies options for querying supported content types.
type ContentTypeParams struct {
	IsEnabled bool // 粘贴板已经处于可用状态
}

// GetContentTypes returns the available clipboard content type identifiers
// for the current platform, optionally gated by `ContentTypeParams`.
func GetContentTypes(params ContentTypeParams) []string {
	return get_content_types(params)
}

// ByteToStrArray decodes a JSON array of strings from bytes.
// Returns the decoded slice or an error.
func ByteToStrArray(b []byte) ([]string, error) {
	var strs []string
	err := json.Unmarshal(b, &strs)
	return strs, err
}

// StrArrayToByte concatenates a slice of strings into a single byte slice
// without separators. This is a utility helper and does not perform encoding.
func StrArrayToByte(strs []string) []byte {
	var total_len int
	for _, str := range strs {
		total_len += len(str)
	}
	result := make([]byte, total_len)
	offset := 0
	for _, str := range strs {
		str_bytes := []byte(str)
		copy(result[offset:], str_bytes)
		offset += len(str_bytes)
	}
	return result
}
