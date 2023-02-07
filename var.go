package ozastutil

import (
	"fmt"
	"go/ast"
	"go/token"
)

// DefineVarAfterVar 在某个变量后面新定义一个变量
//  在test变量后面定义一个：var genExpr *ast.GenExpr
//  使用例子：DefineVarAfterVar(f, "genExpr", "*ast.GenExpr", "test")
func DefineVarAfterVar(f *ast.File, name, kind, afterVar string) bool {
	if name == "" || kind == "" || isVarExist(f, name) != -1 {
		return false
	}
	insertAt := isVarExist(f, afterVar)
	if insertAt == -1 {
		return false
	}
	newVar := getDefineVar(name, kind)
	if newVar == nil {
		return false
	}
	// 插入指定位置
	insertDecls(f, insertAt, newVar)
	return true
}

// AddVarAfterVar 在某个变量后面新增一个变量
//  在test变量后面新增一个：var genExpr = &ast.GenExpr
//  使用例子：AddVarAfterVar(f, "genExpr", "&ast.GenExpr", "test")
func AddVarAfterVar(f *ast.File, name, value, afterVar string) bool {
	if name == "" || value == "" || isVarExist(f, name) != -1 {
		return false
	}
	insertAt := isVarExist(f, afterVar)
	if insertAt == -1 {
		return false
	}
	newVar := getVar(name, value)
	if newVar == nil {
		return false
	}
	// 插入指定位置
	insertDecls(f, insertAt, newVar)
	return true
}

// AddValueToMap 给map变量添加数据，没有做map类型校验
//  测试数据：var mapStr = map[string]string{"cc": "cc"}
//  执行：AddValueToMap(f, "mapStr", AddQuote("key"), AddQuote("test"))
//  结果为：var mapStr = map[string]string{"cc": "cc", "key": "test"}
//
//  测试数据：var mapInt = map[int]string{1: "cc"}
//  执行：AddValueToMap(f, "mapInt", "3", AddQuote("test"))
//  结果为：var mapInt = map[int]string{1: "cc", 3: "test"}
//
//  测试数据：var mapInf = map[string]interface{}{"cc": 1}
//  执行：AddValueToMap(f, "mapInf", AddQuote("hello"), "&aaa")
//  结果为：var mapInf = map[string]interface{}{"cc": 1, "hello": &aaa}
func AddValueToMap(f *ast.File, mapName, key, value string) bool {
	if mapName == "" || key == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			if vs.Names[0].Name == mapName {
				switch vsVal := vs.Values[0].(type) {
				case *ast.CompositeLit:
					if vsVal.Elts == nil {
						vsVal.Elts = []ast.Expr{}
					}
					vsVal.Elts = append(vsVal.Elts, getKVExpr(key, value))
					return true
				}
			}
		}
	}
	return false
}

// AddValueToCaller 给CallerAble变量添加数据，没有做类型校验
//  测试数据：
//  func NewSet(...interface{}) struct {
// 	return struct{}
//  }
//  var tt = NewSet(12,"cc")
//  执行：
//  AddValueToCaller(f, "tt", AddQuote("hello"))
//  AddValueToCaller(f, "tt", "hello")
//  结果为：
//  var tt = NewSet(12, "cc", "hello", hello)
func AddValueToCaller(f *ast.File, varName, value string) bool {
	if varName == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			if vs.Names[0].Name == varName {
				switch vsVal := vs.Values[0].(type) {
				case *ast.CallExpr:
					if vsVal.Args == nil {
						vsVal.Args = []ast.Expr{}
					}
					vsVal.Args = append(vsVal.Args, &ast.BasicLit{
						Kind:  token.STRING,
						Value: value,
						ValuePos: token.NoPos + 1,
					})
					return true
				}
			}
		}
	}
	return false
}

// AddValueToSlice 给slice变量添加数据，没有做类型校验
//  例子：var sliceStr = []string{}
//  执行：
//  AddValueToSlice(f, "sliceStr", AddQuote("hello"))
//  AddValueToSlice(f, "sliceStr", "hello")
//  结果为：
//  var sliceStr = []string{"hello", hello}
func AddValueToSlice(f *ast.File, varName, value string) bool {
	if varName == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			if vs.Names[0].Name == varName {
				switch vsVal := vs.Values[0].(type) {
				case *ast.CompositeLit:
					if vsVal.Elts == nil {
						vsVal.Elts = []ast.Expr{}
					}
					vsVal.Elts = append(vsVal.Elts, &ast.BasicLit{
						Kind:  token.STRING,
						Value: value,
						ValuePos: token.NoPos + 1,
					})
					return true
				}
			}
		}
	}
	return false
}

// AddKVToUnaryStruct 给结构体指针变量添加key,value
//  比如：var xx = &Struct1{
//	 Name: "str",
//  }
// 给上面添加一组key "value"，结果为：
//	var xx = &Struct1{
//	 Name: "str",
//   Key: "value"
//  }
// 当然，对应Struct的应该是：
//  type Struct1 struct {
//	 Name string
//   Key  string
//  }
func AddKVToUnaryStruct(f *ast.File, varName, key, value string) bool {
	if varName == "" || key == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			if vs.Names[0].Name == varName {
				switch vs.Values[0].(type) {
				case *ast.UnaryExpr:
				default:
					return false
				}
				vsVal := vs.Values[0].(*ast.UnaryExpr)
				switch cpl := vsVal.X.(type) {
				case *ast.CompositeLit:
					if cpl.Elts == nil {
						cpl.Elts = []ast.Expr{}
					}
					cpl.Elts = append(cpl.Elts, getKVExpr(key,value))
					return true
				}
			}
		}
	}
	return false
}

// getKVExpr 获取一个*ast.KeyValueExpr
func getKVExpr(key, value string) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key:   &ast.BasicLit{Kind: token.STRING, Value: key, ValuePos: token.NoPos + 1},
		Value: &ast.BasicLit{Kind: token.STRING, Value: value, ValuePos: token.NoPos + 1},
	}
}

// insertDecls 在指定位置插入decls
func insertDecls(f *ast.File, offset int, decl ast.Decl) {
	f.Decls = append(f.Decls, nil)
	copy(f.Decls[offset+2:], f.Decls[offset+1:])
	f.Decls[offset+1] = decl
}

// getDefineVar 获取一个用于生成<var str1 string> 格式的GenDecl，
func getDefineVar(name, kind string) *ast.GenDecl {
	if name == "" || kind == "" {
		return nil
	}
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent(name)},
				Type:  ast.NewIdent(kind),
			},
		},
	}
}

// getVar 获取一个用于生成<var str1 = "string"> 格式的GenDecl，
func getVar(name, value string) *ast.GenDecl {
	if name == "" || value == "" {
		return nil
	}
	return &ast.GenDecl{
		Tok: token.VAR,
		Specs: []ast.Spec{
			&ast.ValueSpec{
				Names: []*ast.Ident{ast.NewIdent(name)},
				Values: []ast.Expr{&ast.BasicLit{
					Kind:     token.STRING,
					ValuePos: token.NoPos + 1,
					Value:    value,
				}},
			},
		},
	}
}

// getAssignVar 获取一个用于生成<str1 := "string"> 格式的 AssignStmt，
func getAssignVar(name, value string,afterLine token.Pos) *ast.AssignStmt {
	if name == "" || value == "" {
		return nil
	}
	return &ast.AssignStmt{
		TokPos: afterLine + 1,
		Lhs: []ast.Expr{&ast.Ident{Name: name}},
		Rhs: []ast.Expr{&ast.BasicLit{Value: value}},
	}
}

// isVarExist 判断变量是否存在
func isVarExist(f *ast.File, name string) int {
	ret := -1
	if name == "" {
		return ret
	}
	for k, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			for _, varIdent := range vs.Names {
				if varIdent.Name == name {
					return k
				}
			}
		}
	}
	return ret
}

func ParseVariable(f *ast.File, varName string) bool {
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.VAR {
			vs := gen.Specs[0].(*ast.ValueSpec)
			for k, varIdent := range vs.Names {
				if varIdent.Name == varName {
					varValue := vs.Values[k]
					fmt.Printf("GenDecl=>%#v\n", gen)
					fmt.Printf("ValueSpec=>%#v\n", vs)
					fmt.Printf("Ident=>%#v\n", varIdent)
					fmt.Printf("Expr=>%#v\n", varValue)
				}
			}
		}
	}
	return true
}
