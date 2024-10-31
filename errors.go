package gotainer

import "fmt"

type PrefetchArgumentError struct {
	ParentTypeName string
	DependencyName string
}

func (e *PrefetchArgumentError) Error() string {
	return fmt.Sprintf("unable to prefetch type %s which is a dependency for %s, please check ordering to ensure this type is registered before the type requiring it.", e.DependencyName, e.ParentTypeName)
}

func NewPrefetchArgumentError(parentTypeName, dependencyName string) *PrefetchArgumentError {
	return &PrefetchArgumentError{
		ParentTypeName: parentTypeName,
		DependencyName: dependencyName,
	}
}

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
