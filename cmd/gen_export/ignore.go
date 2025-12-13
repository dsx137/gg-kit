package main

import (
	"go/ast"
	"strings"

	"golang.org/x/tools/go/packages"
)

const ignoreTag = "ggkit:ignore"

// ignores 返回应被忽略的导出标识符，
// 基于包含 "ggkit:ignore" 的文档注释。
func ignores(p *packages.Package) map[string]bool {
	m := map[string]bool{}
	for _, f := range p.Syntax {
		ast.Inspect(f, func(n ast.Node) bool {
			switch d := n.(type) {
			case *ast.FuncDecl:
				if d.Name != nil && ast.IsExported(d.Name.Name) && hasTag(d.Doc) {
					m[d.Name.Name] = true
				}
			case *ast.GenDecl:
				for _, sp := range d.Specs {
					switch s := sp.(type) {
					case *ast.TypeSpec:
						if s.Name != nil && ast.IsExported(s.Name.Name) && (hasTag(d.Doc) || hasTag(s.Doc)) {
							m[s.Name.Name] = true
						}
					case *ast.ValueSpec:
						if hasTag(d.Doc) || hasTag(s.Doc) {
							for _, id := range s.Names {
								if id != nil && ast.IsExported(id.Name) {
									m[id.Name] = true
								}
							}
						}
					}
				}
			}
			return true
		})
	}
	return m
}

func hasTag(cg *ast.CommentGroup) bool {
	if cg == nil {
		return false
	}
	for _, c := range cg.List {
		if strings.Contains(c.Text, ignoreTag) {
			return true
		}
	}
	return false
}
