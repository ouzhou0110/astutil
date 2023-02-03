package test_demo

type EmptyStruct struct {

}
type TypeStruct struct {
	Name string
}


type EmptyInf interface {

}

type TestInf interface {
	Demo(a string) error
	Tests(a ...string) error
}