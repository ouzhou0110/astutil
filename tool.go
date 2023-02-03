package ozastutil

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
)

func InitEnv(path string) (fset *token.FileSet, f *ast.File) {
	fset = token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		panic(err)
	}
	return
}

func PrintResult(fset *token.FileSet, f *ast.File) {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		panic(err)
	}
	fmt.Printf("%s", buf.Bytes())
}

// AddQuote 给字符串添加双引号
func AddQuote(str string) string {
	return addTag(str, "\"")
}

// AddSingleQuote 给字符串添加单引号
func AddSingleQuote(str string) string  {
	return addTag(str, "'")
}

func addTag(str,tag string) string {
	return fmt.Sprintf("%s%s%s", tag,str,tag)
}


func WriteToFile(fset *token.FileSet, f *ast.File, path string) error {
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, f); err != nil {
		return err
	}
	return ioutil.WriteFile(path, buf.Bytes(), 0777)
}