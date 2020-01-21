/*
 * Copyright 2020 bigpigeon. All rights reserved.
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

package main

import (
	"go/ast"
	"sort"
	"strings"
)

type tagFormatter struct {
	Err error
}

func (s *tagFormatter) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.StructType:
		if n.Fields != nil {

			for _, field := range n.Fields.List {

				if field.Tag != nil {

					quote, keyValues, err := ParseTag(field.Tag.Value)
					if err != nil {
						s.Err = err
						return nil
					}
					sort.Slice(keyValues, func(i, j int) bool {
						return keyValues[i].Key < keyValues[j].Key
					})
					var keyValueRows []string
					for _, kv := range keyValues {
						keyValueRows = append(keyValueRows, kv.KeyValue)
					}

					field.Tag.Value = quote + strings.Join(keyValueRows, " ") + quote

				}
			}
		}
	}
	return s
}

func tagFmt(f *ast.File) error {
	s := &tagFormatter{}
	ast.Walk(s, f)
	return s.Err
}
