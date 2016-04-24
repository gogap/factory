package factory

import (
	"github.com/gogap/errors"
	"reflect"
	"strings"
)

type Scope int

const (
	Singleton Scope = 0
	Prototype Scope = 1
)

type NewObjectFunc func(opts Options) (v interface{}, err error)

type DefinitionOption struct {
	f func(o *ObjectDefinition) (err error)
}

type ObjectDefinition struct {
	name    string
	aliases []string

	scope Scope

	obj             interface{}
	newObjFunc      NewObjectFunc
	reflectVal      reflect.Value
	refs            map[string]*ObjectDefinition
	initialFuncName string
}

func NewObjectDefinition(
	name string,
	scope Scope,
	obj interface{},
	opts ...DefinitionOption) (objDef ObjectDefinition, err error) {

	name = strings.TrimSpace(name)
	if name == "" {
		err = ErrEmptyObjectDefinitionName.New()
		return
	}

	originalVal := reflect.ValueOf(obj)
	val := originalVal
	for {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		} else {
			break
		}
	}

	if val.Kind() != reflect.Struct {
		err = ErrObjectMustBeStruct.New(errors.Params{"name": name})
		return
	}

	def := ObjectDefinition{
		name:       name,
		scope:      scope,
		obj:        obj,
		reflectVal: originalVal,
		refs:       make(map[string]*ObjectDefinition),
	}

	if err = def.options(opts...); err != nil {
		return
	}

	objDef = def

	return
}

func (p *ObjectDefinition) Name() string {
	return p.name
}

func (p *ObjectDefinition) String() string {
	return p.name
}

func (p *ObjectDefinition) Scope() Scope {
	return p.scope
}

func (p *ObjectDefinition) IsTypeMatch(typ reflect.Type) bool {
	return typ == p.reflectVal.Type()
}

func (p *ObjectDefinition) NewObjectFunc() NewObjectFunc {
	return p.newObjFunc
}

func (p *ObjectDefinition) InitialFuncName() string {
	return p.initialFuncName
}

func (p *ObjectDefinition) Aliases() []string {
	return p.aliases
}

func (p *ObjectDefinition) Type() reflect.Type {
	return p.reflectVal.Type()
}

func (p *ObjectDefinition) options(opts ...DefinitionOption) (err error) {
	if opts == nil {
		return
	}

	for _, opt := range opts {
		if err = opt.f(p); err != nil {
			return
		}
	}
	return
}

func DefOptOfNewObjectFunc(fn NewObjectFunc) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {
		od.newObjFunc = fn
		return
	}}
}

func DefOptOfObjectRef(fieldName string, ref ObjectDefinition) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {

		if field, exist := od.reflectVal.Type().Elem().FieldByName(fieldName); !exist {
			err = ErrStructFieldNotExist.New(errors.Params{"name": fieldName})
			return
		} else if field.Type.Kind() != reflect.Ptr {
			err = ErrRefFieldShouldBePtr.New()
			return
		} else if !ref.IsTypeMatch(field.Type) {
			err = ErrRefTypeNotMatch.New(errors.Params{"typeA": ref.Type().String(), "typeB": field.Type.String()})
			return
		}

		od.refs[fieldName] = &ref

		return
	}}
}

func DefOptOfInitialFunc(fnName string) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {
		od.initialFuncName = fnName
		return
	}}
}
