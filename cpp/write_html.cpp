#include <Windows.h>
#include <comdef.h>
#include <vector>
#include <string>
#include <fstream>
#include <iostream>

// 读取文件内容到字符串
std::string ReadFileToString(const std::string& filePath) {
    std::ifstream file(filePath, std::ios::binary);
    if (!file.is_open()) {
        throw std::runtime_error("无法打开文件: " + filePath);
    }

    file.seekg(0, std::ios::end);
    size_t fileSize = file.tellg();
    file.seekg(0, std::ios::beg);

    std::vector<char> buffer(fileSize);
    file.read(buffer.data(), fileSize);

    return std::string(buffer.data(), fileSize);
}

int main() {
    // 初始化COM库
    if (FAILED(CoInitialize(NULL))) {
        std::cerr << "COM初始化失败" << std::endl;
        return 1;
    }

    try {
        // 读取HTML文件
        std::string htmlFilePath = "./debug_clipboard.html";
        std::string clipboardData = ReadFileToString(htmlFilePath);

        // 打开剪贴板
        if (!OpenClipboard(NULL)) {
            throw std::runtime_error("无法打开剪贴板");
        }

        // 清空剪贴板
        if (!EmptyClipboard()) {
            CloseClipboard();
            throw std::runtime_error("无法清空剪贴板");
        }

        // 注册HTML格式
        UINT htmlFormat = RegisterClipboardFormatA("HTML Format");
        UINT size = clipboardData.size() + 1;
        std::cout << "[DEBUG]size" << size << std::endl;
        // 分配全局内存
        HGLOBAL hClipboardData = GlobalAlloc(GMEM_MOVEABLE, size);
        if (hClipboardData == NULL) {
            CloseClipboard();
            throw std::runtime_error("内存分配失败");
        }

        // 锁定内存并复制数据
        char* pchData = static_cast<char*>(GlobalLock(hClipboardData));
        if (pchData == NULL) {
            GlobalFree(hClipboardData);
            CloseClipboard();
            throw std::runtime_error("无法锁定内存");
        }

        memcpy(pchData, clipboardData.c_str(), size);
        GlobalUnlock(hClipboardData);

        // 设置剪贴板数据
        if (SetClipboardData(htmlFormat, hClipboardData) == NULL) {
            GlobalFree(hClipboardData);
            CloseClipboard();
            throw std::runtime_error("无法设置剪贴板数据");
        }

        // 关闭剪贴板
        CloseClipboard();

        std::cout << "HTML格式数据已成功写入剪贴板" << std::endl;
    }
    catch (const std::exception& e) {
        std::cerr << "错误: " << e.what() << std::endl;
        CoUninitialize();
        return 1;
    }

    // 释放COM库
    CoUninitialize();
    return 0;
}
