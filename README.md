# dir2md
依赖于pandoc，用于将指定目录下的文件【html，docx，doc】并发转换成md文件，同时对于doc文件可以在相对路径下解压出图片等资源。

pandoc安装链接：https://www.pandoc.org/installing.html 安装完后，请确认pandoc命令在PATH中。

## Uages
```
Usage: dir2md inPath [outPath]
  -i string
        输入文件路径
  -o string
        输出文件路径 (default "%inPath%_out")

```

## Example
比如有test目录，内容如下：
```
.
├── IOT安全
│   ├── Cisco
│   │   └── （CVE-2019-1663）Cisco 堆栈缓冲区溢出漏洞.docx
│   ├── D-Link
│   │   ├── （CVE-2018-19986）D-Link DIR-818LW&828命令注入漏洞.docx
│   │   ├── （CVE-2018-20056）D-Link DIR-619L&605L 栈溢出漏洞.docx
│   │   ├── （CVE-2018-20057）D-Link DIR-619L&605L 命令注入漏洞.docx
```
运行命令：
```bash
git clone https://github.com/XinRoom/dir2md.git --depth=1
cd dir2md
go get -u .
go run . -i /path/to/test
```
将会在test的同级目录下生成test_out目录和文件：
```
.
├── IOT安全
│   ├── Cisco
│   │   ├── （CVE-2019-1663）Cisco 堆栈缓冲区溢出漏洞.md
│   │   └── media
│   │       ├── rId25.png
│   │       ├── rId26.png
│   │       ├── rId28.png
│   │       ├── rId30.png
│   │       ├── rId31.png
│   ├── D-Link
│   │   ├── （CVE-2018-19986）D-Link DIR-818LW&828命令注入漏洞.md
│   │   ├── （CVE-2018-20056）D-Link DIR-619L&605L 栈溢出漏洞.md
│   │   ├── （CVE-2018-20057）D-Link DIR-619L&605L 命令注入漏洞.md
│   │   └── media
│   │       ├── rId23.png
│   │       ├── rId24.png
│   │       ├── rId25.png
```

## Tip
如需编译使用可以使用以下命令
```
go build -trimpath -ldflags="-s -w -buildid=" .
```