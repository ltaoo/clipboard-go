#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h> 
#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o clipboard_write_file objective-c/write_file.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        // NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        // [pasteboard clearContents];
	    // NSArray <NSString *> *files = @[@"/Users/mayfair/Documents/deploy_step2.png", @"/Users/mayfair/Documents/StatsCard.tsx"];
        // NSMutableArray <NSURL *> *fileURLs = [[NSMutableArray alloc] init];
        // for (NSString *path in files) {
        //     NSURL* fileURL = [[NSURL alloc] initFileURLWithPath:path];
        //     if ([[NSFileManager defaultManager] fileExistsAtPath:path isDirectory:NULL]) {
        //         NSLog(@"新增文件 %@", path);
        //         [fileURLs addObject:fileURL];
        //     } else {
        //         NSLog(@"文件 %@ 不存在", path);
        //     }
        // }
        // // [pasteboard setTypes:@[NSURLPboardType] owner:nil];
        // // [pasteboard declareTypes:@[NSFilenamesPboardType] owner:nil];
        // // CFStringRef fileURLUTI = UTTypeCreatePreferredIdentifierForTag(kUTTagClassFilenameExtension, CFSTR("url"), NULL);
        // if ([pasteboard writeObjects:fileURLs]) {
        //     NSLog(@"多个文件已成功写入粘贴板");
        // } else {
        //     NSLog(@"无法写入文件到粘贴板");
        // }

// --------------
        NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
        [pasteboard clearContents];
        NSMutableArray *filesToCopy = [NSMutableArray array];
        NSString *filePath1 = @"/Users/litao/Downloads/avatar.png";
        // NSString *filePath1 = @"/Users/mayfair/Documents/deploy_step2.png";
        NSString *filePath2 = @"/Users/litao/Downloads/face.png";
        // NSString *filePath2 = @"/Users/mayfair/Documents/deploy_step4.png";
        NSURL *fileURL1 = [NSURL fileURLWithPath:filePath1];
        NSURL *fileURL2 = [NSURL fileURLWithPath:filePath2];
        if (fileURL1) {
            [filesToCopy addObject:fileURL1];
            // [filesToCopy addObject:filePath1];
            NSLog(@"fileURL2 absoluteString: %@", [fileURL1 absoluteString]);
        } else {
            NSLog(@"Error: Could not create URL for %@", filePath1);
        }
        if (fileURL2) {
            [filesToCopy addObject:fileURL2];
            // [filesToCopy addObject:filePath2];
            NSLog(@"fileURL2 absoluteString: %@", [fileURL2 absoluteString]);
        } else {
            NSLog(@"Error: Could not create URL for %@", filePath2);
        }
        // CFStringRef fileURLUTI = CFSTR("public.file-url");
        // UTTypeRef fileURLType = UTTypeCreate(fileURLUTI, NULL);
        // [pasteboard declareTypes:[NSArray arrayWithObject:NSPasteboardTypeFileURL] owner:nil];
        // [pasteboard setPropertyList:filesToCopy forType:NSPasteboardTypeFileURL];
        // [pasteboard setPropertyList:filesToCopy forType:NSPasteboardTypeFileURL];
        BOOL success = [pasteboard writeObjects:filesToCopy];
        if (success) {
            NSLog(@"Successfully copied %lu files to the clipboard.", (unsigned long)[filesToCopy count]);
        } else {
            NSLog(@"Failed to copy files to the clipboard.");
        }
// --------------
        // NSArray *fileList = [NSArray arrayWithObjects:filePath1, filePath2, nil];
        // // NSPasteboard *pboard = [NSPasteboard generalPasteboard];
        // [pasteboard declareTypes:[NSArray arrayWithObject:NSFilenamesPboardType] owner:nil];
        // [pasteboard setPropertyList:fileList forType:NSFilenamesPboardType];
    }
    return 0;
}
