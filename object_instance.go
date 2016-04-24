package factory

import (
	"fmt"
)

type ObjectInstance struct {
	id         string
	object     interface{}
	options    Options
	definition *ObjectDefinition
}

func (p *ObjectInstance) String() string {
	return fmt.Sprintf("<name: %s, id: %s, type: %s>", p.definition.Name(), p.id, p.definition.Type().String())
}

func (p *ObjectInstance) Id() string {
	return p.id
}

func (p *ObjectInstance) Instance() interface{} {
	return p.object
}

func (p *ObjectInstance) Options() Options {
	return p.options
}

func (p *ObjectInstance) Definition() ObjectDefinition {
	return *p.definition
}
