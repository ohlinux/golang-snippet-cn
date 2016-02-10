package main 

import (
    "fmt"
    "os"
)

func PathExist(_path string) (bool, error) {
    fi, err := os.Stat(_path)
    if err != nil && os.IsNotExist(err) {
        return false, err
    }
    if fi.IsDir() {
        return true, err
    }
    return false, err
}

func main() {


    src:="/Users/ajian/Baidu/Project/canary/codecenter/test/packer/tar2.go"

    isdir, _:= PathExist(src)

    var files []os.FileInfo

    if isdir {
        fmt.Println("is dir")
        dir, _:= os.Open(src)
        defer dir.Close()

        //获取文件列表
        files, _= dir.Readdir(0)

    }else{
        fmt.Println("is file")
        fi, _:= os.Stat(src)
        files =[]os.FileInfo{fi}
    }

    for _,file := range files {
        fmt.Println(file.Name())
    }

}
