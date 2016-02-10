package main

import (
    "log"

    "github.com/kelseyhightower/targz"
)

func main() {
    err := targz.Create("/Baidu/Project/canary/codecenter/src/eggs/packer.go", "stuff.tar.gz")
    if err != nil {
        log.Fatalln(err)
    }
}
