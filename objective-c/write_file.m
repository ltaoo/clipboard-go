#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o clipboard_write_file objective-c/write_file.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSString *file1 = @"/Users/mayfair/Documents/deploy_step2.png";
        NSString *file2 = @"/Users/mayfair/Documents/deploy_step4.png";
	    NSArray <NSString *> *files = @[file1, file2];
        NSMutableArray <NSURL *> *fileURLs = [[NSMutableArray alloc] init];
        for (NSString *path in files) {
            NSURL* fileURL = [[NSURL alloc] initFileURLWithPath:path];
            if ([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:NULL]) {
                [fileURLs addObject:fileURL];
            } else {
                NSLog(@"文件 %@ 不存在", path);
            }
        }
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];
        if ([pasteboard writeObjects:fileURLs]) {
            NSLog(@"多个文件已成功写入粘贴板");
        } else {
            NSLog(@"无法写入文件到粘贴板");
        }
    }
    return 0;
}
