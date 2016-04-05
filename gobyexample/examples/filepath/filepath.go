package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {

	// 获取执行文件所在的绝对路径
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	fmt.Println(path)
}
