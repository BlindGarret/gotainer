package main

import (
	"fmt"

	"github.com/BlindGarret/gotainer"
)

type SimpleStruct struct {
	data int
}

func (s *SimpleStruct) DoSomething() {
	fmt.Println("doing something")
}

func NewSimpleStruct() (SimpleInterface, error) {
	return &SimpleStruct{data: 1}, nil
}

type SimpleInterface interface {
	DoSomething()
}

func main() {
	// todo: maintain ordering. right now loops are possible. We need to check that our reqs exists before we allow registration.
	// todo: allow resolution of FNs which are unregistered. For example "fill this fn with services"
	c := gotainer.NewContainer()
	gotainer.MustRegisterTransient[SimpleInterface](c, NewSimpleStruct)

	service, err := gotainer.ResolveInterface[SimpleInterface](c)
	if err != nil {
		panic(err)
	}
	service.DoSomething()
}
