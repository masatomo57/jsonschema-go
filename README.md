# jsonschema-go

Goの構造体からJSON Schemaを生成するライブラリです。JSON Schema Draft 6以降に対応しています。

## インストール

```bash
go get github.com/masatomo.akagawa/jsonschema-go
```

## 使い方

```go
import "github.com/masatomo.akagawa/jsonschema-go"

type User struct {
	ID       int      `json:"id" validate:"required,minimum=1,maximum=100"`
	Name     string   `json:"name" validate:"required,pattern=^[a-zA-Z]+$"`
	Email    string   `json:"email" validate:"pattern=^[a-z]+@[a-z]+\\.[a-z]+$,format=email"`
	Age      int      `json:"age,omitempty" validate:"minimum=0,maximum=150"`
	Score    float64  `json:"score" validate:"multipleOf=0.5,exclusiveMinimum=0,exclusiveMaximum=100"`
	Tags     []string `json:"tags" validate:"minItems=1,maxItems=10"`
}

schema, err := jsonschema.Generate(User{})
```

**出力:**

```json
{
  "type": "object",
  "properties": {
    "id": { "type": "integer", "minimum": 1, "maximum": 100 },
    "name": { "type": "string", "pattern": "^[a-zA-Z]+$" },
    "email": { "type": "string", "pattern": "^[a-z]+@[a-z]+\\.[a-z]+$", "format": "email" },
    "age": { "type": "integer", "minimum": 0, "maximum": 150 },
    "score": { "type": "number", "multipleOf": 0.5, "exclusiveMinimum": 0, "exclusiveMaximum": 100 },
    "tags": { "type": "array", "items": { "type": "string" }, "minItems": 1, "maxItems": 10 }
  },
  "required": ["id", "name", "email", "score", "tags"],
  "additionalProperties": false
}
```

## サポート機能

### データ型

- 基本型: `string`, `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uint8`, `uint16`, `uint32`, `uint64`, `float32`, `float64`, `bool`
- 複合型: `[]T` (配列/スライス), `map[string]T` (マップ), `struct` (構造体), `*T` (ポインタ)

### `json`タグ

- `json:"name"` - フィールド名を指定
- `json:"name,omitempty"` - 空値の場合は省略
- `json:"-"` - フィールドをスキップ

### `validate`タグ

#### 数値型の制約

- `minimum=N` - 最小値（以上）
- `maximum=N` - 最大値（以下）
- `exclusiveMinimum=N` - 排他的最小値（より大きい）※Draft 6以降の形式（数値）
- `exclusiveMaximum=N` - 排他的最大値（より小さい）※Draft 6以降の形式（数値）
- `multipleOf=N` - 倍数制約

#### 文字列型の制約

- `pattern=REGEX` - 正規表現パターン
- `format=FORMAT` - フォーマット指定
  - サポート形式: `date-time`, `time`, `date`, `duration`, `email`, `hostname`, `ipv4`, `ipv6`, `uuid`

#### 配列型の制約

- `minItems=N` - 最小要素数
- `maxItems=N` - 最大要素数

#### その他

- `required` - 必須フィールドとしてマーク

## JSON Schemaバージョン

このライブラリは **JSON Schema Draft 6以降** に対応しています。

主な特徴:
- `exclusiveMinimum`と`exclusiveMaximum`は数値型として扱われます（Draft 6以降の形式）
- `items`は単一のスキーマオブジェクトとして使用されます
