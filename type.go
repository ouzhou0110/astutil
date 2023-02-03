package ozastutil

import (
	"go/ast"
	"go/token"
)

// AddKVToStruct 为 type xx struct 添加成员属性
//  参数key允许为空
//  例子：
//  type EmptyStruct struct {
//  }
//  执行：
//  AddKVToStruct(f,"EmptyStruct", "age", "int")
//  AddKVToStruct(f,"EmptyStruct", "", "TypeStruct")
//  结果为：
//  type EmptyStruct struct {
// 	 age int
//   TypeStruct
//  }
func AddKVToStruct(f *ast.File, name, key, value string) bool {
	if name == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.TYPE {
			varSpec := gen.Specs[0].(*ast.TypeSpec)
			if varSpec.Name.Name == name {
				switch specType := varSpec.Type.(type) {
				case *ast.StructType:
					typeFields := specType.Fields
					if typeFields.List == nil {
						typeFields.List = []*ast.Field{}
					}
					typeFields.List = append(typeFields.List, getField(key, value))
					return true
				}
			}
		}
	}
	return false
}

// AddFuncToInterface 为interface添加函数
//
// 测试数据：
//
// type EmptyInf interface {
// }
//
// 传参：
//
// params := &AstFunc{
//		Name: "test0",
//		Params: []AstKv{
//			{Key: "p1", Value: "string"},
//			{Key: "p2", Value: "...*AstKv"},
//		},
//		Results: []AstKv{
//			{Key: "ret1", Value: "string"},
//		},
// }
//
// 结果为：
//
// type EmptyInf interface {
//	test0(p1 string, p2 ...*AstKv) (ret1 string)
// }
func AddFuncToInterface(f *ast.File, name string, params *AstFunc) bool {
	if name == "" || params == nil {
		return false
	}
	for _, decl := range f.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if ok && gen.Tok == token.TYPE {
			varSpec := gen.Specs[0].(*ast.TypeSpec)
			if varSpec.Name.Name == name {
				switch specType := varSpec.Type.(type) {
				case *ast.InterfaceType:
					typeFields := specType.Methods
					if typeFields.List == nil {
						typeFields.List = []*ast.Field{}
					}
					typeFields.List = append(typeFields.List, &ast.Field{
						Names: []*ast.Ident{ast.NewIdent(params.Name)},
						Type: getFuncType(params),
					})
					return true
				}
			}
		}
	}
	return false
}

func getField(key, value string) *ast.Field {
	ret := &ast.Field{
		Type: ast.NewIdent(value),
	}
	if key != "" {
		ret.Names = []*ast.Ident{ast.NewIdent(key)}
	}
	return ret
}

func getFuncType(params *AstFunc) *ast.FuncType {
	newFunc := &ast.FuncType{}
	if params.Params != nil {
		newFunc.Params = &ast.FieldList{
			List: make([]*ast.Field, 0),
		}
		for _, param := range params.Params {
			newFunc.Params.List = append(newFunc.Params.List, getField(param.Key, param.Value))
		}
	}

	if params.Results != nil {
		newFunc.Results = &ast.FieldList{
			List: make([]*ast.Field, 0),
		}
		for _, param := range params.Results {
			newFunc.Results.List = append(newFunc.Results.List, getField(param.Key, param.Value))
		}
	}
	return newFunc
}
