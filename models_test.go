package gotainer_test

type SimpleStruct struct {
	data int
}

func NewSimpleStruct() (*SimpleStruct, error) {
	return &SimpleStruct{data: 1}, nil
}

func BadCtorForSimpleStruct() *SimpleStruct {
	return &SimpleStruct{data: 2}
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
