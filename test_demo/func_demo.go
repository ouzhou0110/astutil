package test_demo

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Func1(a string) (ret string) {
	var b = a
	ret = b
	return ret
}

func t () {

}

func tStruct() {
	var te = &Stu{}
	stu := &Stu{}
	fmt.Println(te,stu)
}


type Stu struct {
}
func (s *Stu) TT(gin *gin.Context) {
	return
}

func (s *Stu) PStuFunc(a string) (ret string) {
	var b = a
	ret = b
	return ret
}

func (s Stu) StuFunc(a string) (ret string) {
	var b = a
	ret = b
	return ret
}

