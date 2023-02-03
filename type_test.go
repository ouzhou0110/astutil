package ozastutil

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)


func TestToStruct(t *testing.T) {
	fst, f := InitEnv("./test_demo/type_demo.go")
	ast.Print(fst, f)
}

func TestAddKVToStruct(t *testing.T) {
	// 加载测试文件
	fset, f := InitEnv("./test_demo/type_demo.go")

	// 测试 type struct
	//type EmptyStruct struct {
	//}
	res := AddKVToStruct(f,"EmptyStruct", "Name", "string")
	assert.True(t, res)
	res = AddKVToStruct(f,"EmptyStruct", "age", "int")
	assert.True(t, res)
	res = AddKVToStruct(f,"EmptyStruct", "", "TypeStruct")
	assert.True(t, res)
	// 结果为：
	//type EmptyStruct struct {
	//	Name string
	//	age int
	//  TypeStruct
	//}

	PrintResult(fset, f)
}

func TestAddFuncToInterface(t *testing.T) {
	fset, f := InitEnv("./test_demo/type_demo.go")

	// 添加普通函数
	params := &AstFunc{
		Name: "test0",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "...*AstKv"},
		},
		Results: []AstKv{
			{Key: "ret1", Value: "string"},
		},
	}
	ret := AddFuncToInterface(f, "EmptyInf",params)
	assert.True(t, ret)

	params = &AstFunc{
		Name: "test01",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "...*AstKv"},
		},
		Results: []AstKv{
			{Key: "", Value: "string"},
		},
	}
	ret = AddFuncToInterface(f, "TestInf",params)
	assert.True(t, ret)
	PrintResult(fset, f)

}