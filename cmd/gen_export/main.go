package main

import (
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/imports"
)

func GenExportCmdMain() {
	main()
}

func main() {
	// 定义命令行参数
	mod := flag.String("module", "", "module path, e.g. github.com/dsx137/gg-kit")
	out := flag.String("out", ".", "output dir")

	var pkgsArg pkgList
	flag.Var(&pkgsArg, "pkg", "internal package rel path; can be repeated; supports glob like internal/*")

	// 解析命令行参数
	flag.Parse()
	if *mod == "" {
		log.Fatal("missing -module") // 如果没有提供 module 参数，则退出程序
	}

	// 展开包列表（处理通配符）
	pkgs := expandPkgs([]string(pkgsArg))
	if len(pkgs) == 0 {
		log.Fatal("missing -pkg (use -pkg a -pkg b ...; glob supported)") // 如果未指定任何包，则退出程序
	}

	// 配置 Go 包加载选项
	cfg := &packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedCompiledGoFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedSyntax |
			packages.NeedModule,
		Fset: token.NewFileSet(),
	}

	// 1) 快速检查输入包之间的导出名称冲突
	checkNameConflicts(cfg, *mod, pkgs)

	// 2) 为每个包生成代码
	for _, rel := range pkgs {
		p := load1(cfg, *mod+"/"+rel) // 加载单个包
		short := base(rel)            // 获取包的基本名称
		code := gen(p, short)         // 生成代码

		// 构造输出文件路径
		outFile := filepath.Join(*out, fmt.Sprintf("%s_export.gen.go", short))

		// 格式化生成的代码
		formatted, err := imports.Process(outFile, code, nil)
		if err != nil {
			// 如果格式化失败，打印原始代码以便调试
			fmt.Fprintln(os.Stderr, string(code))
			log.Fatal(err)
		}

		// 将格式化后的代码写入文件
		if err := os.WriteFile(outFile, formatted, 0644); err != nil {
			log.Fatal(err)
		}
	}
}
