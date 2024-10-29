package gotainer

import (
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

func Resolve[T any](container *Container) (*T, error) {
	var hardType T
	t := reflect.TypeOf(hardType)
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
	var hardType T
	t := reflect.TypeOf(hardType)

	fnType := reflect.TypeOf(ctor)
	err := testFn(t, fnType, false)
	if err != nil {
		return err
	}

	wrappedCtor := wrapCtor[T, Fn](container, fnType, &ctor)
	container.transientCtors[t.Name()] = unsafe.Pointer(&wrappedCtor)
	return nil
}

func RegisterSingleton[T any, Fn any](container *Container, ctor Fn) error {
	var hardType T
	t := reflect.TypeOf(hardType)

	fnType := reflect.TypeOf(ctor)
	err := testFn(t, fnType, false)
	if err != nil {
		return err
	}

	wrappedCtor := wrapCtor[T, Fn](container, fnType, &ctor)
	wrappedSingletonCtor := wrapSingletonCtor[T](container, t.Name(), wrappedCtor)
	container.singletonCtors[t.Name()] = unsafe.Pointer(&wrappedSingletonCtor)
	return nil
}

func testFn(contentType reflect.Type, fnType reflect.Type, useInterface bool) error {
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
	if useInterface {
		if firstOut.Kind() != reflect.Interface || firstOut.Name() != contentType.Name() {
			return NewConstructorMismatchError("ctor must return an interface to the type it is constructing when registering an interface")
		}
	} else {
		if firstOut.Kind() != reflect.Ptr || firstOut.Elem().Name() != contentType.Name() {
			return NewConstructorMismatchError("ctor must return a pointer to the type it is constructing when registering a struct type")
		}
	}

	return nil
}

func wrapCtor[T any, Fn any](container *Container, funcType reflect.Type, ctor *Fn) unsafeCtor {
	return func() (unsafe.Pointer, error) {
		inputCount := funcType.NumIn()
		if inputCount == 0 {
			return unpackCall(reflect.ValueOf(*ctor).Call([]reflect.Value{}))
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

		return unpackCall(reflect.ValueOf(*ctor).Call(vals))
	}
}

func unpackCall(args []reflect.Value) (unsafe.Pointer, error) {
	if len(args) != 2 {
		return nil, NewConstructorMismatchError("ctor must have 2 return values")
	}
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

	transientCtor, ok := container.transientCtors[name]
	if ok {
		return (*(*unsafeCtor)(transientCtor))()
	}
	return nil, NewTypeNotFoundError(name)
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
