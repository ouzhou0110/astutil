package test_demo

func Func1(a string) (ret string) {
	var b = a
	ret = b
	return ret
}

type Stu struct {
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
