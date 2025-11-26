package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ltaoo/clipboard-go"
	"github.com/ltaoo/clipboard-go/pkg/util"
)

func main() {
	// Use a cancellable context to properly manage the Watch goroutine lifecycle
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure resources are cleaned up

	// Set up signal handling for graceful shutdown (Ctrl+C)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Println("\nReceived interrupt signal, shutting down gracefully...")
		cancel() // Cancel context to stop the watcher
	}()

	ch := clipboard.Watch(ctx)
	fmt.Println("Start watch the clipboard... (Press Ctrl+C to stop)")
	for data := range ch {
		fmt.Println(data.Type)
		// types := clipboard.GetContentTypes()
		// fmt.Println(types)

		if data.Type == "public.utf8-plain-text" {
			if text, ok := data.Data.(string); ok {
				fmt.Println(text)
			}
		}
		if data.Type == "public.html" {
			if text, ok := data.Data.(string); ok {
				fmt.Println(text)
			}
		}
		if data.Type == "public.png" {
			if f, ok := data.Data.([]byte); ok {
				img_filepath, err := util.SaveByteAsLocalImage(f)
				if err == nil {
					fmt.Println("the image save to", img_filepath)
				}
			}
		}
		if data.Type == "public.file-url" {
			if files, ok := data.Data.([]string); ok {
				for _, f := range files {
					fmt.Println(f)
				}
			}
		}
	}
	fmt.Println("Clipboard watcher stopped gracefully")
}
