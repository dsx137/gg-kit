package main

import (
	"fmt"
	"go/types"
	"strings"
)

// genericPlan 决定：
//   - 哪些类型参数应该出现在包装器声明上
//   - 调用是否需要显式的类型参数，如果需要，则是哪些参数
type genericPlan struct {
	sig  *types.Signature
	used map[*types.TypeParam]bool
}

func planGenerics(sig *types.Signature) genericPlan {
	return genericPlan{
		sig:  sig,
		used: usedTypeParamsInSig(sig),
	}
}

// wrapperDecl 返回 "[T any, U ~int]" 或者 ""。
func (p genericPlan) wrapperDecl() string {
	return typeParamsDeclFiltered(p.sig, p.used)
}

// callArgs 返回 "" 或者 "[T, U]"。
// 只有当从调用参数无法推断时才发出显式类型参数。
func (p genericPlan) callArgs() string {
	if !needsExplicitTypeArgs(p.sig, p.used) {
		return ""
	}
	return typeArgsCallFiltered(p.sig, p.used)
}

// usedTypeParamsInSig 返回在参数/结果类型中出现的类型参数（来自 sig.TypeParams()）。
func usedTypeParamsInSig(sig *types.Signature) map[*types.TypeParam]bool {
	used := map[*types.TypeParam]bool{}
	markUsedInTypeTuple(sig.Params(), used)
	markUsedInTypeTuple(sig.Results(), used)
	return used
}

// markUsedInTypeTuple 遍历元组中的每个变量并标记其使用的类型参数
func markUsedInTypeTuple(tup *types.Tuple, used map[*types.TypeParam]bool) {
	if tup == nil {
		return
	}
	for i := 0; i < tup.Len(); i++ {
		markUsedInType(tup.At(i).Type(), used)
	}
}

// markUsedInType 在给定类型中递归查找并标记使用的类型参数
func markUsedInType(t types.Type, used map[*types.TypeParam]bool) {
	switch tt := t.(type) {
	case *types.TypeParam:
		// 如果是类型参数本身，直接标记为已使用
		used[tt] = true

	case *types.Named:
		// 对于具名类型，检查其类型参数
		if ta := tt.TypeArgs(); ta != nil {
			for i := 0; i < ta.Len(); i++ {
				markUsedInType(ta.At(i), used)
			}
		}

	// 处理各种复合类型，继续深入查找其中的类型参数
	case *types.Pointer:
		markUsedInType(tt.Elem(), used)
	case *types.Slice:
		markUsedInType(tt.Elem(), used)
	case *types.Array:
		markUsedInType(tt.Elem(), used)
	case *types.Map:
		markUsedInType(tt.Key(), used)
		markUsedInType(tt.Elem(), used)
	case *types.Chan:
		markUsedInType(tt.Elem(), used)

	case *types.Struct:
		// 结构体：遍历所有字段的类型
		for i := 0; i < tt.NumFields(); i++ {
			markUsedInType(tt.Field(i).Type(), used)
		}

	case *types.Interface:
		// 接口：处理嵌入类型和显式方法
		for i := 0; i < tt.NumEmbeddeds(); i++ {
			markUsedInType(tt.EmbeddedType(i), used)
		}
		for i := 0; i < tt.NumExplicitMethods(); i++ {
			if ms, ok := tt.ExplicitMethod(i).Type().(*types.Signature); ok {
				markUsedInTypeTuple(ms.Params(), used)
				markUsedInTypeTuple(ms.Results(), used)
			}
		}

	case *types.Signature:
		// 函数签名：处理参数和返回值
		markUsedInTypeTuple(tt.Params(), used)
		markUsedInTypeTuple(tt.Results(), used)
	}
}

// typeParamsDeclFiltered 仅为使用的参数生成包装器类型参数声明，保持原始顺序。
func typeParamsDeclFiltered(sig *types.Signature, used map[*types.TypeParam]bool) string {
	tps := sig.TypeParams()
	if tps == nil || tps.Len() == 0 {
		return ""
	}
	var d []string
	for i := 0; i < tps.Len(); i++ {
		tp := tps.At(i)
		if !used[tp] {
			continue
		}
		// 获取约束条件的字符串表示
		c := types.TypeString(tp.Constraint(), func(p *types.Package) string {
			if p == nil {
				return ""
			}
			return p.Name()
		})
		d = append(d, fmt.Sprintf("%s %s", tp.Obj().Name(), c))
	}
	if len(d) == 0 {
		return ""
	}
	return "[" + strings.Join(d, ", ") + "]"
}

// needsExplicitTypeArgs 报告是否有任何使用的类型参数不能从调用参数推断出来
// （即它没有出现在任何参数类型中）。
func needsExplicitTypeArgs(sig *types.Signature, used map[*types.TypeParam]bool) bool {
	if used == nil || len(used) == 0 {
		return false
	}

	// 收集出现在参数类型中的类型参数。
	inParams := map[*types.TypeParam]bool{}
	markUsedInTypeTuple(sig.Params(), inParams)

	// 如果一个使用的TP不在参数中，那么它就不能从参数中推断出来。
	for tp := range used {
		if !inParams[tp] {
			return true
		}
	}
	return false
}

// typeArgsCallFiltered 产生类似 "[T, U]" 的调用位置类型参数列表
// 按照原始顺序，但仅针对 [used](file://f:\Go\gg-kit\cmd\gen\generics.go#L13-L13) 类型参数。
func typeArgsCallFiltered(sig *types.Signature, used map[*types.TypeParam]bool) string {
	tps := sig.TypeParams()
	if tps == nil || tps.Len() == 0 {
		return ""
	}
	var c []string
	for i := 0; i < tps.Len(); i++ {
		tp := tps.At(i)
		if !used[tp] {
			continue
		}
		c = append(c, tp.Obj().Name())
	}
	if len(c) == 0 {
		return ""
	}
	return "[" + strings.Join(c, ", ") + "]"
}
