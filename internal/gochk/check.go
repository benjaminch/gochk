package gochk

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	red    = "\033[1;31m%s\033[0m"
	yellow = "\033[1;33m%s\033[0m"
	green  = "\033[1;32m%s\033[0m"
	teal   = "\033[1;36m%s\033[0m"
)

type dependency struct {
	filepath     string
	currentLayer int
	path         string
	index        int
}

// Check makes sure the direction of dependency is correct
func Check(cfg Config) {
	errorDeps := make([]dependency, 0, 0)
	err := filepath.Walk(cfg.TargetPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if include(cfg.Ignore, path) { // todo
			if info.IsDir() {
				print(yellow, "[Ignored]  "+path)
				return filepath.SkipDir
			}
			print(yellow, "[Ignored]  "+path)
			return nil
		}
		if info.IsDir() || !strings.Contains(info.Name(), ".go") {
			return nil
		}
		tempDeps := checkDependency(cfg.DependencyOrders, path)
		errorDeps = append(errorDeps, tempDeps...)
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	if len(errorDeps) > 0 {
		for _, d := range errorDeps {
			print(red, "[Error]    "+d.filepath+" imports "+d.path)
			print(red, "           \""+cfg.DependencyOrders[d.currentLayer]+"\" depends on \""+cfg.DependencyOrders[d.index]+"\"")
		}
	}
}

func checkDependency(dependencies []string, path string) []dependency {
	currentLayer := search(dependencies, path)
	importLayers := retrieveLayers(dependencies, path, currentLayer)

	if len(importLayers) == 0 {
		print(teal, "[None]     "+path)
		return nil
	}
	redDeps := make([]dependency, 0, len(importLayers))

	for _, d := range importLayers {
		if d.index < currentLayer {
			redDeps = append(redDeps, d)
			continue
		}
	}
	if len(redDeps) > 0 {
		return redDeps
	}
	print(green, "[Verified] "+path)
	return nil
}

func retrieveLayers(dependencies []string, path string, currentLayer int) []dependency {
	filepath, _ := filepath.Abs(path)
	imports := readImports(filepath)
	layers := make([]dependency, 0, len(imports))

	for _, v := range imports {
		l := search(dependencies, v)
		if l != -1 {
			layers = append(layers, dependency{
				filepath:     path,
				currentLayer: currentLayer,
				path:         v,
				index:        l,
			})
		}
	}
	return layers
}

func search(strs []string, elm string) int {
	for i, v := range strs {
		if strings.Contains(elm, v) {
			return i
		}
	}
	return -1
}

func include(strs []string, elm string) bool {
	for _, v := range strs {
		if strings.Contains(elm, v) {
			return true
		}
	}
	return false
}

func print(color string, message string) {
	fmt.Printf(color, message)
	fmt.Println()
}