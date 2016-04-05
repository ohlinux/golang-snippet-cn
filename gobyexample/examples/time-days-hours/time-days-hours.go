package main

import (
	"fmt"
	"time"

	"github.com/jinzhu/now"
)

func main() {
	//t := time.Now()

	t1 := time.Now()         //现在是12点整（假设）,那t1记录的就是12点整
	t2 := t1.Add(-time.Hour) //那t1的时间点 **加上(Add)** 1个小时，是几点呢？
	fmt.Println(t2)          //13点（呵呵）
	fmt.Println(t2.Format("2006010215"))
	fmt.Println(now.Monday())

	day := 7

	var i, h int
	for i = 1; i <= day; i++ {
		for h = 1; h <= 24; h++ {
			// -h*i 不能直接 t1.Add(-h*i*time.Hour) 会报错
			// invalid operation: time.Hour * h (mismatched types time.Duration and int)
			// 需要使用time.Duration来转变类型.
			t2 := t1.Add(time.Duration(-h*i) * time.Hour)
			fmt.Println(t2.Format("2006010215"))
		}
	}

	//time.Sleep(10 * time.Second)

}
