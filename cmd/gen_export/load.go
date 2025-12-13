// load.go
package main

import (
	"go/ast"
	"go/types"
	"log"

	"golang.org/x/tools/go/packages"
)

// load1 加载指定导入路径的包，并返回加载成功的包对象
// 如果加载失败或者存在错误，则会调用 log.Fatalf 终止程序执行
func load1(cfg *packages.Config, ip string) *packages.Package {
	pkgs, err := packages.Load(cfg, ip)
	if err != nil || len(pkgs) != 1 || pkgs[0].Types == nil {
		log.Fatalf("load %s failed: %v", ip, err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		log.Fatalf("load %s has errors", ip)
	}
	p := pkgs[0]
	if p.PkgPath == "" {
		// 回退机制：p.ID 通常类似于导入路径
		p.PkgPath = p.ID
	}
	if p.PkgPath == "" {
		log.Fatalf("load %s: empty PkgPath and ID", ip)
	}
	return p
}

// kind 判断给定的对象是否为类型、函数、变量或常量之一
// 如果是上述四种类型之一则返回 "ok"，否则返回空字符串
func kind(obj types.Object) string {
	switch obj.(type) {
	case *types.TypeName, *types.Func, *types.Var, *types.Const:
		return "ok"
	default:
		return ""
	}
}

// checkNameConflicts 检查多个相关包之间是否存在名称冲突的情况
// 参数 cfg 是用于加载包的配置信息
// 参数 mod 是模块的基本路径
// 参数 relPkgs 是相对于模块基本路径的相关包列表
func checkNameConflicts(cfg *packages.Config, mod string, relPkgs []string) {
	seen := map[string]string{} // 导出名称 -> 包路径 的映射关系

	for _, rel := range relPkgs {
		p := load1(cfg, mod+"/"+rel)
		ign := ignores(p) // 获取当前包中的忽略项

		scope := p.Types.Scope()
		for _, name := range scope.Names() {
			// 跳过未导出的名称以及被标记为忽略的名称
			if !ast.IsExported(name) || ign[name] {
				continue
			}
			// 只检查类型、函数、变量或常量等有效对象
			if kind(scope.Lookup(name)) == "" {
				continue
			}
			// 若发现同名但来自不同包的对象，则报告命名冲突错误
			if prev, ok := seen[name]; ok && prev != p.PkgPath {
				log.Fatalf("name conflict: %s in %s and %s (use //%s)", name, prev, p.PkgPath, ignoreTag)
			}
			seen[name] = p.PkgPath
		}
	}
}
