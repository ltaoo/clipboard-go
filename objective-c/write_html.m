#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o clipboard_write_text objective-c/write_html.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSString *html = @"<!DOCTYPE html>\n"
                        @"<html>\n"
                        @"<body>\n"
                        @"    <h1>这是粘贴板中的HTML</h1>\n"
                        @"</body>\n"
                        @"</html>";
        NSString *plain_text = @"这是HTML内容的纯文本回退版本";

        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];

        // NSData *html_content = [html dataUsingEncoding:NSUTF8StringEncoding];
        // [pasteboard declareTypes:@[NSPasteboardTypeHTML, NSPasteboardTypeString] owner:nil];
        // [pasteboard setData:html_content forType:NSPasteboardTypeHTML];
        if ([pasteboard setString:html forType:NSPasteboardTypeHTML]) {
            NSLog(@"成功写入文本粘贴板");
        } else {
            NSLog(@"无法写入文本粘贴板");
        }
        //  NSDictionary *html_content = @{
        //     NSPasteboardTypeHTML: html,
        //     NSPasteboardTypeString: plain_text ?: @""
        // };
        // 设置要写入粘贴板的文本
        // if ([pasteboard writeObjects:@[html_content] forType:NSPasteboardTypeHTML]) {
        //         NSLog(@"成功写入文本粘贴板");
        // } else {
        //     NSLog(@"无法写入文本粘贴板");
        // }
    }
    return 0;
}
