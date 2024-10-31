package gotainer

import (
	"fmt"
	"reflect"
	"unsafe"
)

type unsafeCtor func() (unsafe.Pointer, error)

type Container struct {
	singletonCtors map[string]unsafe.Pointer
	transientCtors map[string]unsafe.Pointer
	singletons     map[string]unsafe.Pointer
}

func NewContainer() *Container {
	return &Container{
		singletonCtors: make(map[string]unsafe.Pointer),
		transientCtors: make(map[string]unsafe.Pointer),
		singletons:     make(map[string]unsafe.Pointer),
	}
}

func ResolveInterface[T any](container *Container) (T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	var defaultVal T
	ptr, err := resolveNoReflect(container, t.Name())
	if err != nil {
		return defaultVal, err
	}
	herp := (*T)(ptr)
	fmt.Println(herp)
	return *(*T)(ptr), nil
}

func MustResolveInterface[T any](container *Container) T {
	res, err := ResolveInterface[T](container)
	if err != nil {
		panic(err)
	}
	return res
}

func Resolve[T any](container *Container) (*T, error) {
	t := reflect.TypeOf((*T)(nil)).Elem()
	ptr, err := resolveNoReflect(container, t.Name())
	if err != nil {
		return nil, err
	}
	return (*T)(ptr), nil
}

func MustResolve[T any](container *Container) *T {
	res, err := Resolve[T](container)
	if err != nil {
		panic(err)
	}
	return res
}

func MustRegisterTransient[T any, Fn any](container *Container, ctor Fn) {
	err := RegisterTransient[T, Fn](container, ctor)
	if err != nil {
		panic(err)
	}
}

func MustRegisterSingleton[T any, Fn any](container *Container, ctor Fn) {
	err := RegisterSingleton[T, Fn](container, ctor)
	if err != nil {
		panic(err)
	}
}

func RegisterTransient[T any, Fn any](container *Container, ctor Fn) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	fnType := reflect.TypeOf(ctor)
	err := testFn(container, t, fnType)
	if err != nil {
		return err
	}

	wrappedCtor := wrapCtor[T, Fn](container, fnType, &ctor, t.Kind() == reflect.Interface)
	container.transientCtors[t.Name()] = unsafe.Pointer(&wrappedCtor)
	return nil
}

func RegisterSingleton[T any, Fn any](container *Container, ctor Fn) error {
	t := reflect.TypeOf((*T)(nil)).Elem()
	fnType := reflect.TypeOf(ctor)
	err := testFn(container, t, fnType)
	if err != nil {
		return err
	}

	wrappedCtor := wrapCtor[T, Fn](container, fnType, &ctor, t.Kind() == reflect.Interface)
	wrappedSingletonCtor := wrapSingletonCtor[T](container, t.Name(), wrappedCtor)
	container.singletonCtors[t.Name()] = unsafe.Pointer(&wrappedSingletonCtor)
	return nil
}

func testFn(container *Container, contentType reflect.Type, fnType reflect.Type) error {
	if fnType.Kind() != reflect.Func {
		return NewConstructorMismatchError("ctor must be a function")
	}

	if fnType.NumOut() != 2 {
		return NewConstructorMismatchError("ctor must have 2 return values")
	}

	if fnType.Out(1).Name() != "error" {
		return NewConstructorMismatchError("ctor must return an error as the second return value")
	}

	firstOut := fnType.Out(0)
	if firstOut.Kind() != reflect.Ptr && firstOut.Kind() != reflect.Interface {
		return NewConstructorMismatchError("ctor must return a pointer or interface to the type it is constructing")
	}

	if firstOut.Kind() != reflect.Ptr &&
		firstOut.Kind() == reflect.Interface &&
		firstOut.Name() != contentType.Name() {
		return NewConstructorMismatchError("ctor must return an interface to the type it is constructing when registering an interface")
	} else if firstOut.Kind() != reflect.Interface && firstOut.Elem().Name() != contentType.Name() {
		return NewConstructorMismatchError("ctor must return a pointer to the type it is constructing when registering a struct type")
	}

	return findPrefetchErrors(container, fnType, contentType.Name())
}

func findPrefetchErrors(container *Container, funcType reflect.Type, typeName string) error {
	inputCount := funcType.NumIn()
	if inputCount == 0 {
		return nil
	}
	for i := 0; i < inputCount; i++ {
		input := funcType.In(i)
		name := input.Name()
		if input.Kind() == reflect.Ptr {
			name = input.Elem().Name()
		}

		_, isSingleton := container.singletonCtors[name]
		_, isTransient := container.transientCtors[name]
		if !isSingleton && !isTransient {
			return NewPrefetchArgumentError(typeName, name)
		}
	}

	return nil
}

func wrapCtor[T any, Fn any](container *Container, funcType reflect.Type, ctor *Fn, isInterface bool) unsafeCtor {
	return func() (unsafe.Pointer, error) {
		inputCount := funcType.NumIn()
		if inputCount == 0 {
			if isInterface {
				return unpackInterfaceCall(reflect.ValueOf(*ctor).Call([]reflect.Value{}))
			}
			return unpackStructCall(reflect.ValueOf(*ctor).Call([]reflect.Value{}))
		}

		vals := make([]reflect.Value, inputCount)
		for i := 0; i < inputCount; i++ {
			input := funcType.In(i)
			name := input.Name()
			if input.Kind() == reflect.Ptr {
				name = input.Elem().Name()
			}
			resolvedInput, err := resolveNoReflect(container, name)
			if err != nil {
				return nil, err
			}
			vals[i] = reflect.NewAt(input, unsafe.Pointer(&resolvedInput)).Elem()
		}

		if isInterface {
			return unpackInterfaceCall(reflect.ValueOf(*ctor).Call(vals))
		}
		return unpackStructCall(reflect.ValueOf(*ctor).Call(vals))
	}
}

func unpackInterfaceCall(args []reflect.Value) (unsafe.Pointer, error) {
	if !args[1].IsNil() {
		return nil, args[1].Interface().(error)
	}
	int := args[0].Interface()
	return unsafe.Pointer(&int), nil
}

func unpackStructCall(args []reflect.Value) (unsafe.Pointer, error) {
	var argOne unsafe.Pointer
	var argTwo error
	if !args[0].IsNil() {
		argOne = args[0].UnsafePointer()
	}
	if !args[1].IsNil() {
		argTwo = args[1].Interface().(error)
	}
	return argOne, argTwo
}

func resolveNoReflect(container *Container, name string) (unsafe.Pointer, error) {
	singletonCtor, ok := container.singletonCtors[name]
	if ok {
		return (*(*unsafeCtor)(singletonCtor))()
	}

	// no need to check for ok here, if it's not a singleton it must be a transient since we prefetch check
	transientCtor, _ := container.transientCtors[name]
	return (*(*unsafeCtor)(transientCtor))()
}

func wrapSingletonCtor[T any](container *Container, name string, ctor unsafeCtor) unsafeCtor {
	return func() (unsafe.Pointer, error) {
		singleton, ok := container.singletons[name]
		if ok {
			return singleton, nil
		}

		constructed, err := ctor()
		if err != nil {
			return nil, err
		}
		container.singletons[name] = unsafe.Pointer(constructed)
		return constructed, nil
	}
}
