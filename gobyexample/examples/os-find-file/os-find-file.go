package main

import (
	"fmt"
	"os"
	"regexp"
)

func main() {
	// 打开文件夹
	dir, err := os.Open("/tmp")
	if err != nil {
		panic(nil)
	}
	defer dir.Close()

	// 读取文件列表
	fis, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// 遍历文件列表 符合正则的匹配的输出.
	for _, fi := range fis {

		match, _ := regexp.MatchString("mysql.*", fi.Name())
		if match {
			fmt.Println(fi.Name())
		}

	}
}
