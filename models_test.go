package gotainer_test

import "errors"

type SimpleStruct struct {
	data int
}

func NewSimpleStruct() (*SimpleStruct, error) {
	return &SimpleStruct{data: 1}, nil
}

func BadCtorForSimpleStructNonPtr() (SimpleStruct, error) {
	return SimpleStruct{data: 2}, nil
}

func BadCtorForSimpleStructNoErr() *SimpleStruct {
	return &SimpleStruct{data: 2}
}

func BadCtorForSimpleStructNonErr() (*SimpleStruct, int) {
	return &SimpleStruct{data: 3}, 3
}

type TierZeroType struct {
	ref *TierOneType
}

func NewTierZeroType(ref *TierOneType) (*TierZeroType, error) {
	return &TierZeroType{ref: ref}, nil
}

type TierOneType struct {
	ref  *TierTwoTypeOne
	ref2 *TierTwoTypeTwo
}

func NewTierOneType(ref *TierTwoTypeOne, ref2 *TierTwoTypeTwo) (*TierOneType, error) {
	return &TierOneType{ref: ref, ref2: ref2}, nil
}

type TierTwoTypeOne struct {
	data int
}

func NewTierTwoTypeOne() (*TierTwoTypeOne, error) {
	return &TierTwoTypeOne{data: 1}, nil
}

type TierTwoTypeTwo struct {
	data string
}

func NewTierTwoTypeTwo() (*TierTwoTypeTwo, error) {
	return &TierTwoTypeTwo{data: "two"}, nil
}

var TierTwoTypeTwoError = errors.New("tier two type two error")

func NewErroringTierTwoTypeTwo() (*TierTwoTypeTwo, error) {
	return nil, TierTwoTypeTwoError
}

type InterfaceableType struct {
}

func (i *InterfaceableType) DoThing() {
}

type InterfaceType interface {
	DoThing()
}

func NewInterfaceableType() (InterfaceType, error) {
	return &InterfaceableType{}, nil
}
