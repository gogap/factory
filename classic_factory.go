package factory

import (
	"github.com/gogap/errors"
	"github.com/rs/xid"
	"reflect"
	"strings"
	"sync"
)

type ClassicFactory struct {
	objLocker sync.Mutex

	objDefines   map[string]*ObjectDefinition
	objAliases   map[string]string
	objInstances map[string]*ObjectInstance
}

func NewClassicFactory() *ClassicFactory {
	return &ClassicFactory{
		objDefines:   make(map[string]*ObjectDefinition),
		objAliases:   make(map[string]string),
		objInstances: make(map[string]*ObjectInstance),
	}
}

func (p *ClassicFactory) ContainsObject(name string) bool {
	var err error
	if _, err = p.getObjDefinition(name); err != nil {
		return false
	}

	return true
}

func (p *ClassicFactory) GetAliases(name string) (aliases []string, err error) {

	var def *ObjectDefinition
	if def, err = p.getObjDefinition(name); err != nil {
		return
	}

	aliases = def.Aliases()
	if name != def.Name() {
		aliases = append(aliases, def.Name())
	}

	return
}

func (p *ClassicFactory) GetObject(name string, opts ...Options) (obj interface{}, err error) {
	var def *ObjectDefinition

	if def, err = p.getObjDefinition(name); err != nil {
		return
	}

	var opt Options
	if opts != nil || len(opts) > 0 {
		opt = opts[0]
	}

	if obj, err = p.getObject(def, opt); err != nil {
		return
	}

	return
}

func (p *ClassicFactory) GetType(name string) (typ reflect.Type) {

	var def *ObjectDefinition
	var err error
	if def, err = p.getObjDefinition(name); err != nil {
		return
	}

	typ = def.Type()
	return

}

func (p *ClassicFactory) IsPrototype(name string) bool {

	var def *ObjectDefinition
	var err error
	if def, err = p.getObjDefinition(name); err != nil {
		return def.Scope() == Prototype
	}

	return false
}

func (p *ClassicFactory) IsSingleton(name string) bool {
	var def *ObjectDefinition
	var err error
	if def, err = p.getObjDefinition(name); err != nil {
		return def.Scope() == Singleton
	}

	return false
}

func (p *ClassicFactory) IsTypeMatch(name string, typ reflect.Type) bool {
	var def *ObjectDefinition
	var err error
	if def, err = p.getObjDefinition(name); err != nil {
		return def.IsTypeMatch(typ)
	}

	return false
}

func (p *ClassicFactory) RegisterObjectDefinition(definition ObjectDefinition) (err error) {

	p.objLocker.Lock()
	defer p.objLocker.Unlock()

	if _, exist := p.objDefines[definition.Name()]; exist {
		err = ErrObjectDefinitionAlreadyRegistered.New(errors.Params{"name": definition.Name(), "type": definition.Type()})
		return
	}

	p.objDefines[definition.Name()] = &definition

	return
}

func (p *ClassicFactory) getObjDefinition(name string) (def *ObjectDefinition, err error) {
	var exist bool

	if def, exist = p.objDefines[name]; exist {
		return
	}

	var originalName string
	if originalName, exist = p.objAliases[name]; exist {
		if def, exist = p.objDefines[originalName]; exist {
			return
		}
	}

	err = ErrObjectDefintionNotExist.New(errors.Params{"name": name})
	return
}

func (p *ClassicFactory) getObject(def *ObjectDefinition, opts Options) (obj interface{}, err error) {
	// Get ref objects
	var refObjs = make(map[string]interface{})
	for fieldName, refDef := range def.refs {
		var o interface{}
		if o, err = p.getObject(refDef, opts); err != nil {
			return
		}
		refObjs[fieldName] = o
	}

	var retObj interface{}

	if def.Scope() == Singleton {
		var exist bool

		p.objLocker.Lock()
		defer p.objLocker.Unlock()

		var objIns *ObjectInstance
		if objIns, exist = p.objInstances[def.Name()]; exist {
			obj = objIns.Instance()
			return
		}

		// Create new object
		var newInstanceFn NewObjectFunc

		if newInstanceFn, err = p.getNewInstanceFunc(def); err != nil {
			return
		}

		if retObj, err = newInstanceFn(opts); err != nil {
			return
		}

		p.objInstances[def.Name()] = &ObjectInstance{
			id:         xid.New().String(),
			object:     retObj,
			options:    opts,
			definition: def,
		}
	}

	if def.Scope() == Prototype {

		var newInstanceFn NewObjectFunc

		if newInstanceFn, err = p.getNewInstanceFunc(def); err != nil {
			return
		}

		if retObj, err = newInstanceFn(opts); err != nil {
			return
		}
	}

	// Inject dependency object
	for fieldName, fieldValue := range refObjs {
		if err = p.setStructFieldValue(retObj, fieldName, fieldValue); err != nil {
			return
		}
	}

	obj = retObj

	return
}

func (p *ClassicFactory) getNewInstanceFunc(def *ObjectDefinition) (fn NewObjectFunc, err error) {
	fn = def.NewObjectFunc()

	if fn != nil {
		return
	}

	if fn, err = p.newTypeInstance(def.Type()); err != nil {
		return
	}

	return
}

func (p *ClassicFactory) newTypeInstance(typ reflect.Type) (fn NewObjectFunc, err error) {

	val := reflect.New(typ.Elem())

	if !val.IsValid() {
		err = ErrReflectValueNotValid.New()
		return
	}

	fn = func(_ Options) (v interface{}, err error) {
		v = val.Interface()
		return
	}

	return
}

func (p *ClassicFactory) setStructFieldValue(v interface{}, fieldName string, fieldValue interface{}) (err error) {

	if v == nil {
		err = ErrCouldNotSetFiledOfNilObject.New(errors.Params{"field": fieldName})
		return
	}

	val := reflect.ValueOf(v)
	if !val.IsValid() {
		err = ErrReflectValueNotValid.New()
		return
	}

	for {
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		} else {
			break
		}
	}

	if val.Kind() != reflect.Struct {
		err = ErrObjectIsNotStruct.New()
		return
	}

	if val.NumField() == 0 {
		err = ErrCouldNotSetZeroNumFieldObject.New(errors.Params{"field": fieldName})
		return
	}

	fieldNames := strings.Split(fieldName, ".")
	lenfields := len(fieldNames)

	var fieldVal reflect.Value

	for i, fn := range fieldNames {
		if fieldVal = val.FieldByName(fn); !fieldVal.IsValid() {
			err = ErrReflectValueNotValid.New()
			return
		}

		if i+1 >= lenfields && fieldVal.Kind() != reflect.Ptr {
			err = ErrRefObjectShouldBePtr.New()
			return
		}

		val = fieldVal
	}

	newVal := reflect.ValueOf(fieldValue)

	if newVal.Kind() == reflect.Ptr {
		fieldVal.Set(newVal)
	} else if newVal.Kind() == reflect.Struct {
		fieldVal.Set(reflect.Indirect(newVal))
	}

	return
}
