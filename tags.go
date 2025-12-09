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
//   - exclusiveMinimum=N: 数値がこの値より大きい必要がある
//   - exclusiveMaximum=N: 数値がこの値より小さい必要がある
//   - multipleOf=N: 数値がこの値の倍数である必要がある
//   - pattern=REGEX: 文字列のパターンを設定
//   - format=FORMAT: 文字列のフォーマットを設定（date-time, time, date, duration, email, hostname, ipv4, ipv6, uuid）
//   - minItems=N: 配列の最小要素数を設定
//   - maxItems=N: 配列の最大要素数を設定
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
			case "exclusiveMinimum":
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					constraints[PropExclusiveMinimum] = num
				}
			case "exclusiveMaximum":
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					constraints[PropExclusiveMaximum] = num
				}
			case "multipleOf":
				if num, err := strconv.ParseFloat(value, 64); err == nil {
					constraints[PropMultipleOf] = num
				}
			case "pattern":
				constraints[PropPattern] = value
			case "format":
				// サポートされているフォーマット: date-time, time, date, duration, email, hostname, ipv4, ipv6, uuid
				validFormats := map[string]bool{
					"date-time": true,
					"time":      true,
					"date":      true,
					"duration":  true,
					"email":     true,
					"hostname":  true,
					"ipv4":      true,
					"ipv6":      true,
					"uuid":      true,
				}
				if validFormats[value] {
					constraints[PropFormat] = value
				}
			case "minItems":
				if num, err := strconv.Atoi(value); err == nil {
					constraints[PropMinItems] = num
				}
			case "maxItems":
				if num, err := strconv.Atoi(value); err == nil {
					constraints[PropMaxItems] = num
				}
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
