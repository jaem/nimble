package main

import "fmt"

type base struct {
	Name string
}

func (base *base) BaseMethod() string {
	return base.Name
}

type subclass struct {
	*base
}

func main() {

	value := &subclass{&base{Name: "subclass"}}
	fmt.Println(value.Name)

}
