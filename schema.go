package jsonschema

import (
	"fmt"
	"reflect"
)

// Generate はGoの構造体からJSON Schemaを生成する。
func Generate(v any) (map[string]any, error) {
	if v == nil {
		return nil, fmt.Errorf("cannot generate schema from nil value")
	}

	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected struct type, got %s", t.Kind())
	}

	return generateSchema(t), nil
}

// generateSchema はreflect.TypeからJSON Schemaを再帰的に生成する。
func generateSchema(t reflect.Type) map[string]any {
	// ポインタ型の処理
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// 基本型の処理
	if schema := getTypeSchema(t); schema != nil {
		return schema
	}

	// スライスと配列の処理
	if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
		return map[string]any{
			PropType:  TypeArray,
			PropItems: generateSchema(t.Elem()),
		}
	}

	// マップの処理
	if t.Kind() == reflect.Map {
		if t.Key().Kind() != reflect.String {
			return map[string]any{
				PropType: TypeObject,
			}
		}
		return map[string]any{
			PropType:                 TypeObject,
			PropAdditionalProperties: generateSchema(t.Elem()),
		}
	}

	// 構造体の処理
	if t.Kind() == reflect.Struct {
		return generateStructSchema(t)
	}

	// 未知の型のフォールバック
	return map[string]any{
		PropType: TypeObject,
	}
}

// generateStructSchema は構造体型のJSON Schemaを生成する。
func generateStructSchema(t reflect.Type) map[string]any {
	schema := map[string]any{
		PropType:                 TypeObject,
		PropProperties:           map[string]any{},
		PropAdditionalProperties: false,
	}

	var required []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 非公開フィールドをスキップ
		if !field.IsExported() {
			continue
		}

		// jsonタグを解析
		jsonName, omitempty, skip := parseJSONTag(field.Tag.Get("json"))
		if skip {
			continue
		}

		// jsonタグがあればその名前を使用、なければフィールド名を使用
		if jsonName == "" {
			jsonName = field.Name
		}

		// フィールドの型からスキーマを生成
		fieldSchema := generateSchema(field.Type)

		// validateタグを解析
		validationTag := field.Tag.Get("validate")
		if validationTag != "" {
			constraints := parseValidationTag(validationTag)
			for key, value := range constraints {
				fieldSchema[key] = value
			}
		}

		// requiredフィールドかどうかを判定
		if isRequiredField(field, omitempty, validationTag) {
			required = append(required, jsonName)
		}

		// プロパティをスキーマに追加
		props := schema[PropProperties].(map[string]any)
		props[jsonName] = fieldSchema
	}

	if len(required) > 0 {
		schema[PropRequired] = required
	}

	return schema
}
