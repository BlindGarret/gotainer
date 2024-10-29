package main

import (
	"fmt"

	"github.com/BlindGarret/gotainer"
)

var rawTypeId int = 0
var secondaryTypeId int = 0

type RawType struct {
	id int
}

func NewRawType() (*RawType, error) {
	rawTypeId++
	return &RawType{id: rawTypeId}, nil
}

type SecondaryType struct {
	id   int
	data *RawType
}

func NewSecondaryType(r *RawType) (*SecondaryType, error) {
	secondaryTypeId++
	return &SecondaryType{data: r, id: secondaryTypeId}, nil
}

func (s *SecondaryType) GetId() int {
	return s.data.id
}

func main() {
	// todo: maintain ordering. right now loops are possible. We need to check that our reqs exists before we allow registration.
	// todo: allow resolution of FNs which are unregistered. For example "fill this fn with services"

	container := gotainer.NewContainer()
	gotainer.MustRegisterTransient[SecondaryType](container, NewSecondaryType)
	gotainer.MustRegisterSingleton[RawType](container, NewRawType)

	secondaryTypeOne := gotainer.MustResolve[SecondaryType](container)
	fmt.Printf("SecondaryType id: %d has rawType id %d \n", secondaryTypeOne.id, secondaryTypeOne.GetId())

	secondaryTypeTwo := gotainer.MustResolve[SecondaryType](container)
	fmt.Printf("SecondaryType id: %d has rawType id %d \n", secondaryTypeTwo.id, secondaryTypeTwo.GetId())
}
