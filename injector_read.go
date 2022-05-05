package injector

import (
	"fmt"
	"reflect"
)

// Inject tries to resolve a dependency by its type and optionally its name
// If the dependency is unknown, ErrorDepFactoryNotFound is returned
func Inject[T any](name ...string) (T, error) {
	argType := reflect.TypeOf((*T)(nil)).Elem()
	factory := depMap.GetOrDefault(reflectTypeKey(argType), nameOrDefault(name), nil)
	if factory == nil {
		return reflect.Zero(argType).Interface().(T), ErrorDepFactoryNotFound
	}

	dep, err := factory()
	if err != nil {
		return reflect.Zero(argType).Interface().(T), err
	}

	castDep, ok := dep.(T)
	if !ok {
		return reflect.Zero(argType).Interface().(T), ErrorDependencyTypeMismatch
	}
	return castDep, nil
}

// InjectT tries to resolve a dependency by its type and optionally its name
// If the dependency is unknown, ErrorDepFactoryNotFound is returned
func InjectT(depType reflect.Type, name ...string) (interface{}, error) {
	factory := depMap.GetOrDefault(reflectTypeKey(depType), nameOrDefault(name), nil)
	if factory == nil {
		return nil, ErrorDepFactoryNotFound
	}
	return factory()
}

// InjectT tries to resolve a dependency by its type and optionally its name
// If the dependency is unknown, ErrorDepFactoryNotFound is returned
func InjectInto(out interface{}, name ...string) error {
	argType := reflect.TypeOf(out)
	if argType.Kind() != reflect.Pointer {
		return ErrorDepNotAPointer
	}
	factory := depMap.GetOrDefault(reflectTypeKey(argType.Elem()), nameOrDefault(name), nil)
	if factory == nil {
		return ErrorDepFactoryNotFound
	}
	dep, err := factory()
	if err != nil {
		return err
	}
	reflect.ValueOf(out).Elem().Set(reflect.ValueOf(dep))
	return nil
}

// MustInject tries to resolve a dependency by its type and optionally its name or panics
func MustInject[T any](name ...string) T {
	dep, err := Inject[T](name...)
	if err != nil {
		panic(err)
	}
	return dep
}

// MustInjectT tries to resolve a dependency by its type and optionally its name or panics
func MustInjectT(depType reflect.Type, name ...string) interface{} {
	dep, err := InjectT(depType, name...)
	if err != nil {
		panic(err)
	}
	return dep
}

// MustInjectInto tries to resolve a dependency by its type and optionally its name or panics
func MustInjectInto(out interface{}, name ...string) {
	if err := InjectInto(out, name...); err != nil {
		panic(err)
	}
}

// Call will attempt to resolve all arguments of the function and then call it
func Call(fn interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("dependency injection failed because factory paniced, recovered value: %v", r)
		}
	}()
	args, err := injectFuncArgs(fn)
	if err != nil {
		return
	}
	resVal := reflect.ValueOf(fn).Call(args)
	if len(resVal) > 0 {
		if err, ok := resVal[0].Interface().(error); ok && err != nil {
			return err
		}
	}
	return
}

// MustCall will attempt to resolve all arguments of the function and then call it or panic
func MustCall(fn interface{}) {
	if err := Call(fn); err != nil {
		panic(err)
	}
}

// Fill will attempt to resolve all tagged fields of a struct with their matching dependency
func Fill(strct interface{}) error {
	return injectStructFields(strct)
}

// MustFill will attempt to resolve all tagged fields of a struct with their matching dependency or panic
func MustFill(strct interface{}) {
	if err := Fill(strct); err != nil {
		panic(err)
	}
}
