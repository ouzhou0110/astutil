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
			Type: ast.NewIdent(params.Recv.Value),
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
