Add-Type -AssemblyName System.Windows.Forms

# 读取已经包含完整剪贴板格式的HTML文件
$htmlFilePath = "./powershell/debug_clipboard.html"
$clipboardData = [System.IO.File]::ReadAllText($htmlFilePath)

# 注册HTML格式
$htmlFormat = [System.Windows.Forms.DataFormats]::GetFormat("HTML Format")

# 创建DataObject并设置数据
$dataObject = New-Object System.Windows.Forms.DataObject
$dataObject.SetData($htmlFormat.Name, $clipboardData)

# 设置剪贴板内容
[System.Windows.Forms.Clipboard]::SetDataObject($dataObject, $true)

Write-Host "HTML格式数据已成功写入剪贴板"
