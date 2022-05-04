package injector

import (
	"fmt"
	"reflect"
)

type FactoryFunction func() (interface{}, error)

// TransientFactory is a function wrapper for a transient dependency
// to provide dependency injection inside the factory function
func TransientFactory(depFactory interface{}) (FactoryFunction, error) {
	depFactoryType := reflect.TypeOf(depFactory)
	depFactoryValue := reflect.ValueOf(depFactory)

	if depFactoryType.Kind() != reflect.Func {
		return nil, ErrorDepFactoryNotAFunc
	}

	if depFactoryType.NumOut() != 1 {
		return nil, ErrorDepFactoryReturnCount
	}

	// Dependency-Inject arguments for factory function
	args, err := injectFuncArgs(depFactory)
	if err != nil {
		return nil, err
	}

	// Create factory wrapper to inject dependencies
	return func() (dep interface{}, err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("dependency injection failed because factory paniced, recovered value: %v", r)
				dep = nil
			}
		}()
		returnVals := depFactoryValue.Call(args)
		dep = returnVals[0].Interface()
		return
	}, nil
}

// SingletonFactory is a function wrapper for a singleton dependency factory
// to provide dependency injection inside the factory function and to retain
// the singleton instance once instantiated.
func SingletonFactory(depFactory interface{}) (FactoryFunction, error) {
	factory, err := TransientFactory(depFactory)
	if err != nil {
		return nil, err
	}

	// Wrap factory wrapper to ensure existing singleton value is used if
	// it already exists.
	var singleton interface{}
	return func() (interface{}, error) {
		if singleton != nil {
			return singleton, nil
		}

		// Singleton is not ready, call transient factory to instantiate dependency
		dep, err := factory()
		if err != nil {
			return nil, err
		}
		singleton = dep
		return singleton, nil
	}, nil
}
