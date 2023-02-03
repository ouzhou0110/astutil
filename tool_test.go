package ozastutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWriteToFile(t *testing.T) {
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
	err := WriteToFile(fst, f, "./test_demo/func_demo_1.go")
	assert.NoError(t, err)
}
