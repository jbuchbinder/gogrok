// GOGROK
// https://github.com/jbuchbinder/gogrok

package grok

// #cgo LDFLAGS: -lgrok
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "grok.h"

// Make it easier to handle in Go bindings
typedef grok_t * grok_obj;
typedef grok_match_t * grok_match_obj;

*/
import "C"

import (
	"errors"
	//"fmt"
	"unsafe"
)

const (
	GROK_OK                         = 0
	GROK_ERROR_FILE_NOT_ACCESSIBLE  = 1
	GROK_ERROR_PATTERN_NOT_FOUND    = 2
	GROK_ERROR_UNEXPECTED_READ_SIZE = 3
	GROK_ERROR_COMPILE_FAILED       = 4
	GROK_ERROR_UNINITIALIZED        = 5
	GROK_ERROR_PCRE_ERROR           = 6
	GROK_ERROR_NOMATCH              = 7
)

type Grok struct {
	Obj C.grok_obj
}

type GrokMatch struct {
	Matched        bool
	Subject        string
	Start          int
	End            int
	OriginalObject C.grok_match_t
}

func NewGrok() (obj Grok) {
	obj = Grok{}
	obj.Obj = C.grok_new()
	return obj
}

func (this *Grok) AddPattern(name, pattern string) (err error) {
	name_c := C.CString(name)
	defer C.free(unsafe.Pointer(name_c))
	pattern_c := C.CString(pattern)
	defer C.free(unsafe.Pointer(pattern_c))

	ret := C.grok_pattern_add(this.Obj, name_c, C.size_t(len(name)), pattern_c, C.size_t(len(pattern)))
	if ret != GROK_OK {
		err = errors.New("Failed to add pattern " + name)
	}
	return
}

func (this *Grok) AddPatternsFromFile(path string) (err error) {
	path_c := C.CString(path)
	defer C.free(unsafe.Pointer(path_c))

	ret := C.grok_patterns_import_from_file(this.Obj, path_c)
	if ret != GROK_OK {
		err = errors.New("Failed to add patterns from file " + path)
	}
	return
}

func (this *Grok) Compile(pattern string) (err error) {
	pattern_c := C.CString(pattern)
	defer C.free(unsafe.Pointer(pattern_c))

	ret := C.grok_compilen(this.Obj, pattern_c, C.int(len(pattern)))
	if ret != GROK_OK {
		err = errors.New("Compile failed")
	}
	return
}

func (this *Grok) Match(text string) (gm GrokMatch, err error) {
	gm = GrokMatch{}
	text_c := C.CString(text)
	defer C.free(unsafe.Pointer(text_c))

	var m C.grok_match_t
	ret := C.grok_execn(this.Obj, text_c, C.int(len(text)), &m)
	switch ret {
	case GROK_OK:
		{
			gm.Matched = true
			gm.Subject = C.GoString(m.subject)
			gm.Start = int(m.start)
			gm.End = int(m.end)
			gm.OriginalObject = m
		}
	case GROK_ERROR_NOMATCH:
		{
			gm.Matched = false
		}
	}
	err = errors.New("Unknown return from grok_execn")
	return
}
