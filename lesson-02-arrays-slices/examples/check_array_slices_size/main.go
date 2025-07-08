package main

import (
	"fmt"
	"unsafe"
)

func main() {
	array := [255]byte{}
	slice := make([]byte, 255, 255)

	fmt.Printf("array size: %d; pointer to array size: %d\n", unsafe.Sizeof(array), unsafe.Sizeof(&array))
	fmt.Printf("slice size: %d; pointer to slice size: %d\n", unsafe.Sizeof(slice), unsafe.Sizeof(&slice))
}
