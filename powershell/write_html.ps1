Add-Type -AssemblyName System.Windows.Forms

# 准备HTML内容
$htmlContent = @"
Version:0.9
StartHTML:0000000151
EndHTML:0000001349
StartFragment:0000000187
EndFragment:0000001313
SourceURL:https://sabrogden.github.io/Ditto/
<html>
<body>
<!--StartFragment--><h1 id="ditto---clipboard-manager" style="box-sizing: border-box; font-size: 2em; margin-top: 0px !important; margin-right: 0px; margin-bottom: 16px; margin-left: 0px; font-weight: 600; line-height: 1.25; padding-bottom: 0.3em; border-bottom: 1px solid rgb(234, 236, 239); color: rgb(36, 41, 46); font-family: -apple-system, BlinkMacSystemFont, &quot;Segoe UI&quot;, Helvetica, Arial, sans-serif, &quot;Apple Color Emoji&quot;, &quot;Segoe UI Emoji&quot;, &quot;Segoe UI Symbol&quot;; font-style: normal; font-variant-ligatures: normal; font-variant-caps: normal; letter-spacing: normal; orphans: 2; text-align: start; text-indent: 0px; text-transform: none; widows: 2; word-spacing: 0px; -webkit-text-stroke-width: 0px; white-space: normal; background-color: rgb(255, 255, 255); text-decoration-thickness: initial; text-decoration-style: initial; text-decoration-color: initial;"><a href="https://github.com/sabrogden/Ditto/releases/download/3.25.113.0/DittoSetup_3_25_113_0.exe" style="box-sizing: border-box; background-color: transparent; color: rgb(3, 102, 214); text-decoration: none;">Ditto - Clipboard Manager</a></h1><!--EndFragment-->
</body>
</html>
"@

# 组合完整数据
$clipboardData = $htmlContent

# 注册HTML格式
$htmlFormat = [System.Windows.Forms.DataFormats]::GetFormat("HTML Format")

# 创建DataObject并设置数据
$dataObject = New-Object System.Windows.Forms.DataObject
$dataObject.SetData($htmlFormat.Name, $clipboardData)

# 可选：添加纯文本后备
$plainText = "这是加粗文本，这是斜体文本"
$dataObject.SetData([System.Windows.Forms.DataFormats]::Text, $plainText)

# 设置剪贴板内容
[System.Windows.Forms.Clipboard]::SetDataObject($dataObject, $true)

Write-Host "HTML格式数据已成功写入剪贴板"
Write-Host "现在可以粘贴到Word、Outlook等支持富文本的应用程序中了"

