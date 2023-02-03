package ozastutil

import (
	"go/ast"
	"go/token"
	"strconv"
)

// AddImport 添加包名
//  参数 name 可以为空
func AddImport(fset *token.FileSet, f *ast.File, name, path string) bool {
	if !canImport(f, name, path) {
		return false
	}

	// 注册一个新增import实例
	newImport := &ast.ImportSpec{
		Path: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(path),
		},
	}

	// 避免为空时，前方出现一个空格
	if name != "" {
		newImport.Name = &ast.Ident{Name: name}
	}

	// 判断原文件中是否存在一个或多个import关键字
	impDecl := &ast.GenDecl{
		Tok: token.IMPORT,
	}
	lastImport := -1
	lastIndex := -1
	for i, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.IMPORT {
			lastImport = i
			lastIndex = len(gen.Specs) - 1
			impDecl = gen
		}
	}

	// 原文件没有import关键字
	if lastImport == -1 {
		impDecl.TokPos = f.Package
		file := fset.File(f.Package)
		pkgLine := file.Line(f.Package)
		for _, c := range f.Comments {
			if file.Line(c.Pos()) > pkgLine {
				break
			}
			// +2 for a blank line
			impDecl.TokPos = c.End() + 2
		}
		f.Decls = append(f.Decls, nil)
		copy(f.Decls[lastImport+2:], f.Decls[lastImport+1:])
		f.Decls[lastImport+1] = impDecl
	} else {
		impDecl.TokPos = f.Decls[lastImport].End()
	}

	// 将import实例注入decl的specs空间中
	insertAt := lastIndex + 1
	impDecl.Specs = append(impDecl.Specs, nil)
	copy(impDecl.Specs[insertAt+1:], impDecl.Specs[insertAt:])
	impDecl.Specs[insertAt] = newImport
	return true
}

func canImport(f *ast.File, name, path string) bool {
	if path == "" {
		return false
	}
	for _, s := range f.Imports {
		if importName(s) == name && importPath(s) == path {
			return false
		}
	}
	return true
}

// importName returns the name of s,
// or "" if the import is not named.
func importName(s *ast.ImportSpec) string {
	if s.Name == nil {
		return ""
	}
	return s.Name.Name
}

// importPath returns the unquoted import path of s,
// or "" if the path is not properly quoted.
func importPath(s *ast.ImportSpec) string {
	t, err := strconv.Unquote(s.Path.Value)
	if err != nil {
		return ""
	}
	return t
}
