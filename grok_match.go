// GOGROK
// https://github.com/jbuchbinder/gogrok

package grok

// #cgo LDFLAGS: -lgrok
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "grok.h"

typedef grok_match_t * grok_match_obj;

*/
import "C"

import (
	//"fmt"
	"unsafe"
)

func GetCaptures() (ret map[string]string) {
	var gm C.grok_match_obj
	C.grok_match_walk_init(gm)

	var name_ptr *C.char
	defer C.free(unsafe.Pointer(name_ptr))
	var namelen_ptr C.int
	var data_ptr *C.char
	defer C.free(unsafe.Pointer(data_ptr))
	var datalen_ptr C.int

	ret = map[string]string{}

	for int(C.grok_match_walk_next(gm, &name_ptr, &namelen_ptr, &data_ptr, &datalen_ptr)) == GROK_OK {
		ret[C.GoString(name_ptr)[0:namelen_ptr]] = C.GoString(data_ptr)[0:datalen_ptr]
	}
	C.grok_match_walk_end(gm)

	return
}
