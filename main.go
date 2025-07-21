package main

import (
	"fmt"
	"unsafe"
)

type node struct {
	key   nodeKey
	value any
	// isSoftDeleted bool
}

func main() {
	fmt.Printf("size: %d bytes", unsafe.Sizeof(node{}))
}

// 0-15, 16-32, 33-39

// 00X00  key
// 00X01  key
// ...
// 00X-15 key
// 00X-16 value
// 00X-17 value
// ....
// 00X-32 value
// 00X-33 isSoftDeleted
// 00X-34 padding (1 byte)
// 00X-35 padding (1 byte)
// ...
// 00X39  padding
