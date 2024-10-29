package gotainer_test

import (
	"errors"
	"testing"

	"github.com/BlindGarret/gotainer"
)

func TestContainer_Constructed_NotNil(t *testing.T) {
	c := gotainer.NewContainer()

	if c == nil {
		t.Error("Container is nil")
	}
}

func TestContainer_RegisterTransientWithBadCtorSignature_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, BadCtorForSimpleStruct)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering transient with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterSingletonWithBadCtorSignature_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[SimpleStruct](c, BadCtorForSimpleStruct)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering singleton with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterSingletonSingleSimpleStruct_ResolvesDifferentTransients(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[SimpleStruct](c, NewSimpleStruct)
	if err != nil {
		t.Error(err)
		return
	}

	firstStruct, err := gotainer.Resolve[SimpleStruct](c)
	if err != nil {
		t.Error(err)
		return
	}
	secondStruct, err := gotainer.Resolve[SimpleStruct](c)
	if err != nil {
		t.Error(err)
		return
	}

	if firstStruct != secondStruct {
		t.Error("singletons when resolved should be the same reference")
	}
}

func TestContainer_RegisterTransientSingleSimpleStruct_ResolvesDifferentTransients(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, NewSimpleStruct)
	if err != nil {
		t.Error(err)
		return
	}

	firstStruct, err := gotainer.Resolve[SimpleStruct](c)
	if err != nil {
		t.Error(err)
		return
	}
	secondStruct, err := gotainer.Resolve[SimpleStruct](c)
	if err != nil {
		t.Error(err)
		return
	}

	if firstStruct == secondStruct {
		t.Error("transients when resolved should not be the same reference")
	}
}

func TestContainer_RegisterSingletonComplexObject_ResolvesSingletons(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[TierTwoTypeOne](c, NewTierTwoTypeOne)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterSingleton[TierTwoTypeTwo](c, NewTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterSingleton[TierOneType](c, NewTierOneType)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterSingleton[TierZeroType](c, NewTierZeroType)
	if err != nil {
		t.Error(err)
		return
	}

	resolvedOne, err := gotainer.Resolve[TierZeroType](c)
	if err != nil {
		t.Error(err)
		return
	}

	resolveTwo, err := gotainer.Resolve[TierZeroType](c)
	if err != nil {
		t.Error(err)
		return
	}

	if resolvedOne != resolveTwo {
		t.Error("singletons when resolved should be the same reference")
		return
	}

	if resolvedOne.ref != resolveTwo.ref {
		t.Error("singletons when resolved should be the same reference")
		return
	}

	if resolvedOne.ref.ref != resolveTwo.ref.ref {
		t.Error("singletons when resolved should be the same reference")
		return
	}

	if resolvedOne.ref.ref2 != resolveTwo.ref.ref2 {
		t.Error("singletons when resolved should be the same reference")
		return
	}
}

func TestContainer_RegisterTransientComplexObject_ResolvesTransients(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[TierTwoTypeOne](c, NewTierTwoTypeOne)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterTransient[TierTwoTypeTwo](c, NewTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterTransient[TierOneType](c, NewTierOneType)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterTransient[TierZeroType](c, NewTierZeroType)
	if err != nil {
		t.Error(err)
		return
	}

	resolvedOne, err := gotainer.Resolve[TierZeroType](c)
	if err != nil {
		t.Error(err)
		return
	}

	resolvedTwo, err := gotainer.Resolve[TierZeroType](c)
	if err != nil {
		t.Error(err)
		return
	}

	if resolvedOne == resolvedTwo {
		t.Error("transients when resolved should not be the same reference")
		return
	}

	if resolvedOne.ref == resolvedTwo.ref {
		t.Error("transients when resolved should not be the same reference")
		return
	}

	if resolvedOne.ref.ref == resolvedTwo.ref.ref {
		t.Error("transients when resolved should not be the same reference")
		return
	}

	if resolvedOne.ref.ref2 == resolvedTwo.ref.ref2 {
		t.Error("transients when resolved should not be the same reference")
		return
	}
}

func TestContainer_ResolveTypeNotFound_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	_, err := gotainer.Resolve[SimpleStruct](c)
	typeErr := &gotainer.TypeNotFoundError{}
	if err == nil {
		t.Error("expected error when resolving type not found")
		return
	}

	if !errors.As(err, &typeErr) {
		t.Error("expected error to be ConstructorMismatchError")
		return
	}

	if typeErr.TypeName != "SimpleStruct" {
		t.Error("expected error to be for SimpleStruct")
		return
	}
}

func TestContainer_ResolveTypeNotFoundInComplexType_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[TierTwoTypeTwo](c, NewTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterTransient[TierOneType](c, NewTierOneType)
	if err != nil {
		t.Error(err)
		return
	}
	err = gotainer.RegisterTransient[TierZeroType](c, NewTierZeroType)
	if err != nil {
		t.Error(err)
		return
	}

	_, err = gotainer.Resolve[TierZeroType](c)
	typeErr := &gotainer.TypeNotFoundError{}
	if err == nil {
		t.Error("expected error when resolving type not found")
		return
	}

	if !errors.As(err, &typeErr) {
		t.Error("expected error to be ConstructorMismatchError")
		return
	}

	if typeErr.TypeName != "TierTwoTypeOne" {
		t.Error("expected error to be for TierTwoTypeOne")
		return
	}
}
