package ozastutil

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestParseVariable(t *testing.T) {
	// 加载测试文件
	_, f := InitEnv("./test_demo/var_demo.go")

	res := ParseVariable(f, "cc")
	assert.True(t, res)
}

func TestAddVariable(t *testing.T) {
	// 加载测试文件
	fst, f := InitEnv("./test_demo/var_add_demo.go")

	res := DefineVarAfterVar(f, "cc", "*ast.GenExpr", "test")
	assert.True(t, res)
	// 结果为：var cc *ast.GenExpr

	res = AddVarAfterVar(f, "ccb", "&ast.GenExpr", "test")
	assert.True(t, res)
	// 结果为：var ccb = &ast.GenExpr

	PrintResult(fst, f)

}

func TestAddVar(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_demo.go")
	ast.Print(fst, f)
}

func TestAddValueToMap(t *testing.T) {
	// 加载测试文件
	fst, f := InitEnv("./test_demo/var_add_value_to_map.go")

	// 测试数据：var mapStr = map[string]string{"cc": "cc"}
	res := AddValueToMap(f, "mapStr", AddQuote("key"), AddQuote("test"))
	assert.True(t, res)
	// 结果为：var mapStr = map[string]string{"cc": "cc", "key": "test"}

	// 测试数据：var mapInt = map[int]string{1: "cc"}
	res = AddValueToMap(f, "mapInt", "3", AddQuote("test"))
	assert.True(t, res)
	// 结果为：var mapInt = map[int]string{1: "cc", 3: "test"}

	// 测试数据：var mapInf = map[string]interface{}{"cc": 1}
	res = AddValueToMap(f, "mapInf", AddQuote("hello"), "&aaa")
	assert.True(t, res)
	// 结果为：var mapInf = map[string]interface{}{"cc": 1, "hello": &aaa}

	PrintResult(fst, f)

}

func TestToMap(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_map.go")
	ast.Print(fst, f)
}

func TestAddValueToCaller(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_caller.go")

	// 测试数据：var tt = NewSet(12,"cc")

	// 添加字符串："hello"
	res := AddValueToCaller(f, "tt", AddQuote("hello"))
	assert.True(t, res)

	// 添加变量名称：hello
	res = AddValueToCaller(f, "tt", "hello")
	assert.True(t, res)
	PrintResult(fst, f)
	// 结果为：var tt = NewSet(12, "cc", "hello", hello)
}

func TestAddKVToUnaryStruct(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_caller.go")

	// 测试数据：
	//var xx = &Struct1{
	//	Name: "str",
	//}

	// 添加成员属性：Key: "value"
	res := AddKVToUnaryStruct(f, "xx", "Key", AddQuote("value"))
	assert.True(t, res)
	PrintResult(fst, f)
	// 结果为：var xx = &Struct1{
	//	Name: "str", Key: "value",
	//}
}

func TestToCaller(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_caller.go")
	ast.Print(fst, f)
}

func TestAddValueToSlice(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_slice.go")

	// 添加字符串："hello"
	res := AddValueToSlice(f, "sliceStr", AddQuote("hello"))
	assert.True(t, res)

	// 添加变量名称：hello
	res = AddValueToSlice(f, "sliceStr", "hello")
	assert.True(t, res)
	PrintResult(fst, f)
	// 结果为：var sliceStr = []string{"1", "hello", hello}
}

func TestToSlice(t *testing.T) {
	fst, f := InitEnv("./test_demo/var_add_value_to_slice.go")
	ast.Print(fst, f)
}
