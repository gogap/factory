package factory

import (
	"reflect"
	"testing"
)

type testObjectB struct {
	BValue string
}

type testObjectC struct {
	CValue *testObjectB
}

func newTestObjectB(opts Options) (v interface{}, err error) {
	return &testObjectB{BValue: "VB"}, nil
}

type testObject struct {
	isFromNewFunc bool

	ObjB *testObjectB
	ObjC testObjectC
}

func newTestObject(opts Options) (v interface{}, err error) {
	return &testObject{isFromNewFunc: true}, nil
}

func TestClassicFactory(t *testing.T) {

	var objDef ObjectDefinition
	var err error

	if objDef, err = NewObjectDefinition("testObjName", Prototype, new(testObject), DefOptOfNewObjectFunc(newTestObject)); err != nil {
		t.Error(err)
		return
	}

	factory := NewClassicFactory()

	factory.RegisterObjectDefinition(objDef)

	var obj interface{}
	if obj, err = factory.GetObject("testObjName"); err != nil {
		t.Error(err)
		return
	}

	objIns := obj.(*testObject)

	if !objIns.isFromNewFunc {
		t.Error("object is not created by func of newTestObject")
		return
	}
}

func TestClassicFactoryWithoutNewFunc(t *testing.T) {

	var objDef ObjectDefinition
	var err error

	if objDef, err = NewObjectDefinition("testObjName", Prototype, new(testObject)); err != nil {
		t.Error(err)
		return
	}

	factory := NewClassicFactory()

	factory.RegisterObjectDefinition(objDef)

	var obj interface{}
	if obj, err = factory.GetObject("testObjName"); err != nil {
		t.Error(err)
		return
	}

	objIns := obj.(*testObject)

	if objIns.isFromNewFunc {
		t.Error("object is created by func of newTestObject")
		return
	}
}

func TestClassicFactorySingleton(t *testing.T) {

	var objDef ObjectDefinition
	var err error

	if objDef, err = NewObjectDefinition("testObjName", Singleton, new(testObject), DefOptOfNewObjectFunc(newTestObject)); err != nil {
		t.Error(err)
		return
	}

	factory := NewClassicFactory()

	factory.RegisterObjectDefinition(objDef)

	var obj interface{}
	if obj, err = factory.GetObject("testObjName"); err != nil {
		t.Error(err)
		return
	}

	objIns1 := obj.(*testObject)

	var obj2 interface{}
	if obj2, err = factory.GetObject("testObjName"); err != nil {
		t.Error(err)
		return
	}

	objIns2 := obj2.(*testObject)

	v1 := reflect.ValueOf(objIns1)
	v2 := reflect.ValueOf(objIns2)

	if v1.Pointer() != v2.Pointer() {
		t.Error("test singleton failure")
		return
	}
}

func TestClassicFactoryOfObjRef(t *testing.T) {

	var objDef ObjectDefinition
	var objBDef ObjectDefinition
	var err error

	if objBDef, err = NewObjectDefinition("testObjBName", Prototype, new(testObjectB), DefOptOfNewObjectFunc(newTestObjectB)); err != nil {
		t.Error(err)
		return
	}

	if objDef, err = NewObjectDefinition("testObjName",
		Prototype,
		new(testObject),
		DefOptOfNewObjectFunc(newTestObject),
		DefOptOfObjectRef("ObjB", objBDef),
		DefOptOfObjectRef("ObjC.CValue", objBDef)); err != nil {
		t.Error(err)
		return
	}

	factory := NewClassicFactory()

	factory.RegisterObjectDefinition(objDef)

	var obj interface{}
	if obj, err = factory.GetObject("testObjName"); err != nil {
		t.Error(err)
		return
	}

	objIns := obj.(*testObject)

	if objIns.ObjB == nil {
		t.Error("inject ObjB failure")
		return
	}

	if objIns.ObjB.BValue != "VB" {
		t.Error("field of ObjB ref object B's value is not 'VB'")
		return
	}

	if objIns.ObjC.CValue.BValue != "VB" {
		t.Error("field of ObjC.CValue ref object B's value is not 'VB'")
		return
	}

}
