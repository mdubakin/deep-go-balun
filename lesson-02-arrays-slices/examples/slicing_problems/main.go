package main

import (
	"fmt"
)

func main() {
	sl1 := make([]int, 2, 5) // [0, 0, -, -, -] len = 2; cap = 5;
	sl2 := sl1[1:3]          //    [0, 0, -, -] len = 2; cap = 4;

	fmt.Printf("Начало:\n\tsl1 len: %d; cap: %d\n\tsl2 len: %d; cap: %d\n", len(sl1), cap(sl1), len(sl2), cap(sl2))

	fmt.Printf("\n1 итерация (sl2 append 2):\n")
	sl2 = append(sl2, 2)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, -, -, -] len = 2; cap = 5;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 0, 2, -] len = 3; cap = 4;

	fmt.Printf("\n2 итерация (sl1 append 1):\n")
	sl1 = append(sl1, 1)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, -, -] len = 3; cap = 5;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 2, -] len = 3; cap = 4;

	fmt.Printf("\n3 итерация (sl1 append 1):\n")
	sl1 = append(sl1, 1)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, 1, -] len = 4; cap = 5;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 1, -] len = 3; cap = 4;

	fmt.Printf("\n4 итерация (sl1 append 1):\n")
	sl1 = append(sl1, 1)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, 1, 1] len = 5; cap = 5;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 1, -] len = 3; cap = 4;

	fmt.Printf("\n5 итерация (sl2 append 2):\n")
	sl2 = append(sl2, 2)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, 1, 2] len = 5; cap = 5;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 1, 2] len = 4; cap = 4;

	fmt.Printf("\nДо \"РАЗЛОМА\":\n\tsl1 len: %d; cap: %d\n\tsl2 len: %d; cap: %d\n", len(sl1), cap(sl1), len(sl2), cap(sl2))

	fmt.Printf("\n6 итерация (sl1 append 1) \"РАЗЛОМ\":\n")
	sl1 = append(sl1, 1)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, 1, 2, 1] len = 6; cap = 10; (new array)
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 1, 2]    len = 4; cap = 4;

	fmt.Printf("\n7 итерация (sl2 append 2) \"РАЗЛОМ\":\n")
	sl2 = append(sl2, 2)
	fmt.Printf("\tsl1: %v\n", sl1) // [0, 0, 1, 1, 2, 1] len = 6; cap = 10;
	fmt.Printf("\tsl2: %v\n", sl2) //    [0, 1, 1, 2, 2] len = 5; cap = 8;  (new array)

	fmt.Printf("\nПосле \"РАЗЛОМА\":\n\tsl1 len: %d; cap: %d\n\tsl2 len: %d; cap: %d\n", len(sl1), cap(sl1), len(sl2), cap(sl2))
}
