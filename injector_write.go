package injector

import (
	"reflect"

	"github.com/infinytum/structures"
)

var (
	depMap structures.Table[string, string, FactoryFunction] = structures.NewTable[string, string, FactoryFunction]()
)

// Singleton will register a dependency that is only instantiated once, then re-used
// for all future resolve calls
func Singleton(resolver interface{}, name ...string) error {
	factory, err := SingletonFactory(resolver)
	if err != nil {
		return err
	}
	depFactoryType := reflect.TypeOf(resolver)
	return depMap.Set(reflectTypeKey(depFactoryType.Out(0)), nameOrDefault(name), factory)
}

// Transient will register a dependency that is instantiated every time it's resolved.
func Transient(resolver interface{}, name ...string) error {
	factory, err := TransientFactory(resolver)
	if err != nil {
		return err
	}
	depFactoryType := reflect.TypeOf(resolver)
	return depMap.Set(reflectTypeKey(depFactoryType.Out(0)), nameOrDefault(name), factory)
}
