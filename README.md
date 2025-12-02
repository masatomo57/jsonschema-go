# jsonschema-go

Goの構造体からJSON Schemaを生成する

## インストール

```bash
go get github.com/masatomo.akagawa/jsonschema-go
```

## 使い方

```go
import "github.com/masatomo.akagawa/jsonschema-go"

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required,minLength=1,maxLength=100"`
	Email string `json:"email" validate:"pattern=^.+@.+$"`
	Age   int    `json:"age,omitempty" validate:"minimum=0,maximum=150"`
}

schema, err := jsonschema.Generate(User{})
```

**出力:**

```json
{
  "type": "object",
  "properties": {
    "id": { "type": "integer" },
    "name": { "type": "string", "minLength": 1, "maxLength": 100 },
    "email": { "type": "string", "pattern": "^.+@.+$" },
    "age": { "type": "integer", "minimum": 0, "maximum": 150 }
  },
  "required": ["id", "name", "email"],
  "additionalProperties": false
}
```

## サポート機能

**データ型:** `string`, `int`, `float64`, `bool`, `[]T`, `map[string]T`, `struct`

**`json`タグ:** `json:"name"`, `json:"name,omitempty"`, `json:"-"`

**`validate`タグ:** `required`, `minimum`, `maximum`, `minLength`, `maxLength`, `pattern`
