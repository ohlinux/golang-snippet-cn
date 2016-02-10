package main

import (
    "fmt"
    "sync"
)

func main() {
    var wg sync.WaitGroup
    for i := 0; i < 100; i++ {
        fmt.Println("add",i)
        wg.Add(i)
    }

    for i := 0; i < 100; i++ {
        go wg.Done()
    }
    fmt.Println("exit")
    wg.Wait()
}

func add(wg sync.WaitGroup) {
    wg.Add(1)
}

func done(wg sync.WaitGroup) {
    wg.Done()
}
