// Copyright 2021 The golang.design Initiative Authors.
// All rights reserved. Use of this source code is governed
// by a MIT license that can be found in the LICENSE file.
//
// Written by Changkun Ou <changkun.de>

//go:build darwin && !ios

// Interact with NSPasteboard using Objective-C
// https://developer.apple.com/documentation/appkit/nspasteboard?language=objc

#import <Foundation/Foundation.h>
#import <Cocoa/Cocoa.h>

unsigned int clipboard_read_string(void **out) {
	NSPasteboard * pasteboard = [NSPasteboard generalPasteboard];
	NSData *data = [pasteboard dataForType:NSPasteboardTypeString];
	if (data == nil) {
		return 0;
	}
	NSUInteger siz = [data length];
	*out = malloc(siz);
	[data getBytes: *out length: siz];
	return siz;
}

// let png_data = unsafe { self.pasteboard.dataForType(NSPasteboardTypePNG) };
// if let Some(data) = png_data {
// 	return RustImageData::from_bytes(&data.to_vec());
// };
// // if no png data, read NSImage;
// let ns_image =
// 	unsafe { NSImage::initWithPasteboard(NSImage::alloc(), &self.pasteboard) };
// if let Some(image) = ns_image {
// 	let tiff_data = unsafe { image.TIFFRepresentation() };
// 	if let Some(data) = tiff_data {
// 		return RustImageData::from_bytes(&data.to_vec());
// 	}
// };
// Err("no image data".into())
unsigned int clipboard_read_image(void **out) {
	NSPasteboard * pasteboard = [NSPasteboard generalPasteboard];
	NSData *data = [pasteboard dataForType:NSPasteboardTypePNG];
	if (data == nil) {
		return 0;
	}
	NSUInteger siz = [data length];
	*out = malloc(siz);
	[data getBytes: *out length: siz];
	return siz;
}

// fn get_files(&self) -> Result<Vec<String>> {
// 	let mut res = vec![];
// 	let ns_array = unsafe { self.pasteboard.propertyListForType(NSFilenamesPboardType) };
// 	unsafe {
// 		if let Some(array) = ns_array {
// 			// cast to NSArray<NSString>
// 			let array: Retained<NSArray<NSString>> = Retained::cast_unchecked(array);
// 			array.iter().for_each(|item| {
// 				res.push(item.to_string());
// 			});
// 		}
// 	}
// 	if res.is_empty() {
// 		return Err("no files".into());
// 	}
// 	Ok(res)
// }
unsigned int clipboard_get_files(void **out) {
	*out = NULL;
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	// NSArray *files = [pasteboard propertyListForType:NSFilenamesPboardType];
	// id pbData = [pasteboard propertyListForType:NSPasteboardTypeFileURL];
	id pbData = [pasteboard propertyListForType:NSFilenamesPboardType];
	// 检查数据有效性：必须是 NSArray 类型且非空
	if (!pbData || ![pbData isKindOfClass:[NSArray class]]) {
		return 1; // 错误：无数据或数据类型错误
	}
	NSArray *filePaths = (NSArray *)pbData;
	if (filePaths.count == 0) {
		return 2; // 错误：剪贴板中无文件
	}
	// 分配 C 字符串数组内存（+1 用于结尾的 NULL 标记）
	char **result = (char **)malloc((filePaths.count + 1) * sizeof(char *));
	if (!result) {
		return 3; // 错误：内存分配失败
	}
	// 遍历 NSArray 转换每个路径为 C 字符串
	for (NSUInteger i = 0; i < filePaths.count; i++) {
		id item = filePaths[i];
		// 检查数组元素是否为 NSString 类型
		if (![item isKindOfClass:[NSString class]]) {
			// 清理已分配的内存
			for (NSUInteger j = 0; j < i; j++) {
				free(result[j]);
			}
			free(result);
			return 4; // 错误：数组包含非字符串元素
		}

		// 转换 NSString 为 C 字符串并复制（避免原对象释放后失效）
		NSString *path = (NSString *)item;
		const char *cPath = [path UTF8String];
		if (!cPath) { // 理论上 NSString 的 UTF8String 不会为 NULL，但防御性检查
			for (NSUInteger j = 0; j < i; j++) {
				free(result[j]);
			}
			free(result);
			return 5; // 错误：字符串无法转换为 C 格式
		}

		// 复制字符串到新分配的内存（需手动释放）
		result[i] = strdup(cPath);
		if (!result[i]) { // strdup 可能因内存不足失败
			for (NSUInteger j = 0; j < i; j++) {
				free(result[j]);
			}
			free(result);
			return 6; // 错误：内存分配失败
		}
	}

	// 添加 NULL 结尾标记（方便遍历）
	result[filePaths.count] = NULL;
	*out = result;

	return 0; // 成功
}

int clipboard_write_string(const void *bytes, NSInteger n) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSData *data = [NSData dataWithBytes: bytes length: n];
	[pasteboard clearContents];
	BOOL ok = [pasteboard setData: data forType:NSPasteboardTypeString];
	if (!ok) {
		return -1;
	}
	return 0;
}
int clipboard_write_image(const void *bytes, NSInteger n) {
	NSPasteboard *pasteboard = [NSPasteboard generalPasteboard];
	NSData *data = [NSData dataWithBytes: bytes length: n];
	[pasteboard clearContents];
	BOOL ok = [pasteboard setData: data forType:NSPasteboardTypePNG];
	if (!ok) {
		return -1;
	}
	return 0;
}

NSInteger clipboard_change_count() {
	return [[NSPasteboard generalPasteboard] changeCount];
}
