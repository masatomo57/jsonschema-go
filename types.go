package jsonschema

import "reflect"

// JSON Schemaの型定数
const (
	TypeString  = "string"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeBoolean = "boolean"
	TypeObject  = "object"
	TypeArray   = "array"
)

// JSON Schemaのプロパティ名定数
const (
	PropType                 = "type"
	PropProperties           = "properties"
	PropAdditionalProperties = "additionalProperties"
	PropRequired             = "required"
	PropItems                = "items"
	PropMinimum              = "minimum"
	PropMaximum              = "maximum"
	PropExclusiveMinimum     = "exclusiveMinimum"
	PropExclusiveMaximum     = "exclusiveMaximum"
	PropMultipleOf           = "multipleOf"
	PropPattern              = "pattern"
	PropFormat               = "format"
	PropMinItems             = "minItems"
	PropMaxItems             = "maxItems"
)

// getTypeSchema はGoの基本型のJSON Schemaを返す。
func getTypeSchema(t reflect.Type) map[string]any {
	switch t.Kind() {
	case reflect.String:
		return map[string]any{PropType: TypeString}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return map[string]any{PropType: TypeInteger}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return map[string]any{PropType: TypeInteger}
	case reflect.Float32, reflect.Float64:
		return map[string]any{PropType: TypeNumber}
	case reflect.Bool:
		return map[string]any{PropType: TypeBoolean}
	default:
		return nil
	}
}
