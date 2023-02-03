package ozastutil

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddImport(t *testing.T) {
	// 加载测试文件
	fset, f := InitEnv("./test_demo/import_demo.go")

	ret := AddImport(fset, f, "demo22", "demo1")
	assert.True(t, ret)

	ret = AddImport(fset, f, "", "fmt")
	assert.False(t, ret)

	ret = AddImport(fset, f, "", "")
	assert.False(t, ret)

	ret = AddImport(fset, f, "tt", "")
	assert.False(t, ret)

	ret = AddImport(fset, f, "", "demo1")
	assert.True(t, ret)

	// 结果为：
	// import (
	//	"fmt"
	//	demo22 "demo1"
	//	"demo1"
	// )

	// 显示结果，没有写入文件的
	PrintResult(fset,f)
}
