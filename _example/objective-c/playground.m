#import <Cocoa/Cocoa.h>

// gcc -framework Cocoa -o playground objective-c/playground.m
int main(int argc, const char * argv[]) {
    @autoreleasepool {
        NSString *file1 = @"/Users/mayfair/Documents/deploy_step2.png";

	NSUInteger length = [file1 lengthOfBytesUsingEncoding:NSUTF8StringEncoding];
	NSMutableData *v = [NSMutableData dataWithLength:length];
    }
    return 0;
}
