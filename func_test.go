package ozastutil

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestAddFunc1(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	ast.Print(fst, f)


}

func TestAddFunc(t *testing.T) {

	// 添加普通函数
	fst, f := InitEnv("./test_demo/func_demo.go")
	params := &AstFunc{
		Name: "test",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "*AstKv"},
		},
		Results: []AstKv{
			{Key: "ret", Value: "string"},
		},
		Return: []string{"ret"},
	}
	ret := AddFunc(f, params)
	assert.True(t, ret)

	// 添加普通函数
	params = &AstFunc{
		Name: "test0",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "*AstKv"},
		},
		Results: []AstKv{
			{Key: "", Value: "string"},
		},
		Return: []string{"ret"},
	}
	ret = AddFunc(f, params)
	assert.True(t, ret)

	// 为struct等添加函数
	params = &AstFunc{
		Name: "test1",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "*AstKv"},
		},
		Results: []AstKv{
			{Key: "ret", Value: "string"},
		},
		Return: []string{"ret"},
		Recv: &AstKv{
			Key: "s",
			Value: "Stu",
		},
	}
	ret = AddFunc(f, params)

	// 为struct等添加函数
	params = &AstFunc{
		Name: "test2",
		Params: []AstKv{
			{Key: "p1", Value: "string"},
			{Key: "p2", Value: "*AstKv"},
		},
		Results: []AstKv{
			{Key: "ret", Value: "string"},
		},
		Return: []string{"ret"},
		Recv: &AstKv{
			Key: "s",
			Value: "*Stu",
		},
	}
	ret = AddFunc(f, params)
	assert.True(t, ret)

	PrintResult(fst,f)
}

func TestAddParamToFunc(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	ret := AddParamToFunc(f, "t", "demo", "demo.IDemo")
	assert.True(t, ret)

	ret = AddParamToFunc(f, "t", "demo", "demo.IDemo")
	assert.False(t, ret)

	PrintResult(fst,f)
}

func TestAddKVToFuncUnaryStruct(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	ret := AddKVToFuncUnaryStruct(f, "tStruct", "te", "key", "value")
	assert.True(t, ret)

	ret = AddKVToFuncUnaryStruct(f, "tStruct", "stu", "key", "value")
	assert.True(t, ret)

	//ret = AddKVToFuncUnaryStruct(f, "tStruct", "stu", "key", "value")
	//assert.False(t, ret)

	PrintResult(fst,f)
}

func TestAddVarToFunc(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	//ret := AddVarToFunc(f, "t", "te", "key", "", "var")
	//assert.True(t, ret)
	//ret = AddVarToFunc(f, "t", "te1", "key", "", "define")
	//assert.True(t, ret)
	ret := AddVarToFunc(f, "t", "te2222", "key", "", "assign")
	assert.True(t, ret)
	PrintResult(fst,f)

}

func TestAddCallBlockToFunc(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	data := []AstCallExpr{
		{FunName: "rdemo",FunSel: "POST", Args: []string{AddQuote(""),"demo.Create"}},
		{FunName: "rdemo",FunSel: "POST", Args: []string{AddQuote("/"),"demo.Create"}},
	}
	ret := AddCallBlockToFunc(f, "t", data, "key")
	assert.True(t, ret)
	PrintResult(fst,f)
}

func TestGetLastVarFormFunc(t *testing.T) {
	fst, f := InitEnv("./test_demo/func_demo.go")
	ret := GetLastVarFormFunc(f, "StuFunc")
	assert.Equal(t, ret, "ret")
	PrintResult(fst,f)
}