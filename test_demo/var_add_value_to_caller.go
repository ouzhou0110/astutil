package test_demo

type Struct1 struct {
	Name string
}

func NewSet(...interface{}) Struct1 {
	return Struct1{}
}

var tt = NewSet(12,"cc")

var xx = &Struct1{
	Name: "str",
}