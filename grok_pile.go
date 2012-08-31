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
