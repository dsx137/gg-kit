package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// pkgList 是字符串切片的别名，用于表示包列表
type pkgList []string

// String 返回 pkgList 的字符串表示形式，各元素用逗号连接
func (p *pkgList) String() string {
	return strings.Join(*p, ",")
}

// Set 将值添加到 pkgList 中
func (p *pkgList) Set(v string) error {
	*p = append(*p, v)
	return nil
}

// hasGlob 检查字符串是否包含通配符 (*, ?, [)
func hasGlob(s string) bool {
	return strings.ContainsAny(s, "*?[")
}

// expandPkgs 将 -pkg 输入参数（支持通配符模式）展开成一个唯一且排序后的列表
func expandPkgs(inputs []string) []string {
	seen := map[string]bool{} // 用于去重
	var out []string          // 输出结果

	// 内部函数 add 用来处理单个路径
	add := func(s string) {
		s = strings.TrimSpace(s) // 去除空格
		if s == "" {
			return
		}
		// 标准化为斜杠风格并去除尾部的 "/"
		s = filepath.ToSlash(s)
		s = strings.TrimSuffix(s, "/")
		if s == "" {
			return
		}
		// 如果未见过该路径，则加入结果列表
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}

	// 遍历所有输入项
	for _, raw := range inputs {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}

		// 如果是通配符模式
		if hasGlob(raw) {
			matches, err := filepath.Glob(raw) // 匹配文件路径
			if err != nil {
				log.Fatalf("bad -pkg glob %q: %v", raw, err)
			}
			// 对每个匹配项检查是否为目录
			for _, m := range matches {
				fi, err := os.Stat(m)
				if err == nil && !fi.IsDir() {
					continue // 跳过非目录项
				}
				add(m) // 添加目录项
			}
			continue
		}

		add(raw) // 直接添加非通配符路径
	}

	sort.Strings(out) // 排序输出结果
	return out
}

// base 获取相对路径中的最后一个目录或文件名
func base(rel string) string {
	rel = strings.TrimSuffix(rel, "/")
	if i := strings.LastIndex(rel, "/"); i >= 0 {
		return rel[i+1:]
	}
	return rel
}
