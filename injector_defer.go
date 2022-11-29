package injector

import (
	"reflect"

	"github.com/infinytum/structures"
)

var (
	deferredFactories structures.Table[string, string, *deferredFactory] = structures.NewTable[string, string, *deferredFactory]()
)

type deferredFactory struct {
	IsSingleton bool
	Resolver    interface{}
}

// DeferredSingleton will register a dependency that is only instantiated once, then re-used
// for all future resolve calls. The resolver function will be introspected at the first time its needed.
func DeferredSingleton(resolver interface{}, name ...string) error {
	depFactoryType := reflect.TypeOf(resolver)
	return deferredFactories.Set(reflectTypeKey(depFactoryType.Out(0)), nameOrDefault(name), &deferredFactory{
		IsSingleton: true,
		Resolver:    resolver,
	})
}

// DeferredTransient will register a dependency that is instantiated every time it's resolved.
// The resolver function will be introspected at the first time its needed.
func DeferredTransient(resolver interface{}, name ...string) error {
	depFactoryType := reflect.TypeOf(resolver)
	return deferredFactories.Set(reflectTypeKey(depFactoryType.Out(0)), nameOrDefault(name), &deferredFactory{
		IsSingleton: false,
		Resolver:    resolver,
	})
}

func activateDeferredFactories(t reflect.Type, name ...string) error {
	deferredFactory, ok := deferredFactories.Get(reflectTypeKey(t), nameOrDefault(name))
	if ok == structures.TableKeysNotFound {
		return nil
	}
	if ok != nil {
		return ok
	}
	defer deferredFactories.Delete(reflectTypeKey(t), nameOrDefault(name))
	if deferredFactory.IsSingleton {
		return Singleton(deferredFactory.Resolver, name...)
	}
	return Transient(deferredFactory.Resolver, name...)
}
