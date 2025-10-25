// g++ cpp/read_image.cpp -o read_image.exe

#include <iostream>
#include <windows.h>
#include <fstream>
#include <cstdlib>
#include <cstring>

// 用于将DIB转换为JPEG的库（这里假设使用libjpeg库，需要链接相应库文件）
#include <jpeglib.h>

// 自定义函数：将DIB数据转换为JPEG格式
bool dibToJpeg(const char* dibData, DWORD dibSize, const char* jpegFilePath) {
    // 分配JPEG压缩对象
    struct jpeg_compress_struct cinfo;
    struct jpeg_error_mgr jerr;
    FILE* outfile;

    cinfo.err = jpeg_std_error(&jerr);
    jpeg_create_compress(&cinfo);

    // 打开输出文件
    if ((outfile = fopen(jpegFilePath, "wb")) == NULL) {
        std::cerr << "无法打开输出文件: " << jpegFilePath << std::endl;
        return false;
    }
    jpeg_stdio_dest(&cinfo, outfile);

    // 设置压缩参数
    cinfo.image_width = 0; // 需要从DIB头中获取
    cinfo.image_height = 0; // 需要从DIB头中获取
    cinfo.input_components = 3; // 假设为RGB
    cinfo.in_color_space = JCS_RGB;

    jpeg_set_defaults(&cinfo);
    jpeg_set_quality(&cinfo, 75, TRUE);

    // 开始压缩
    jpeg_start_compress(&cinfo, TRUE);

    // 这里需要解析DIB头获取图像尺寸等信息
    BITMAPV5HEADER* bmpHeader = reinterpret_cast<BITMAPV5HEADER*>(const_cast<char*>(dibData));
    cinfo.image_width = bmpHeader->bV5Width;
    cinfo.image_height = bmpHeader->bV5Height;

    // 每行扫描线的字节数，需要根据图像宽度和颜色模式计算
    int rowStride = (cinfo.image_width * 3 + 3) & ~3;
    JSAMPROW row_pointer[1];

    DWORD dataOffset = sizeof(BITMAPV5HEADER);
    for (int i = 0; i < cinfo.image_height; i++) {
        row_pointer[0] = dibData + dataOffset + (cinfo.image_height - i - 1) * rowStride;
        jpeg_write_scanlines(&cinfo, row_pointer, 1);
    }

    // 完成压缩
    jpeg_finish_compress(&cinfo);
    fclose(outfile);
    jpeg_destroy_compress(&cinfo);

    return true;
}

int main() {
    // 获取剪贴板数据
    if (!OpenClipboard(NULL)) {
        std::cerr << "无法打开剪贴板" << std::endl;
        return 1;
    }

    HGLOBAL hClipboardData = GetClipboardData(CF_DIB);
    if (!hClipboardData) {
        std::cerr << "剪贴板中没有DIB格式的数据" << std::endl;
        CloseClipboard();
        return 1;
    }

    char* pData = (char*)GlobalLock(hClipboardData);
    if (!pData) {
        std::cerr << "无法锁定剪贴板数据" << std::endl;
        CloseClipboard();
        return 1;
    }

    DWORD dwSize = GlobalSize(hClipboardData);

    // 生成JPEG文件
    if (!dibToJpeg(pData, dwSize, "clipboard_image.jpg")) {
        GlobalUnlock(hClipboardData);
        CloseClipboard();
        return 1;
    }

    // 解锁内存块
    GlobalUnlock(hClipboardData);

    // 关闭剪贴板
    CloseClipboard();

    std::cout << "已将剪贴板数据保存为clipboard_image.jpg" << std::endl;

    return 0;
}
