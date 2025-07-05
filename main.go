package main

import "fmt"

func main() {

	//fmt.Println("hellow world!")
	var length uint64 = 101
	h := New(length)
	// fmt.Print("slots:\n", h.slots)
	// foo-111 -> 3
	// foo-1 -> 3
	// h.Insert()
	h.Insert("foo-1", 40)
	fmt.Print("slots:\n", h.slots)
	fmt.Printf("item: %v\n", h.slots[3])

}
