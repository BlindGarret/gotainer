package gotainer

import "fmt"

type ConstructorMismatchError struct {
	Reason string
}

func (e *ConstructorMismatchError) Error() string {
	return fmt.Sprintf("ctor must return a pointer or interface to the type it is constructing, along with an error: (T*, error) || (Ti, error): %s", e.Reason)
}

func NewConstructorMismatchError(reason string) *ConstructorMismatchError {
	return &ConstructorMismatchError{
		Reason: reason,
	}
}

type TypeNotFoundError struct {
	TypeName string
}

func (e *TypeNotFoundError) Error() string {
	return fmt.Sprintf("type not found in container: %s", e.TypeName)
}

func NewTypeNotFoundError(typeName string) *TypeNotFoundError {
	return &TypeNotFoundError{
		TypeName: typeName,
	}
}
