package factory

import (
	"reflect"
)

type Factory interface {
	ContainsObject(name string) bool
	GetAliases(name string) (aliases []string, err error)
	GetObject(name string, opts ...Options) (obj interface{}, err error)
	GetType(name string) (typ reflect.Type)

	IsPrototype(name string) bool
	IsSingleton(name string) bool
	IsTypeMatch(name string, typ reflect.Type) bool

	RegisterObjectDefinition(definition ObjectDefinition) error
}
