package factory

import (
	"fmt"
	"github.com/gogap/errors"
	"reflect"
	"sync"
)

var defaultModelProvider ModelProvider = NewClassicModelProvider()

type ModelProvider interface {
	Register(model interface{}, aliases ...string) (err error)
	Get(name string) (typ reflect.Type, exist bool)
}

type ClassicModelProvider struct {
	modelTypes map[string]reflect.Type
	aliases    map[string]string

	locker sync.Mutex
}

func NewClassicModelProvider() ModelProvider {
	return &ClassicModelProvider{
		aliases:    make(map[string]string),
		modelTypes: make(map[string]reflect.Type),
	}
}

func RegisterModel(model interface{}, aliases ...string) (err error) {
	if err = defaultModelProvider.Register(model, aliases...); err != nil {
		return
	}
	return
}

func (p *ClassicModelProvider) Register(model interface{}, aliases ...string) (err error) {
	p.locker.Lock()
	defer p.locker.Unlock()

	typ := reflect.TypeOf(model)

	for {
		if typ.Kind() == reflect.Ptr {
			typ = typ.Elem()
		} else {
			break
		}
	}

	name := fmt.Sprintf("%s::%s", typ.PkgPath(), typ.String())

	if originalType, exist := p.modelTypes[name]; exist {
		err = ErrModelAlreayRegistered.New(errors.Params{"name": name, "type": originalType.String()})
		return
	}

	p.modelTypes[name] = typ

	for _, alias := range aliases {
		if originalName, exist := p.aliases[alias]; exist {
			if originalName != alias {
				err = ErrModleAliasAlreadyExist.New(errors.Params{"originalName": originalName, "newName": alias})
				return
			}
		}
		p.aliases[alias] = name
	}

	p.aliases[name] = name

	return
}

func (p *ClassicModelProvider) Get(name string) (typ reflect.Type, exist bool) {
	var alias string

	if alias, exist = p.aliases[name]; exist {
		typ, exist = p.modelTypes[alias]
		return
	}

	return
}
