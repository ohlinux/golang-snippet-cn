package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	skillfolder := `/tmp`
	// 获取所有文件
	files, _ := ioutil.ReadDir(skillfolder)
	for _, file := range files {
		if file.IsDir() {
			continue
		} else {
			fmt.Println(file.Name())
		}
	}
}
