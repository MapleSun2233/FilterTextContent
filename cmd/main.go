package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	url := flag.String("url", "", "文本流请求地址，url或file必指定其一，同时指定时优先url")
	filePath := flag.String("file", "", "文本文件地址，url或file必指定其一，同时指定时优先url")
	charset := flag.String("charset", "UTF-8", "文本字符集，默认UTF-8")
	cookieValue := flag.String("cookieStr", "", "携带cookie")
	prefix := flag.String("prefix", "", "过滤行前缀")
	endPrefix := flag.String("endPrefix", "", "过滤行前缀")
	feature := flag.String("feature", "", "过滤行特征")
	outputFile := flag.String("out", "out.log", "日志输出路径")
	maxBuff := flag.Int("buff", 1048576, "行缓冲大小，行太大时指定，默认1MB")
	flag.Parse()
	if len(*url) == 0 && len(*filePath) == 0 {
		fmt.Println("url或file必指定其一")
		return
	}
	fmt.Printf(`参数信息：
	url: %s
	file: %s
	charset: %s
	cookieValue: %s
	prefix: %s
	endPrefix: %s
	feature: %s
	maxBuff: %d
	outFilePath: %s
`, *url, *filePath, *charset, *cookieValue, *prefix, *endPrefix, *feature, *maxBuff, *outputFile)
	file, err := os.Create(*outputFile)
	if err != nil {
		fmt.Printf("Failed to create file: %s\n", err)
		return
	}
	defer file.Close()

	var scanner *bufio.Scanner
	var decoder *encoding.Decoder
	switch strings.ToLower(*charset) {
	case "gbk":
		decoder = simplifiedchinese.GBK.NewDecoder()
	case "gb2312":
		decoder = simplifiedchinese.HZGB2312.NewDecoder()
	case "gb18030":
		decoder = simplifiedchinese.GB18030.NewDecoder()
	default:
	}

	if len(*url) != 0 {
		client := &http.Client{}
		req, err := http.NewRequest("GET", *url, nil)
		if err != nil {
			fmt.Printf("Failed to create request: %s\n", err)
			return
		}
		if len(*cookieValue) > 0 {
			*cookieValue = strings.ReplaceAll(*cookieValue, " ", "")
			cookies := strings.Split(*cookieValue, ";")
			for _, cookieStr := range cookies {
				if strings.Contains(cookieStr, "=") {
					fmt.Println("AddToken:" + cookieStr)
					tokens := strings.Split(cookieStr, "=")
					cookie := &http.Cookie{Name: tokens[0], Value: tokens[1]}
					req.AddCookie(cookie)
				}
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Failed to send request: %s\n", err)
			return
		}
		defer resp.Body.Close()
		scanner = bufio.NewScanner(resp.Body)
	} else if len(*filePath) != 0 {
		sourceFile, err := os.OpenFile(*filePath, os.O_RDONLY, 0777)
		if err != nil {
			fmt.Printf("Fail to read file: %s\n", err.Error())
		}
		scanner = bufio.NewScanner(bufio.NewReader(sourceFile))
	}

	buf := make([]byte, *maxBuff)
	scanner.Buffer(buf, *maxBuff)
	isWrite := false
	isMatchPrefix := len(*prefix) > 0
	isMatchEndPrefix := len(*endPrefix) > 0
	isMatchFeature := len(*feature) > 0
	if !isMatchPrefix && !isMatchFeature {
		fmt.Println("未指定前缀特征，读取日志全文...")
		isWrite = true
	}
	for scanner.Scan() {
		line := scanner.Text()
		if !isWrite && isMatchPrefix && (strings.HasPrefix(line, *prefix)) {
			fmt.Println("发现匹配前缀特征，开始写入...")
			isWrite = true
		}
		if !isWrite && isMatchFeature && strings.Contains(line, *feature) {
			fmt.Println("发现匹配关键词特征，开始写入...")
			isWrite = true
		}
		if isWrite {
			if decoder != nil {
				line, _, err = transform.String(decoder, line)
				if err != nil {
					fmt.Printf("Failed to decode string: %s\n", err)
					return
				}
			}
			_, err = fmt.Fprintln(file, line)
			if err != nil {
				fmt.Printf("Failed to write to file: %s\n", err)
				return
			}
		}
		if isWrite && isMatchEndPrefix && (strings.HasPrefix(line, *endPrefix)) {
			fmt.Println("发现匹配结束前缀特征，结束写入...")
			return
		}
	}

	if err := scanner.Err(); err != nil && err != io.EOF {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Println("Successfully wrote response to file")
}
