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

type testObject struct {
	isFromNewFunc bool

	ObjB *testObjectB
	ObjC testObjectC
}

func init() {
	RegisterModel((*testObjectB)(nil), "testObjectB")
	RegisterModel((*testObjectC)(nil), "testObjectC")
	RegisterModel((*testObject)(nil), "testObject")
}

func newTestObjectB(opts Options) (v interface{}, err error) {
	return &testObjectB{BValue: "VB"}, nil
}

func newTestObject(opts Options) (v interface{}, err error) {
	return &testObject{isFromNewFunc: true}, nil
}

func TestClassicFactory(t *testing.T) {

	var err error

	factory := NewClassicFactory(nil)

	if err = factory.Define("testObjName", Prototype, "testObjectB", DefOptOfNewObjectFunc(newTestObject)); err != nil {
		t.Error(err)
		return
	}

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

	var err error

	factory := NewClassicFactory(nil)

	if err = factory.Define("testObjName", Prototype, "testObject"); err != nil {
		t.Error(err)
		return
	}

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

	var err error
	factory := NewClassicFactory(nil)

	if err = factory.Define("testObjName", Singleton, "testObject", DefOptOfNewObjectFunc(newTestObject)); err != nil {
		t.Error(err)
		return
	}

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

	var err error

	factory := NewClassicFactory(nil)

	if err = factory.Define("testObjBName", Prototype, "testObjectB", DefOptOfNewObjectFunc(newTestObjectB)); err != nil {
		t.Error(err)
		return
	}

	if err = factory.Define("testObjName",
		Prototype,
		"testObject",
		DefOptOfNewObjectFunc(newTestObject),
		DefOptOfObjectRef("ObjB", "testObjBName"),
		DefOptOfObjectRef("ObjC.CValue", "testObjBName")); err != nil {
		t.Error(err)
		return
	}

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
