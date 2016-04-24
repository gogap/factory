package factory

import (
	"github.com/gogap/errors"
)

const (
	ErrNamespace = "gogap-factory"
)

var (
	ErrEmptyObjectDefinitionName         = errors.TN(ErrNamespace, 1000, "empty object definition name")
	ErrObjectDefinitionAlreadyRegistered = errors.TN(ErrNamespace, 1001, "object definition already registered, name: {{.name}}, type: {{.type}}")
	ErrObjectNotExist                    = errors.TN(ErrNamespace, 1002, "object not exist, name: {{.name}}")
	ErrObjectDefintionNotExist           = errors.TN(ErrNamespace, 1003, "object definition not exist, name: {{.name}}")
	ErrReflectValueNotValid              = errors.TN(ErrNamespace, 1004, "reflect value not valid")
	ErrStructFieldNotExist               = errors.TN(ErrNamespace, 1005, "struct field not exist, field name: {{.name}}")
	ErrRefTypeNotMatch                   = errors.TN(ErrNamespace, 1006, "ref type not match, typeA: {{.typeA}}, typeB: {{.typeB}}")
	ErrObjectMustBeStruct                = errors.TN(ErrNamespace, 1007, "object must be struct, name: {{.name}}")
	ErrCouldNotSetFiledOfNilObject       = errors.TN(ErrNamespace, 1008, "object is nil, could not inject filed value, field: {{.field}}")
	ErrCouldNotSetZeroNumFieldObject     = errors.TN(ErrNamespace, 1009, "file number is zero, could not inject filed value, filed: {{.field}}")
	ErrObjectIsNotStruct                 = errors.TN(ErrNamespace, 1010, "the object must be a struct or ptr to struct")
	ErrRefObjectShouldBePtr              = errors.TN(ErrNamespace, 1011, "ref object should be ptr")
	ErrRefFieldShouldBePtr               = errors.TN(ErrNamespace, 1012, "ref field should be ptr")
)
