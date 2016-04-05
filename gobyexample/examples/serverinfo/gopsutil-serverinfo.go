package main

import (
	"fmt"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	v, _ := mem.VirtualMemory()

	// almost every return value is a struct
	fmt.Printf("Total: %v, Free:%v, UsedPercent:%f%%\n", v.Total, v.Free, v.UsedPercent)

	// convert to JSON. String() is also implemented
	fmt.Println(v)

	h, err := host.HostInfo()
	if err != nil {
		fmt.Println("err:", err)
	} else {

		fmt.Printf("hostname %v", h)
	}

	c, err := cpu.CPUInfo()
	if err != nil {
		fmt.Println("err:", err)
	}
	for _, v := range c {
		fmt.Printf("cpu info %v \n ", v)
	}

}
