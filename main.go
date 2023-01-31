package main

import "fmt"

type Test1 struct{}

func (Test1) Show() {
	fmt.Println(123)
}

func (t Test1) ShowByAnother() {
	t.Show()
}

type Test2 struct {
	Test1
}

func (t Test2) Show() {
	//t.Test1.Show()
	fmt.Println(321)
}

func main() {
	t := Test2{}
	t.Show()
	fmt.Println("===")

	t.ShowByAnother()
}
