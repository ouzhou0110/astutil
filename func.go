package ozastutil

import (
	"go/ast"
	"go/token"
	"strings"
)

type AstKv struct {
	Key   string
	Value string
}

type AstCallExpr struct {
	FunName string
	FunSel  string
	Args    []string
}

type AstFunc struct {
	Name    string
	Params  []AstKv
	Results []AstKv
	Return  []string
	Recv    *AstKv
}

// AddFunc 新增一个函数，目前只支持如下格式：
//
// 参数：
// &AstFunc{
//		Name: "test",
//		Params: []AstKv{
//			{Key: "p1", Value: "string"},
//			{Key: "p2", Value: "*AstKv"},
//		},
//		Results: []AstKv{
//			{Key: "ret", Value: "string"}, // results 这里的可以省略
//		},
//		Return: []string{"ret"},
//	}
//
// 结果为：
// func test(p1 string, p2 *AstKv) (ret string) {
// return ret
// }
//
// 在上诉AstFunc中添加如下 &AstFunc{Recv:{Key:"s", Value:"*Stu"}} 的结果为：
//
// func (s Stu) test1(p1 string, p2 *AstKv) (ret string) {
//	return ret
// }
//
// 在上诉AstFunc中添加如下 &AstFunc{Recv:{Key:"s", Value:"*Stu"}} 的结果为：
//
// func (s *Stu) test2(p1 string, p2 *AstKv) (ret string) {
//	return ret
// }
func AddFunc(f *ast.File, params *AstFunc) bool {
	if params.Name == "" {
		return false
	}
	f.Decls = append(f.Decls, getFuncDecl(params))
	return true
}

// AddParamToFunc 给函数添加参数
//
// 测试数据：func t () {}
//
// 执行：AddParamToFunc(f, "t", "demo", "demo.IDemo")
//
// 结果为：func t(demo demo.IDemo,) {}
func AddParamToFunc(f *ast.File, funcName, paramName, paramType string) bool {
	if funcName == "" || paramName == "" || paramType == "" {
		return false
	}
	for _, decl := range f.Decls {
		switch fd := decl.(type) {
		case *ast.FuncDecl:
			if fd.Name.Name != funcName {
				continue
			}
			if fd.Type.Params.List == nil {
				fd.Type.Params.List = make([]*ast.Field, 0)
			}
			// 判断变量名称是否存在
			for _, field := range fd.Type.Params.List {
				if field.Names[0].Name == paramName {
					return false
				}
			}
			fd.Type.Params.List = append(fd.Type.Params.List, getField(paramName, paramType))
			return true
		}
	}
	return false
}

// AddKVToFuncUnaryStruct 为函数的struct指针变量添加key value数据。目前没有做变量重复判断
//
// 测试数据：
//
// func tStruct() {
//	var te = &Stu{}
//	stu := &Stu{}
//	fmt.Println(te,stu)
// }
//
// 执行：
//
// AddKVToFuncUnaryStruct(f, "tStruct", "te", "key", "value")
//
// AddKVToFuncUnaryStruct(f, "tStruct", "stu", "key", "value")
//
// 结果为：
//
// func tStruct() {
//	var te = &Stu{key: value}
//	stu := &Stu{key: value}
//	fmt.Println(te, stu)
// }
func AddKVToFuncUnaryStruct(f *ast.File, funcName, varName, key, value string) bool {
	if funcName == "" || varName == "" || key == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		switch fd := decl.(type) {
		case *ast.FuncDecl:
			if fd.Name.Name != funcName {
				continue
			}
			if fd.Body.List == nil {
				return false
			}
			for _, stmt := range fd.Body.List {
				switch s := stmt.(type) {
				case *ast.DeclStmt:
					// 处理 var te = &Stu{} 结构
					switch g := s.Decl.(type) {
					case *ast.GenDecl:
						vs := g.Specs[0].(*ast.ValueSpec)
						if vs.Names[0].Name != varName {
							continue
						}
						switch vs.Values[0].(type) {
						case *ast.UnaryExpr:
							vsVal := vs.Values[0].(*ast.UnaryExpr)
							return addKVToUnaryExpr(vsVal, key, value)
						}
					}
				case *ast.AssignStmt:
					// 处理 stu := &Stu{}
					if s.Lhs[0].(*ast.Ident).Name != varName {
						continue
					}
					switch s.Rhs[0].(type) {
					case *ast.UnaryExpr:
						vsVal := s.Rhs[0].(*ast.UnaryExpr)
						return addKVToUnaryExpr(vsVal, key, value)
					}
				}
			}
		}
	}
	return false
}

// AddVarToFunc 添加变量到函数中
//
// 支持格式有：
// var xx string
//
// var xx = "xxx"
//
// xx := ccc 有问题，暂不修复
//
// 不支持复杂的定义
func AddVarToFunc(f *ast.File, funcName, varName, value, afterVar, tag string) bool {
	if funcName == "" || varName == "" || value == "" {
		return false
	}
	for _, decl := range f.Decls {
		switch fd := decl.(type) {
		case *ast.FuncDecl:
			if fd.Name.Name != funcName {
				continue
			}
			if fd.Body.List == nil {
				fd.Body.List = make([]ast.Stmt, 0)
			}
			afterIndex := -1
			afterLine := token.NoPos
			if afterVar != "" && len(fd.Body.List) > 0 {
				for k, stmt := range fd.Body.List {
					switch s := stmt.(type) {
					// 处理 var te = &Stu{} 结构
					case *ast.DeclStmt:
						switch g := s.Decl.(type) {
						case *ast.GenDecl:
							if g.Specs[0].(*ast.ValueSpec).Names[0].Name == afterVar {
								afterIndex = k
								afterLine = g.Specs[0].(*ast.ValueSpec).Names[0].NamePos
								break
							}
						}
					// 处理 stu := &Stu{}
					case *ast.AssignStmt:
						switch v := s.Lhs[0].(type) {
						case *ast.Ident:
							if v.Name == afterVar {
								afterIndex = k
								afterLine = v.NamePos
								break
							}
						case *ast.SelectorExpr:
							if v.X.(*ast.Ident).Name == afterVar {
								afterIndex = k
								afterLine = v.X.(*ast.Ident).NamePos
								break
							}
						}
					}
				}
			}
			if afterVar == "" && len(fd.Body.List) > 0 {
				switch fd.Body.List[len(fd.Body.List) - 1].(type) {
				case *ast.ReturnStmt:
					afterIndex = len(fd.Body.List) - 2
				}
			}
			// 生产数据
			var newVar ast.Stmt
			switch tag {
			case "var":
				newVar = &ast.DeclStmt{Decl: getVar(varName, value)}
			case "define":
				newVar = &ast.DeclStmt{Decl: getDefineVar(varName, value)}
			case "assign":
				newVar = getAssignVar(varName, value, afterLine)
			}
			if afterIndex == -1 {
				fd.Body.List = append(fd.Body.List, newVar)
			} else {
				insertAt := afterIndex + 1
				fd.Body.List = append(fd.Body.List, nil)
				copy(fd.Body.List[insertAt+1:], fd.Body.List[insertAt:])
				fd.Body.List[insertAt] = newVar
			}
			return true
		}
	}
	return false
}

// AddCallBlockToFunc 添加一个如下所示的代码块到函数中
//
// {
//		rdemo.POST("", demo.Create)
//		rdemo.POST("/", demo.Create)
//	}
func AddCallBlockToFunc(f *ast.File, funcName string, data []AstCallExpr, afterVar string) bool {
	if funcName == "" || len(data) == 0 {
		return false
	}
	for _, decl := range f.Decls {
		switch fd := decl.(type) {
		case *ast.FuncDecl:
			if fd.Name.Name != funcName {
				continue
			}
			if fd.Body.List == nil {
				fd.Body.List = make([]ast.Stmt, 0)
			}
			afterIndex := -1
			if afterVar != "" && len(fd.Body.List) > 0 {
				for k, stmt := range fd.Body.List {
					switch s := stmt.(type) {
					// 处理 var te = &Stu{} 结构
					case *ast.DeclStmt:
						switch g := s.Decl.(type) {
						case *ast.GenDecl:
							if g.Specs[0].(*ast.ValueSpec).Names[0].Name == afterVar {
								afterIndex = k
								break
							}
						}
					// 处理 stu := &Stu{}
					case *ast.AssignStmt:
						switch v := s.Lhs[0].(type) {
						case *ast.Ident:
							if v.Name == afterVar {
								afterIndex = k
								break
							}
						case *ast.SelectorExpr:
							if v.X.(*ast.Ident).Name == afterVar {
								afterIndex = k
								break
							}
						}
					}
				}
			}
			// 生产数据
			newVar := &ast.BlockStmt{
				List: make([]ast.Stmt, 0),
			}
			for _, datum := range data {
				args := make([]ast.Expr, 0)
				if len(datum.Args) > 0 {
					for _, arg := range datum.Args {
						args = append(args, &ast.BasicLit{Value: arg})
					}
				}
				d := &ast.ExprStmt{
					X: &ast.CallExpr{
						Fun: &ast.SelectorExpr{
							X:   &ast.Ident{Name: datum.FunName},
							Sel: &ast.Ident{Name: datum.FunSel},
						},
						Args: args,
					},
				}
				newVar.List = append(newVar.List, d)
			}
			if afterIndex == -1 {
				fd.Body.List = append(fd.Body.List, newVar)
			} else {
				insertAt := afterIndex + 1
				fd.Body.List = append(fd.Body.List, nil)
				copy(fd.Body.List[insertAt+1:], fd.Body.List[insertAt:])
				fd.Body.List[insertAt] = newVar
			}
			return true
		}
	}
	return false
}

// GetLastVarFormFunc 获取函数中最后一个变量名称
func GetLastVarFormFunc(f *ast.File, funcName string) string {
	if funcName == "" {
		return ""
	}
	retName := ""
	for _, decl := range f.Decls {
		switch fd := decl.(type) {
		case *ast.FuncDecl:
			if fd.Name.Name != funcName {
				continue
			}
			if fd.Body.List == nil {
				return ""
			}
			for _, stmt := range fd.Body.List {
				switch s := stmt.(type) {
				// 处理 var te = &Stu{} 结构
				case *ast.DeclStmt:
					switch g := s.Decl.(type) {
					case *ast.GenDecl:
						retName = g.Specs[0].(*ast.ValueSpec).Names[0].Name
					}
				// 处理 stu := &Stu{}
				case *ast.AssignStmt:
					switch v := s.Lhs[0].(type) {
					case *ast.Ident:
						retName = v.Name
					case *ast.SelectorExpr:
						retName = v.X.(*ast.Ident).Name
					}
				}
			}
		}
	}
	return retName
}

func addKVToUnaryExpr(ue *ast.UnaryExpr, key, value string) bool {
	switch cpl := ue.X.(type) {
	case *ast.CompositeLit:
		if cpl.Elts == nil {
			cpl.Elts = []ast.Expr{}
		}
		// 判断参数是否存在 -- 不完善，暂时屏蔽
		//for _, elt := range cpl.Elts {
		//	if elt.(*ast.KeyValueExpr).Key.(*ast.Ident).Name == key {
		//		return false
		//	}
		//}
		cpl.Elts = append(cpl.Elts, getKVExpr(key, value))
		return true
	}
	return false
}

func getFuncDecl(params *AstFunc) *ast.FuncDecl {
	newFunc := &ast.FuncDecl{Name: ast.NewIdent(params.Name)}
	if params.Params != nil {
		if newFunc.Type == nil {
			newFunc.Type = &ast.FuncType{
				Params: &ast.FieldList{
					List: make([]*ast.Field, 0),
				},
			}
		}
		for _, param := range params.Params {
			newFunc.Type.Params.List = append(newFunc.Type.Params.List, getField(param.Key, param.Value))
		}
	}

	if params.Results != nil {
		if newFunc.Type == nil {
			newFunc.Type = &ast.FuncType{
				Results: &ast.FieldList{
					List: make([]*ast.Field, 0),
				},
			}
		}
		if newFunc.Type.Results == nil {
			newFunc.Type.Results = &ast.FieldList{
				List: make([]*ast.Field, 0),
			}
		}
		for _, param := range params.Results {
			newFunc.Type.Results.List = append(newFunc.Type.Results.List, getField(param.Key, param.Value))
		}
	}

	if params.Return != nil {
		if newFunc.Body == nil {
			newFunc.Body = &ast.BlockStmt{
				List: make([]ast.Stmt, 0),
			}
		}
		newFunc.Body.List = append(newFunc.Body.List, getReturnStmt(params.Return))
	}

	if params.Recv != nil {
		newFunc.Recv = &ast.FieldList{
			List: []*ast.Field{},
		}
		field := &ast.Field{
			Names: []*ast.Ident{ast.NewIdent(params.Recv.Key)},
			Type:  ast.NewIdent(params.Recv.Value),
		}
		if strings.HasPrefix(params.Recv.Value, "*") {
			field.Type = &ast.StarExpr{
				X: ast.NewIdent(strings.TrimPrefix(params.Recv.Value, "*")),
			}
		}
		newFunc.Recv.List = append(newFunc.Recv.List, field)
	}
	return newFunc
}

func getReturnStmt(rtns []string) *ast.ReturnStmt {
	ret := &ast.ReturnStmt{Return: token.NoPos}
	if len(rtns) == 0 {
		return ret
	}
	ret.Results = make([]ast.Expr, 0)
	for _, rtn := range rtns {
		ret.Results = append(ret.Results, ast.NewIdent(rtn))
	}
	return ret
}
