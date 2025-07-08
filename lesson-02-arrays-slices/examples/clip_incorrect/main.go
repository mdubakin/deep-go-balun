package main

import (
	"fmt"
	"runtime"
	"slices"
	"unsafe"
)

func printAllocs() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("%d MB\n", m.Alloc/1024/1024)
}

func main() {
	data := make([]int, 10, 100<<20)
	fmt.Println("data:", unsafe.SliceData(data), len(data), cap(data))

	printAllocs()

	temp := slices.Clip(data) // data[:10:10]
	fmt.Println("temp2:", unsafe.SliceData(temp), len(temp), cap(temp))

	runtime.GC()
	printAllocs()

	runtime.KeepAlive(temp)
}
