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

func TestContainer_RegisterSingletonWithBadCtorNonPtr_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[SimpleStruct](c, BadCtorForSimpleStructNonPtr)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering singleton with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterTransientWithBadCtorNonPtr_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, BadCtorForSimpleStructNonPtr)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering transient with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterTransientWithBadCtorNonFunc_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, 3)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering transient with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterSingletonWithBadCtorNonFunc_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[SimpleStruct](c, 3)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering singleton with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterTransientWithBadCtorSignatureNonErr_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, BadCtorForSimpleStructNonErr)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering transient with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterSingletonWithBadCtorSignatureNonError_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[SimpleStruct](c, BadCtorForSimpleStructNonErr)
	ctorErr := &gotainer.ConstructorMismatchError{}
	if err == nil {
		t.Error("expected error when registering singleton with bad ctor signature")
	}
	if !errors.As(err, &ctorErr) {
		t.Error("expected error to be ConstructorMismatchError")
	}
}

func TestContainer_RegisterTransientWithBadCtorSignature_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, BadCtorForSimpleStructNoErr)
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
	err := gotainer.RegisterSingleton[SimpleStruct](c, BadCtorForSimpleStructNoErr)
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

func TestContainer_ResolveErrConstructorInComplexSingletonType_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[TierTwoTypeTwo](c, NewErroringTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}

	err = gotainer.RegisterSingleton[TierTwoTypeOne](c, NewTierTwoTypeOne)
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

	_, err = gotainer.Resolve[TierZeroType](c)
	if err == nil {
		t.Error("expected error when resolving with downstream constructor error")
		return
	}

	if !errors.Is(err, TierTwoTypeTwoError) {
		t.Error("expected error to be TierTwoTypeTwoError")
		return
	}
}

func TestContainer_ResolveErrConstructorInComplexTransientType_ReturnsError(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[TierTwoTypeTwo](c, NewErroringTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}

	err = gotainer.RegisterTransient[TierTwoTypeOne](c, NewTierTwoTypeOne)
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
	if err == nil {
		t.Error("expected error when resolving with downstream constructor error")
		return
	}

	if !errors.Is(err, TierTwoTypeTwoError) {
		t.Error("expected error to be TierTwoTypeTwoError")
		return
	}
}

func TestContainer_MustRegisterTransientWithBadCtorSignature_Panics(t *testing.T) {
	c := gotainer.NewContainer()
	defer func() { _ = recover() }()

	gotainer.MustRegisterTransient[SimpleStruct](c, BadCtorForSimpleStructNoErr)

	t.Error("expected panic")
}

func TestContainer_MustRegisterSingletonWithBadCtorSignature_Panics(t *testing.T) {
	c := gotainer.NewContainer()
	defer func() { _ = recover() }()

	gotainer.MustRegisterSingleton[SimpleStruct](c, BadCtorForSimpleStructNoErr)

	t.Error("expected panic")
}

func TestContainer_MustResolveTypeNotFound_Panics(t *testing.T) {
	c := gotainer.NewContainer()
	defer func() { _ = recover() }()

	gotainer.MustResolve[SimpleStruct](c)

	t.Error("expected panic")
}

func TestContainer_MustResolveTypeWithCtorError_Panics(t *testing.T) {
	c := gotainer.NewContainer()
	defer func() { _ = recover() }()
	err := gotainer.RegisterSingleton[TierTwoTypeTwo](c, NewErroringTierTwoTypeTwo)
	if err != nil {
		t.Error(err)
		return
	}

	gotainer.MustResolve[TierTwoTypeTwo](c)

	t.Error("expected panic")
}

func TestContainer_MustResolveHappyPath_ResolvesAsExpected(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[SimpleStruct](c, NewSimpleStruct)
	if err != nil {
		t.Error(err)
		return
	}

	s := gotainer.MustResolve[SimpleStruct](c)

	if s == nil {
		t.Error("expected resolved type to not be nil")
	}
}

func TestContainer_RegisterSingletonWithTypesRegisteredOutOfOrder_FailsPrefetchCheck(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterSingleton[TierZeroType](c, NewTierZeroType)

	if err == nil {
		t.Error("expected error when resolving with types registered out of order, if this is not an error cycles can occur in resolution")
		return
	}

	prefetchErr := &gotainer.PrefetchArgumentError{}
	if !errors.As(err, &prefetchErr) {
		t.Error("expected error to be PrefetchArgumentError")
		return
	}

	if prefetchErr.ParentTypeName != "TierZeroType" {
		t.Error("expected error to be for TierZeroType was for ", prefetchErr.ParentTypeName)
		return
	}

	if prefetchErr.DependencyName != "TierOneType" {
		t.Error("expected error to be for TierOneType was for ", prefetchErr.DependencyName)
		return
	}
}

func TestContainer_RegisterTransientWithTypesRegisteredOutOfOrder_FailsPrefetchCheck(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[TierZeroType](c, NewTierZeroType)

	if err == nil {
		t.Error("expected error when resolving with types registered out of order, if this is not an error cycles can occur in resolution")
		return
	}

	prefetchErr := &gotainer.PrefetchArgumentError{}
	if !errors.As(err, &prefetchErr) {
		t.Errorf("expected error to be PrefetchArgumentError, got %v", err)
		return
	}

	if prefetchErr.ParentTypeName != "TierZeroType" {
		t.Errorf("expected error to be for TierZeroType was for %s", prefetchErr.ParentTypeName)
		return
	}

	if prefetchErr.DependencyName != "TierOneType" {
		t.Errorf("expected error to be for TierOneType was for %s", prefetchErr.DependencyName)
		return
	}
}

func TestContainer_ResolveTransientInterface_Resolves(t *testing.T) {
	c := gotainer.NewContainer()
	err := gotainer.RegisterTransient[InterfaceType](c, NewInterfaceableType)
	if err != nil {
		t.Error(err)
		return
	}

	resolved, err := gotainer.ResolveInterface[InterfaceType](c)
	if err != nil {
		t.Error(err)
		return
	}

	if resolved == nil {
		t.Error("expected resolved type to not be nil")
	}

	resolved.DoThing()
}
