package clipboard_test

import (
	"context"
	"fmt"
	"time"

	clipboard "github.com/ltaoo/clipboard-go"
)

// ExampleInit demonstrates initializing the clipboard.
func ExampleInit() {
	if err := clipboard.Init(); err != nil {
		fmt.Println("init failed:", err)
		return
	}
	fmt.Println("ok")
	// Output:
	// ok
}

// ExampleReadText demonstrates reading text from the clipboard.
func ExampleReadText() {
	_ = clipboard.Init()
	text, err := clipboard.ReadText()
	if err != nil {
		fmt.Println("read failed:", err)
		return
	}
	fmt.Println(len(text) >= 0)
	// Output:
	// true
}

// ExampleWriteText demonstrates writing plain text to the clipboard.
func ExampleWriteText() {
	_ = clipboard.Init()
	if err := clipboard.WriteText("hello clipboard"); err != nil {
		fmt.Println("write failed:", err)
		return
	}
	// Output:
}

// ExampleWatch demonstrates watching for clipboard changes.
func ExampleWatch() {
	_ = clipboard.Init()
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	for change := range clipboard.Watch(ctx) {
		_ = change // handle ClipboardContent
	}
	// Output:
}
