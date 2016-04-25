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
	ErrObjectMustBeStruct                = errors.TN(ErrNamespace, 1007, "object must be struct, name: {{.name}}")
	ErrCouldNotSetFiledOfNilObject       = errors.TN(ErrNamespace, 1008, "object is nil, could not inject filed value, field: {{.field}}")
	ErrCouldNotSetZeroNumFieldObject     = errors.TN(ErrNamespace, 1009, "file number is zero, could not inject filed value, filed: {{.field}}")
	ErrObjectIsNotStruct                 = errors.TN(ErrNamespace, 1010, "the object must be a struct or ptr to struct")
	ErrRefObjectShouldBePtr              = errors.TN(ErrNamespace, 1011, "ref object should be ptr")
	ErrRefFieldShouldBePtr               = errors.TN(ErrNamespace, 1012, "ref field should be ptr")
	ErrEmptyFieldName                    = errors.TN(ErrNamespace, 1013, "empty field name")
	ErrBadFieldName                      = errors.TN(ErrNamespace, 1014, "bad field name, field name: {{.name}}")
	ErrModelAlreayRegistered             = errors.TN(ErrNamespace, 1015, "model already registered, name: {{.name}}, type: {{.type}}")
	ErrModleAliasAlreadyExist            = errors.TN(ErrNamespace, 1016, "model alias already exist and model name not match, original name: {{.originalName}}, new name: {{.newName}}")
	ErrModelNameIsEmpty                  = errors.TN(ErrNamespace, 1017, "model name is empty")
	ErrModelNotExist                     = errors.TN(ErrNamespace, 1018, "model of {{.name}} not exist")
	ErrRefDefinitionNameIsEmpty          = errors.TN(ErrNamespace, 1019, "ref definition name is empty, def name: {{.name}}")
	ErrFiledAreadyRef                    = errors.TN(ErrNamespace, 1020, "field already ref others definition, original ref defition name: {{.name}}")
	ErrFieldIsZeroValue                  = errors.TN(ErrNamespace, 1021, "filed is zero value, field name: {{.name}}")
	ErrBadRefOrderLength                 = errors.TN(ErrNamespace, 1022, "ref order does not equal definition refs")
	ErrRefOrderContainNonExistRef        = errors.TN(ErrNamespace, 1023, "ref order contain non exist def ref, name: {{.name}}")
)
