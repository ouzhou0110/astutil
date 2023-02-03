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
