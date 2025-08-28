#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o clipboard_write_text objective-c/write_text.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSString *text = @"hello world";
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	[pasteboard clearContents];
	// 设置要写入粘贴板的文本
	if ([pasteboard setString:text forType:NSPasteboardTypeString]) {
            NSLog(@"成功写入文本粘贴板");
        } else {
            NSLog(@"无法写入文本粘贴板");
        }
    }
    return 0;
}
