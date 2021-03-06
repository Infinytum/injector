package injector

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	DefaultForType = ""

	ErrorDependencyTypeMismatch = errors.New("the resolved dependency does not match the generic type")

	ErrorDepFactoryNotAFunc    = errors.New("the provided dependency factory is not a function")
	ErrorDepFactoryNotFound    = errors.New("the requested type/name combination is not a registered dependency")
	ErrorDepFactoryReturnCount = errors.New("the provided dependency factory must return exactly 1 value")

	ErrorDepNotAPointer = errors.New("the provided value must be a pointer to the struct you want to inject into")
	ErrorDepNotAStruct  = errors.New("the provided value must be a struct")

	ErrorInvalidTag = errors.New("the provided injector tag is not valid")
)

// nameOrDefault will return the default name if the name array is nil or empty
func nameOrDefault(name []string) string {
	depName := DefaultForType
	if name != nil && len(name) > 0 {
		depName = strings.Join(name, "")
	}
	return depName
}

func reflectTypeKey(t reflect.Type) string {
	nameType := t
	for nameType != nil && nameType.Kind() == reflect.Pointer {
		nameType = nameType.Elem()
	}

	pkg, name := "UNKNOWN_PACKAGE", t.String()
	if nameType != nil {
		pkg = nameType.PkgPath()
	}
	return fmt.Sprintf("%s/%s", pkg, name)
}
