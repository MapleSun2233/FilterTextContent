# Filter Text Content
一个轻量级的流文本过滤器，数据来源支持文件或URL，旨为解决大日志的查看问题，从目标位置开始截取日志。亦可用于其他大文本的截取处理场景。

# Build
```shell
go build -o FilterTextContent ./cmd/main.go
```
# Usage
```shell
Usage of FilterTextContent:
  -url string
    	文本流请求地址，url或file必指定其一，同时指定时优先url
  -file string
    	文本文件地址，url或file必指定其一，同时指定时优先url
  -feature string
    	过滤行关键词特征
  -prefix string
    	过滤行开始前缀
  -endPrefix string
    	过滤行结束前缀
  -cookieStr string
    	携带cookie
  -buff int
    	行缓冲大小，行太大时指定，默认1MB (default 1048576)
  -charset string
    	文本字符集，默认UTF-8 (default "UTF-8")
  -out string
    	日志输出路径 (default "out.log")
```
# Tips
1. 来源为URL时可携带Cookies信息，绕过认证问题。
2. 同时指定过滤行开始前缀和结束前缀可实现范围截取。
3. 仅指定过滤行开始前缀或过滤行关键词特征时，发现匹配行开始截取，直至文件结束。
4. 若文件扫描到末尾仍未发现匹配内容，则结果文件为空。