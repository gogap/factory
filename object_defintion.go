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

	newObjFunc      NewObjectFunc
	typ             reflect.Type
	refs            map[string]string
	refsOptions     map[string]Options
	refsOrder       []string
	initialFuncName string
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
	if typ.Kind() != reflect.Ptr {
		return typ == p.typ
	}

	return typ.Elem() == p.typ
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
	return p.typ
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

func DefOptOfObjectRef(fieldName string, refDefName string, opts ...Options) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {

		fieldName = strings.TrimSpace(fieldName)

		if fieldName == "" {
			err = ErrEmptyFieldName.New()
			return
		}

		if refDefName == "" {
			err = ErrRefDefinitionNameIsEmpty.New(errors.Params{"name": refDefName})
			return
		}

		if originalRefDefName, exist := od.refs[fieldName]; exist {
			if originalRefDefName != refDefName {
				err = ErrFiledAreadyRef.New(errors.Params{"name": originalRefDefName})
				return
			}
		}

		fieldNames := strings.Split(fieldName, ".")
		lenfields := len(fieldNames)

		typ := od.typ

		for i, fn := range fieldNames {

			fn = strings.TrimSpace(fn)
			if fn == "" {
				err = ErrBadFieldName.New(errors.Params{"name": fieldName})
				return
			}

			var field reflect.StructField
			var exist bool

			for {
				if typ.Kind() == reflect.Ptr {
					typ = typ.Elem()
				} else {
					break
				}
			}

			if field, exist = typ.FieldByName(fn); !exist {
				err = ErrStructFieldNotExist.New(errors.Params{"name": fieldName})
				return
			}

			if i+1 >= lenfields {
				if field.Type.Kind() != reflect.Ptr {
					err = ErrRefFieldShouldBePtr.New()
					return
				}
			}

			typ = field.Type
		}

		od.refs[fieldName] = refDefName
		if opts != nil && len(opts) > 0 {
			od.refsOptions[fieldName] = opts[0]
		}

		od.refsOrder = append(od.refsOrder, fieldName)

		return
	}}
}

func DefOptOfInitialFunc(fnName string) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {
		od.initialFuncName = fnName
		return
	}}
}

func DefOptOfRefOrder(check bool, order ...string) DefinitionOption {
	return DefinitionOption{func(od *ObjectDefinition) (err error) {
		if check {

			tmpOrder := removeDuplicates(order)

			if len(tmpOrder) != len(od.refs) {
				err = ErrBadRefOrderLength.New()
				return
			}

			for _, filedName := range tmpOrder {
				if _, exist := od.refs[filedName]; !exist {
					err = ErrRefOrderContainNonExistRef.New(errors.Params{"name": filedName})
					return
				}
			}
		}

		od.refsOrder = order
		return
	}}
}

func removeDuplicates(elements []string) []string {
	encountered := map[string]bool{}

	for v := range elements {
		encountered[elements[v]] = true
	}

	result := []string{}
	for key, _ := range encountered {
		result = append(result, key)
	}
	return result
}
