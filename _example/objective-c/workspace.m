#import <Foundation/Foundation.h>
#import <AppKit/AppKit.h> 
#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o workspace objective-c/workspace.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
	// 获取 NSWorkspace 共享实例
        NSWorkspace *workspace = [NSWorkspace sharedWorkspace];
        // 获取前台应用程序
        NSRunningApplication *frontApp = [workspace frontmostApplication];
        if (frontApp) {
            // 获取应用程序的本地化名称
            NSString *appName = frontApp.localizedName;
            NSLog(@"当前前台应用: %@", appName);
        } else {
            NSLog(@"未能获取到前台应用");
        }
    }
    return 0;
}
