//go:build mage

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/magefile/mage/sh"
)

var Default = Gen
var modulePath = "github.com/dsx137/gg-kit"
var outputDir = "./pkg/ggkit/"

func Gen() error {
	Clean()
	fmt.Println("Generating ggkit facade...")

	sh.Run("go", "mod", "tidy")

	if err := sh.Run(
		"go", "run", "./cmd/gen_export",
		"-module", modulePath,
		"-out", outputDir,
		"-srcPkg", "internal/*",
	); err != nil {
		return err
	}

	return nil
}

func Clean() error {
	fmt.Println("Cleaning...")
	matches, err := filepath.Glob(filepath.Join(outputDir, "*_export.gen.go"))
	if err != nil {
		return err
	}
	for _, f := range matches {
		_ = os.Remove(f)
	}
	return nil
}
