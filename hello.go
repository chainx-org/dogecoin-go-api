package main

/*
#cgo LDFLAGS: -L./lib -ldogecoin_dll
#include <stdlib.h>
#include "./lib/DogecoinHeader.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func main() {
	s := "Go say: Hello Rust"

	input := C.CString(s)
	defer C.free(unsafe.Pointer(input))
	o := C.generate_my_privkey_dogecoin(input, input)
	output := C.GoString(o)
	fmt.Printf("%s\n", output)
}
