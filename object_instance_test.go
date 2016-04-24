package factory

import (
	"fmt"
	"testing"

	"github.com/rs/xid"
)

func TestObjectInstance(t *testing.T) {
	id := xid.New().String()

	var objDef ObjectDefinition
	var err error

	if objDef, err = NewObjectDefinition("testObjName", Prototype, new(testObject)); err != nil {
		t.Error(err)
		return
	}

	objIns := &ObjectInstance{
		id:         id,
		object:     testObject{},
		options:    nil,
		definition: &objDef,
	}

	if objIns.Id() != id {
		t.Error("id not match")
		return
	}

	objStr := fmt.Sprintf("<name: %s, id: %s, type: %s>", objIns.definition.name, objIns.id, objIns.definition.Type().String())

	if objIns.String() != objStr {
		t.Error("object String() error")
		return
	}
}
