package main

import (
	"fmt"
	"go/types"
	"strconv"
	"strings"
)

// typeWriter 为包装器签名渲染类型。
//
// 关键行为：
//   - 如果一个*命名*类型来自被包装的包，并且它的名称在 localTypes 中，
//     则将其渲染为裸名称"Name"（因为 ggkit 将其重新导出为 `type Name = alias.Name`）。
//   - 其他包渲染为"pkgname.Type"，goimports 可以通过导入来修复。
type typeWriter struct {
	wrappedPkgPath string
	localTypes     map[string]bool
}

func newTypeWriter(wrappedPkgPath string, localTypes map[string]bool) *typeWriter {
	return &typeWriter{
		wrappedPkgPath: wrappedPkgPath,
		localTypes:     localTypes,
	}
}

func (w *typeWriter) Type(t types.Type) string { return w.render(t) }

func (w *typeWriter) render(t types.Type) string {
	switch tt := t.(type) {
	case *types.Basic:
		return tt.Name()

	case *types.TypeParam:
		return tt.Obj().Name()

	case *types.Named:
		obj := tt.Obj()
		name := obj.Name()

		// 如果这个命名类型来自被包装的包并且我们重新导出它，
		// 则渲染为裸名称"Name"（无限定符）。
		if obj.Pkg() != nil && obj.Pkg().Path() == w.wrappedPkgPath && w.localTypes[name] {
			return name + w.typeArgs(tt.TypeArgs())
		}

		// 对于其他包，使用包名限定（goimports 可以解析）。
		if obj.Pkg() != nil {
			return obj.Pkg().Name() + "." + name + w.typeArgs(tt.TypeArgs())
		}
		return name + w.typeArgs(tt.TypeArgs())

	case *types.Pointer:
		return "*" + w.render(tt.Elem())

	case *types.Slice:
		return "[]" + w.render(tt.Elem())

	case *types.Array:
		return fmt.Sprintf("[%d]%s", tt.Len(), w.render(tt.Elem()))

	case *types.Map:
		return "map[" + w.render(tt.Key()) + "]" + w.render(tt.Elem())

	case *types.Chan:
		switch tt.Dir() {
		case types.SendOnly:
			return "chan<- " + w.render(tt.Elem())
		case types.RecvOnly:
			return "<-chan " + w.render(tt.Elem())
		default:
			return "chan " + w.render(tt.Elem())
		}

	case *types.Signature:
		return "func" + w.signature(tt)

	case *types.Struct:
		// 对于导出的 API 很少见；实现以保证完整性。
		var fields []string
		for i := 0; i < tt.NumFields(); i++ {
			f := tt.Field(i)
			part := ""
			if f.Embedded() {
				part = w.render(f.Type())
			} else {
				part = f.Name() + " " + w.render(f.Type())
			}
			if tag := tt.Tag(i); tag != "" {
				part += " " + strconv.Quote(tag)
			}
			fields = append(fields, part)
		}
		return "struct{ " + strings.Join(fields, "; ") + " }"

	case *types.Interface:
		// 保持可读性，不完全等同于 go/types 打印器。
		if tt.NumExplicitMethods() == 0 && tt.NumEmbeddeds() == 0 {
			return "interface{}"
		}
		var parts []string
		for i := 0; i < tt.NumEmbeddeds(); i++ {
			parts = append(parts, w.render(tt.EmbeddedType(i)))
		}
		for i := 0; i < tt.NumExplicitMethods(); i++ {
			m := tt.ExplicitMethod(i)
			sig, _ := m.Type().(*types.Signature)
			parts = append(parts, m.Name()+w.signature(sig))
		}
		return "interface{ " + strings.Join(parts, "; ") + " }"

	default:
		// 对于上面未处理的特殊类型，回退到 go/types 打印器。
		return types.TypeString(t, func(p *types.Package) string {
			if p == nil {
				return ""
			}
			return p.Name()
		})
	}
}

func (w *typeWriter) typeArgs(list *types.TypeList) string {
	if list == nil || list.Len() == 0 {
		return ""
	}
	var args []string
	for i := 0; i < list.Len(); i++ {
		args = append(args, w.render(list.At(i)))
	}
	return "[" + strings.Join(args, ", ") + "]"
}

func (w *typeWriter) signature(sig *types.Signature) string {
	if sig == nil {
		return "()"
	}

	// 参数
	ps := sig.Params()
	var params []string
	for i := 0; i < ps.Len(); i++ {
		v := ps.At(i)
		n := v.Name()
		if n == "" || n == "_" {
			n = fmt.Sprintf("_p%d", i)
		}
		if sig.Variadic() && i == ps.Len()-1 {
			sl := v.Type().(*types.Slice)
			params = append(params, fmt.Sprintf("%s ...%s", n, w.render(sl.Elem())))
		} else {
			params = append(params, fmt.Sprintf("%s %s", n, w.render(v.Type())))
		}
	}

	// 返回值
	rs := sig.Results()
	if rs.Len() == 0 {
		return "(" + strings.Join(params, ", ") + ")"
	}
	var results []string
	for i := 0; i < rs.Len(); i++ {
		v := rs.At(i)
		t := w.render(v.Type())
		if v.Name() != "" {
			results = append(results, v.Name()+" "+t)
		} else {
			results = append(results, t)
		}
	}
	return "(" + strings.Join(params, ", ") + ") (" + strings.Join(results, ", ") + ")"
}
