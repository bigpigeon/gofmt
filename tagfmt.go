/*
 * Copyright 2020 bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

package main

import (
	"go/ast"
	"strings"
)

type tagFormatter struct {
	Err error
}

func (s *tagFormatter) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.StructType:
		if n.Fields != nil {
			var longestList []int
			var groupStart int
			for i, field := range n.Fields.List {
				if field.Tag != nil {
					_, keyValues, err := ParseTag(field.Tag.Value)
					if err != nil {
						s.Err = err
						return nil
					}
					for i, kv := range keyValues {
						if len(longestList) <= i {
							longestList = append(longestList, 0)
						}
						longestList[i] = max(len(kv.KeyValue), longestList[i])
					}
				} else {
					fieldsTagFormat(n.Fields.List[i:groupStart], longestList)
					groupStart = i + 1
					longestList = nil
				}
			}
			fieldsTagFormat(n.Fields.List[groupStart:], longestList)
		}
	}
	return s
}

func fieldsTagFormat(fields []*ast.Field, longestList []int) {
	for _, f := range fields {
		quote, keyValues, err := ParseTag(f.Tag.Value)
		if err != nil {
			// must be nil error
			panic(err)
		}
		var keyValueRaw []string
		for i, kv := range keyValues {
			keyValueRaw = append(keyValueRaw, kv.KeyValue+strings.Repeat(" ", longestList[i]-len(kv.KeyValue)))
		}

		f.Tag.Value = quote + strings.Join(keyValueRaw, " ") + quote
		f.Tag.ValuePos = 0
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func tagFmt(f *ast.File) error {
	s := &tagFormatter{}
	ast.Walk(s, f)
	return s.Err
}
