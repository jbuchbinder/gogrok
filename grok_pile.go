// GOGROK
// https://github.com/jbuchbinder/gogrok

package grok

// #cgo LDFLAGS: -lgrok
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "grok.h"
*/
import "C"

import (
	"os"
)

// A grok pile is an easy way to have multiple patterns together so
// that you can try to match against each one.
// The API provided should be similar to the normal Grok
// interface, but you can compile multiple patterns and match will
// try each one until a match is found.
type Pile struct {
	Groks        []Grok
	Patterns     map[string]string
	PatternFiles []string
}

func (this *Pile) AddPattern(name, pattern string) {
	this.Patterns[name] = pattern
}

func (this *Pile) AddPatternsFromFile(path string) {
	if _, err := os.Stat(path); err != os.ErrNotExist {
		this.PatternFiles = append(this.PatternFiles, path)
	}
}

func (this *Pile) Compile(pattern string) {
	grok := NewGrok()
	for k, v := range this.Patterns {
		grok.AddPattern(k, v)
	}
	for _, v := range this.PatternFiles {
		grok.AddPatternsFromFile(v)
	}
	grok.Compile(pattern)
	this.Groks = append(this.Groks, grok)
}

func (this *Pile) Match(text string) (g Grok, m GrokMatch) {
	for _, v := range this.Groks {
		match, _ := v.Match(text)
		if match.Matched {
			return v, match
		}
	}
	return
}
