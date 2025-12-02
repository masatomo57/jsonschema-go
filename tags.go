package jsonschema

import (
	"reflect"
	"strconv"
	"strings"
)

// parseJSONTag はjsonタグ文字列を解析し、名前、omitemptyフラグ、スキップフラグを返す。
// 例:
//   - `json:"name"` -> ("name", false, false)
//   - `json:"name,omitempty"` -> ("name", true, false)
//   - `json:"-"` -> ("", false, true)
func parseJSONTag(tag string) (name string, omitempty bool, skip bool) {
	if tag == "" {
		return "", false, false
	}

	if tag == "-" {
		return "", false, true
	}

	parts := strings.Split(tag, ",")
	name = parts[0]

	for _, part := range parts[1:] {
		if strings.TrimSpace(part) == "omitempty" {
			omitempty = true
		}
	}

	return name, omitempty, false
}

// parseValidationTag はvalidateタグを解析し、JSON Schema制約のマップを返す。
// サポートする制約:
//   - required: required配列に追加（別途処理）
//   - minimum=N: 数値の最小値を設定
//   - maximum=N: 数値の最大値を設定
//   - minLength=N: 文字列の最小長を設定
//   - maxLength=N: 文字列の最大長を設定
//   - pattern=REGEX: 文字列のパターンを設定
func parseValidationTag(tag string) map[string]any {
	constraints := make(map[string]any)

	if tag == "" {
		return constraints
	}

	parts := strings.Split(tag, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		// key=value形式の制約を処理
		if strings.Contains(part, "=") {
			kv := strings.SplitN(part, "=", 2)
			key := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])

			switch key {
			case "minimum", "min":
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					constraints[PropMinimum] = num
				}
			case "maximum", "max":
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					constraints[PropMaximum] = num
				}
			case "minLength":
				if num, err := strconv.Atoi(value); err == nil {
					constraints[PropMinLength] = num
				}
			case "maxLength":
				if num, err := strconv.Atoi(value); err == nil {
					constraints[PropMaxLength] = num
				}
			case "pattern":
				constraints[PropPattern] = value
			}
		}
		// NOTE: "required"はisRequiredFieldで別途処理
	}

	return constraints
}

// isRequiredField はフィールドがrequiredかどうかを判定する。
// 判定ルール:
//   - json:omitempty が指定されている場合は required 扱いしない
//   - validate:required が明示的に指定されている場合は required
//   - validateタグに omitempty がなく、かつポインタ型でない場合は required
func isRequiredField(field reflect.StructField, omitempty bool, validationTag string) bool {
	// json:omitempty が指定されている場合はrequired扱いしない
	if omitempty {
		return false
	}

	// validate:required が明示的に指定されている
	if strings.Contains(validationTag, "required") {
		return true
	}

	// validateタグにomitemptyがない、かつポインタ型でない場合はrequired
	if !strings.Contains(validationTag, "omitempty") && field.Type.Kind() != reflect.Ptr {
		return true
	}

	return false
}
